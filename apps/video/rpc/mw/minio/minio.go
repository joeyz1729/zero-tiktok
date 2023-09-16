package minio

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"mime/multipart"
	"net/url"
	"time"

	minio "github.com/minio/minio-go/v7"
)

var (
	MinioVideoBucketName = "tiktokvideo"
	MinioImgBucketName   = "tiktokimage"

	endpoint             = "localhost:9090"
	accessKeyId          = "minio"
	secretAccessKey      = "minio@123"
	useSSL          bool = false

	Client *minio.Client
	err    error
)

func MakeBucket(ctx context.Context, bucketName string) {
	exists, err := Client.BucketExists(ctx, bucketName)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !exists {
		err = Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Successfully created mybucket %v\n", bucketName)
	}
}

func PutToBucket(ctx context.Context, bucketName string, file *multipart.FileHeader) (info minio.UploadInfo, err error) {
	fileObj, _ := file.Open()
	defer fileObj.Close()
	info, err = Client.PutObject(ctx, bucketName, file.Filename, fileObj, file.Size, minio.PutObjectOptions{})
	return info, err
}

func GetObjURL(ctx context.Context, bucketName, filename string) (u *url.URL, err error) {
	exp := time.Hour * 24
	reqParams := make(url.Values)
	u, err = Client.PresignedGetObject(ctx, bucketName, filename, exp, reqParams)
	return u, err
}

func PutToBucketByBuf(ctx context.Context, bucketName, filename string, buf *bytes.Buffer) (info minio.UploadInfo, err error) {
	info, err = Client.PutObject(ctx, bucketName, filename, buf, int64(buf.Len()), minio.PutObjectOptions{})
	return info, err
}

func PutToBucketByFilePath(ctx context.Context, bucketName, filename, filepath string) (info minio.UploadInfo, err error) {
	// 是否需要加options：contentType？
	info, err = Client.FPutObject(ctx, bucketName, filename, filepath, minio.PutObjectOptions{})
	return info, err
}

func Init() {
	ctx := context.Background()
	var err error
	Client, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyId, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Println("minio client connect failed")
		log.Fatalln(err)
	}
	log.Println("minio client connect success")
	log.Printf("%#v\n", Client)

	MakeBucket(ctx, MinioVideoBucketName)
	MakeBucket(ctx, MinioImgBucketName)
}
