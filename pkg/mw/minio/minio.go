package minio

import (
	"bytes"
	"context"
	"fmt"
	"github.com/joeyz1729/zero-tiktok/pkg/mw/ffmpeg"
	"github.com/joeyz1729/zero-tiktok/pkg/utils"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"mime/multipart"
	"net/url"
	"time"

	minio "github.com/minio/minio-go/v7"
)

var (
	ImageTypeSuffix      = ".png"
	MinioVideoBucketName = "tiktokvideo"
	MinioImgBucketName   = "tiktokimage"

	endpoint             = "localhost:9090"
	accessKeyId          = "minio"
	secretAccessKey      = "minio@123"
	useSSL          bool = false

	globalClient *minio.Client
)

func GetClient() *minio.Client {
	return globalClient
}

func MakeBucket(ctx context.Context, bucketName string) {
	exists, err := globalClient.BucketExists(ctx, bucketName)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !exists {
		err = globalClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
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
	info, err = globalClient.PutObject(ctx, bucketName, file.Filename, fileObj, file.Size, minio.PutObjectOptions{})
	return info, err
}

func GetObjURL(ctx context.Context, bucketName, filename string) (u *url.URL, err error) {
	exp := time.Hour * 24
	reqParams := make(url.Values)
	u, err = globalClient.PresignedGetObject(ctx, bucketName, filename, exp, reqParams)
	return u, err
}

func PutToBucketByBuf(ctx context.Context, bucketName, filename string, buf *bytes.Buffer) (info minio.UploadInfo, err error) {
	info, err = globalClient.PutObject(ctx, bucketName, filename, buf, int64(buf.Len()), minio.PutObjectOptions{})
	return info, err
}

func PutToBucketByFilePath(ctx context.Context, bucketName, filename, filepath string) (info minio.UploadInfo, err error) {
	// 是否需要加options：contentType？
	info, err = globalClient.FPutObject(ctx, bucketName, filename, filepath, minio.PutObjectOptions{})
	return info, err
}

func Init() {
	ctx := context.Background()
	var err error
	globalClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyId, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Println("minio client connect failed")
		log.Fatalln(err)
	}
	log.Println("minio client connect success")
	log.Printf("%#v\n", globalClient)

	MakeBucket(ctx, MinioVideoBucketName)
	MakeBucket(ctx, MinioImgBucketName)
}

func UploadVideo(ctx context.Context, userId int64, fileType string, data []byte) (string, string, error) {
	// 视频文件提交到minio
	timeNow := time.Now()
	filename := utils.NewFileName(userId, timeNow.Unix()) + fileType
	_, err := PutToBucketByBuf(
		ctx,
		MinioVideoBucketName,
		filename,
		bytes.NewBuffer(data),
	)
	if err != nil {
		return "", "", err
	}

	// 获取视频路径，并截取视频帧作为封面
	filepath, err := globalClient.PresignedGetObject(ctx, MinioVideoBucketName, filename, time.Minute*1, nil)
	if err != nil {
		return "", "", err
	}

	buf, err := ffmpeg.GetSnapshot(filepath.String())
	if err != nil {
		return "", "", err
	}

	// 将封面文件上传至minio
	_, err = PutToBucketByBuf(ctx, MinioImgBucketName, filename+".png", buf)
	if err != nil {
		return "", "", err
	}

	playURL := MinioVideoBucketName + "/" + filename + fileType
	coverURL := MinioImgBucketName + "/" + filename + ImageTypeSuffix
	return playURL, coverURL, nil
}
