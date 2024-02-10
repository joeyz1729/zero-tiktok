package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"net/url"
	"testing"
)

func TestGetObjURL(t *testing.T) {
	Init()
	ctx := context.Background()
	var u = new(url.URL)
	var err error
	u, err = GetObjURL(ctx, MinioVideoBucketName, "1000.1676403991.mp4")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(u.String())
	}
}

func TestBuckMake(t *testing.T) {
	ctx := context.Background()
	exists, err := globalClient.BucketExists(ctx, MinioVideoBucketName)
	if err != nil {
		fmt.Println(err)
		return
	}
	if exists {
		fmt.Printf("%v found\n", MinioVideoBucketName)
	} else {
		err = globalClient.MakeBucket(ctx, MinioVideoBucketName, minio.MakeBucketOptions{})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Successfully created mybucket %v\n", MinioVideoBucketName)
	}
}

func TestBucketExist(t *testing.T) {
	ctx := context.Background()
	exists, err := globalClient.BucketExists(ctx, MinioVideoBucketName)
	if err != nil {
		fmt.Println(err)
		return
	}
	if exists {
		fmt.Printf("%v found\n", MinioVideoBucketName)
	} else {
		fmt.Println("not found!")
	}
}

func TestPutToBucketByFilePath(t *testing.T) {
	Init()
	ctx := context.Background()
	info, err := PutToBucketByFilePath(ctx, MinioVideoBucketName, "miniotest.mp4", "/Users/zouyi/Movies/miniotest.mp4")
	if err != nil {
		fmt.Printf("put to bucket by file path failed, err: %v\n", err)

	} else {
		fmt.Printf("put to bucket by file path success, info: %v\n", info)
	}
}
