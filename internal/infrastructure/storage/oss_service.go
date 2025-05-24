package storage

import (
	"context"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gogf/gf/v2/frame/g"
	"path"
	"time"
)

type ossService struct {
	client     *oss.Client
	bucket     *oss.Bucket
	bucketName string
}

// NewOSSService 创建OSS存储服务实例
func NewOSSService() (*ossService, error) {
	endpoint := g.Cfg().MustGet(context.Background(), "oss.endpoint").String()
	accessKeyID := g.Cfg().MustGet(context.Background(), "oss.accessKeyId").String()
	accessKeySecret := g.Cfg().MustGet(context.Background(), "oss.accessKeySecret").String()
	bucketName := g.Cfg().MustGet(context.Background(), "oss.bucketName").String()

	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, err
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	return &ossService{
		client:     client,
		bucket:     bucket,
		bucketName: bucketName,
	}, nil
}

// UploadFile 上传文件到OSS
func (s *ossService) UploadFile(objectKey string, filePath string) (string, error) {
	err := s.bucket.PutObjectFromFile(objectKey, filePath)
	if err != nil {
		return "", err
	}

	// 生成文件访问URL
	url := s.bucket.SignURL(objectKey, oss.HTTPGet, time.Hour*24*365)
	return url, nil
}

// UploadContent 上传内容到OSS
func (s *ossService) UploadContent(objectKey string, content []byte) (string, error) {
	err := s.bucket.PutObject(objectKey, content)
	if err != nil {
		return "", err
	}

	// 生成文件访问URL
	url := s.bucket.SignURL(objectKey, oss.HTTPGet, time.Hour*24*365)
	return url, nil
}

// DeleteFile 从OSS删除文件
func (s *ossService) DeleteFile(objectKey string) error {
	return s.bucket.DeleteObject(objectKey)
}

// GetFileURL 获取文件访问URL
func (s *ossService) GetFileURL(objectKey string) string {
	return s.bucket.SignURL(objectKey, oss.HTTPGet, time.Hour*24*365)
}

// GenerateObjectKey 生成对象键
func (s *ossService) GenerateObjectKey(prefix string, filename string) string {
	ext := path.Ext(filename)
	return prefix + "/" + time.Now().Format("2006/01/02") + "/" + filename + ext
} 