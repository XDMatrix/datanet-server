package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)


func createBucket(w io.Writer, projectID, bucketName string, storageClass, region, location string) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}

	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, 30 * time.Second)
	defer cancel()

	bucketAttrs := &storage.BucketAttrs{
		StorageClass: storageClass,
		Location: location,
		CustomPlacementConfig: &storage.CustomPlacementConfig{
			DataLocations: []string{ region1 },
		},
	}

	bucket := client.Bucket(bucketName)
	if err := bucket.Create(ctx, projectID, bucketAttrs); err != nil {
		return fmt.Errorf("Bucket(%q). Create: %v", bucketName, err)
	}

	fmt.Fprintf(w, "Bucket %v created\n", bucketName)
	fmt.Fprintf(w, " - storageClass: %v", bucketAttrs.StorageClass)
	fmt.Fprintf(w, " - location: %v", bucketAttrs.Location)
	fmt.Fprintf(w, " - locationType: %v", bucketAttrs.LocationType)
	fmt.Fprintf(w, " - customPlacementConfig.dataLocations: %v", bucketAttrs.CustomPlacementConfig.DataLocations)

	return nil
}