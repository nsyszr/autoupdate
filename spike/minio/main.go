package main

import (
	"fmt"
	"log"

	"github.com/minio/minio-go"
)

func main() {
	endpoint := "127.0.0.1:9000"
	accessKeyID := "VI460PQTSGN8B5YOWASB"
	secretAccessKey := "zacft43zIfvHUpN6PcgYUkOV3nc8RfTwkyEX18Hy"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	// Make a new bucket called mymusic.
	bucketName := "packages-1234"
	//location := "us-east-1"

	/*err = minioClient.MakeBucket(bucketName, location)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, err := minioClient.BucketExists(bucketName)
		if err == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	}
	log.Printf("Successfully created %s\n", bucketName)

	// Upload the zip file
	objectName := "firmware/update-2.7-full.tar"
	filePath := "/Users/tlx3m3j/Downloads/update-2.7-full.tar"
	contentType := "application/tar"

	// Upload the zip file with FPutObject
	n, err := minioClient.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)*/

	// Create a done channel to control 'ListObjects' go routine.
	doneCh := make(chan struct{})

	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)

	// List all objects from a bucket-name with a matching prefix.
	for object := range minioClient.ListObjectsV2(bucketName, "/firmware", true, doneCh) {
		if object.Err != nil {
			fmt.Println(object.Err)
			return
		}
		fmt.Println(object)
	}
	return

}
