package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"pinche/config"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// UploadBizType defines the business type for upload
type UploadBizType string

const (
	BizTypeImage  UploadBizType = "images" // chat images
	BizTypeAvatar UploadBizType = "avatar" // user avatar
	BizTypeTrip   UploadBizType = "trip"   // trip images
	BizTypeVoice  UploadBizType = "voices" // voice messages
	BizTypeVideo  UploadBizType = "videos" // video messages
	BizTypeThumb  UploadBizType = "thumbs" // thumbnails for images/videos
)

// allowed file extensions for each biz type
var allowedImageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
}

var allowedVoiceExts = map[string]bool{
	".webm": true, ".mp3": true, ".wav": true, ".m4a": true, ".ogg": true, ".aac": true,
}

var allowedVideoExts = map[string]bool{
	".mp4": true, ".webm": true, ".mov": true, ".m4v": true,
}

// IsValidBizType checks if the biz type is valid
func IsValidBizType(bizType UploadBizType) bool {
	return bizType == BizTypeImage || bizType == BizTypeAvatar || bizType == BizTypeTrip ||
		bizType == BizTypeVoice || bizType == BizTypeVideo || bizType == BizTypeThumb
}

type UploadService struct {
	cosClient *cos.Client
	bucket    string
	region    string
}

func NewUploadService(cfg *config.Config) *UploadService {
	// COS configuration from config
	bucketURL, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", cfg.COS.Bucket, cfg.COS.Region))
	serviceURL, _ := url.Parse(fmt.Sprintf("https://cos.%s.myqcloud.com", cfg.COS.Region))

	client := cos.NewClient(&cos.BaseURL{BucketURL: bucketURL, ServiceURL: serviceURL}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.COS.SecretID,
			SecretKey: cfg.COS.SecretKey,
		},
	})

	return &UploadService{
		cosClient: client,
		bucket:    cfg.COS.Bucket,
		region:    cfg.COS.Region,
	}
}

// Upload handles file upload to COS based on biz type
// Returns object key for private files, or public URL for avatar/trip
func (s *UploadService) Upload(file *multipart.FileHeader, bizType UploadBizType) (string, error) {
	// validate biz type
	if !IsValidBizType(bizType) {
		bizType = BizTypeImage
	}

	// validate file type based on biz type
	ext := strings.ToLower(filepath.Ext(file.Filename))

	if bizType == BizTypeVoice {
		if !allowedVoiceExts[ext] {
			return "", errors.New("不支持的音频格式，仅支持webm/mp3/wav/m4a/ogg/aac")
		}
	} else if bizType == BizTypeVideo {
		if !allowedVideoExts[ext] {
			return "", errors.New("不支持的视频格式，仅支持mp4/webm/mov/m4v")
		}
	} else {
		if !allowedImageExts[ext] {
			return "", errors.New("不支持的图片格式，仅支持jpg/jpeg/png/gif/webp")
		}
	}

	// validate file size (max 100MB)
	if file.Size > 100*1024*1024 {
		return "", errors.New("图片大小不能超过100MB")
	}

	// open uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	// read file content for hash calculation
	content, err := io.ReadAll(src)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}

	// generate unique filename using file content md5
	hash := md5.Sum(content)
	hashStr := hex.EncodeToString(hash[:])

	// generate COS object key: {bizType}/hash.ext
	filename := fmt.Sprintf("%s%s", hashStr, ext)
	objectKey := fmt.Sprintf("%s/%s", bizType, filename)

	// upload to COS
	ctx := context.Background()

	// trip images should be public-read
	if bizType == BizTypeTrip {
		opt := &cos.ObjectPutOptions{
			ACLHeaderOptions: &cos.ACLHeaderOptions{
				XCosACL: "public-read",
			},
		}
		_, err = s.cosClient.Object.Put(ctx, objectKey, bytes.NewReader(content), opt)
		if err != nil {
			return "", fmt.Errorf("上传到COS失败: %w", err)
		}
		// return public URL for trip images
		publicURL := fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s", s.bucket, s.region, objectKey)
		return publicURL, nil
	}

	// set Content-Type for voice files to ensure proper playback
	var putOpt *cos.ObjectPutOptions
	if bizType == BizTypeVoice {
		contentType := "audio/webm"
		switch ext {
		case ".mp3":
			contentType = "audio/mpeg"
		case ".wav":
			contentType = "audio/wav"
		case ".m4a":
			contentType = "audio/mp4"
		case ".ogg":
			contentType = "audio/ogg"
		case ".aac":
			contentType = "audio/aac"
		}
		putOpt = &cos.ObjectPutOptions{
			ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
				ContentType: contentType,
			},
		}
	} else if bizType == BizTypeVideo {
		contentType := "video/mp4"
		switch ext {
		case ".webm":
			contentType = "video/webm"
		case ".mov":
			contentType = "video/quicktime"
		case ".m4v":
			contentType = "video/mp4"
		}
		putOpt = &cos.ObjectPutOptions{
			ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
				ContentType: contentType,
			},
		}
	}

	_, err = s.cosClient.Object.Put(ctx, objectKey, bytes.NewReader(content), putOpt)
	if err != nil {
		return "", fmt.Errorf("上传到COS失败: %w", err)
	}

	// return object key (not full URL)
	return objectKey, nil
}

// UploadAvatar uploads avatar with openID as filename, returns the public URL
func (s *UploadService) UploadAvatar(file *multipart.FileHeader, openID string) (string, error) {
	// validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
	}
	if !allowedExts[ext] {
		return "", errors.New("不支持的图片格式，仅支持jpg/jpeg/png/gif/webp")
	}

	// validate file size (max 2MB for avatar)
	if file.Size > 10*1024*1024 {
		return "", errors.New("头像图片大小不能超过10MB")
	}

	// open uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	// read file content
	content, err := io.ReadAll(src)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}

	// generate COS object key: avatar/{openID}.ext
	objectKey := fmt.Sprintf("avatar/%s%s", openID, ext)

	// upload to COS with public-read ACL
	ctx := context.Background()
	opt := &cos.ObjectPutOptions{
		ACLHeaderOptions: &cos.ACLHeaderOptions{
			XCosACL: "public-read",
		},
	}
	_, err = s.cosClient.Object.Put(ctx, objectKey, bytes.NewReader(content), opt)
	if err != nil {
		return "", fmt.Errorf("上传到COS失败: %w", err)
	}

	// return public URL (no signature needed)
	avatarURL := fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s", s.bucket, s.region, objectKey)
	return avatarURL, nil
}

// GetSignedURL generates a presigned URL for the given object key
func (s *UploadService) GetSignedURL(objectKey string, expireSeconds int) (string, error) {
	if objectKey == "" {
		return "", errors.New("资源key不能为空")
	}

	// default expire time is 1 hour
	if expireSeconds <= 0 {
		expireSeconds = 3600
	}
	// max expire time is 1 day
	if expireSeconds > 1*24*3600 {
		expireSeconds = 1 * 24 * 3600
	}

	ctx := context.Background()
	expireDuration := time.Duration(expireSeconds) * time.Second

	signedURL, err := s.cosClient.Object.GetPresignedURL(ctx, http.MethodGet, objectKey, s.cosClient.GetCredential().SecretID, s.cosClient.GetCredential().SecretKey, expireDuration, nil)
	if err != nil {
		return "", fmt.Errorf("生成签名URL失败: %w", err)
	}

	return signedURL.String(), nil
}
