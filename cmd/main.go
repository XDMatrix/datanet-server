

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"cloud.google.com/go/logging"
	"example.com/micro/metadata"
	"github.com/gorilla/mux"

	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang/org/grpc/credentials/insecure"
)

type App struct {
	*http.Server
	projectID string
	log *logging.Logger
}

func main() {
	ctx := context.Background()
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on PORT %s", port)
	porjectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	app, err := newApp(ctx, port, projectID)
	if err != nil {
		log.Fatalf("Unable to initialize application : %v", err)
	}

	log.Println("Starting HTTP Server")

	go func() {
		if err := app.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server closed : %v", err)
		}
	}()

	nctx, stop := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer stop()
	<-nctx.Done()

	log.Println("Shutdown initiated")

	ctx, cancel := context.WithTimeout(ctx, 10 * time.Second)
	defer cancel()
	app.Shutdown(ctx)
	log.Println("Shutdown")
}

func newApp(ctx context.Context, port, projectID string) (*App, error) {
	app := &App{
		Server: &http.Server {
			Addr: ":" + port,
			ReadTimeout: 10 * time.Second,
			WriteTimeout: 10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}

	if projectID == "" {
		projID, err := metadata.ProjectID()
		if err != nil {
			return nil, fmt.Errorf("Unable to detect Project ID from GOOGLE_CLOUD_PROJECT or metadata server : %w", err)
		}
		projectID = projID
	}

	app.projectID = projectID

	client, err := logging.NewClient(ctx, fmt.Sprintf("projects/%s", app.projectID), option.WithoutAuthentication(), option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials()),))
	if err != nil {
		return nil, fmt.Errorf("Unable to initialize logging client : %v", err)
	}

	app.log = client.Logger("test-log", logging.RedirectAsJSON(os.Stderr))

	r := mux.NewRouter()
	r.HandleFunc("/", app.Handler).Methods("GET")
	app.Server.Handler = r

	return app, nil
}