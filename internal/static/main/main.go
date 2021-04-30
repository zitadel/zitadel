package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/caos/zitadel/internal/static/s3"
)

func main() {
	config := s3.S3Config{
		Endpoint:        "storage.googleapis.com",
		AccessKeyID:     "",
		SecretAccessKey: "",
		SSL:             true,
		Location:        "europe-west6",
	}
	ctx := context.Background()
	minio, err := s3.NewMinio(config)
	if err != nil {
		log.Fatalln(err)
	}

	// Make a new bucket called mymusic.
	bucketName := "hodor2"
	location := "europe-west6"

	err = minio.CreateBucket(ctx, bucketName, location)
	if err != nil {
		log.Printf("Error in create bucket %s: %v", bucketName, err)
	}

	file, err := os.Open("gigi01-sep.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		log.Printf("Error in create bucket %s: %v", bucketName, err)
		return
	}

	info, err := minio.PutObject(ctx, bucketName, file.Name(), "application/octet-stream", file, fileStat.Size())
	if err != nil {
		log.Printf("Error in put object %s: %v", file.Name(), err)
	}
	log.Printf("Object Info: %v", info)
}
