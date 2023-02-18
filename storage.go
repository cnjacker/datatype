package datatype

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

type StorageService interface {
	GetToken(config StorageConfig) string
	UploadFile(data []byte, token string, config StorageConfig) string
}

type StorageConfig struct {
	// 是否签名
	Signatured bool
	// 服务器URL
	ServerURL string
	// Bucket名称
	BucketName string
	// 访问Key
	AccessKey string
	// 签名Key
	SignatureKey string
	// 过期时间
	Expired time.Duration
	// 缓存
	Cache *cache.Cache
	// 存储服务
	Service StorageService
}

func (c StorageConfig) HasSignature() bool {
	if !c.Signatured {
		return false
	}

	if c.ServerURL == "" {
		return false
	}

	if c.BucketName == "" {
		return false
	}

	if c.AccessKey == "" {
		return false
	}

	if c.SignatureKey == "" {
		return false
	}

	return true
}

var StorageOptions = StorageConfig{
	Signatured: false,
	BucketName: "app",
	Expired:    5 * time.Minute,
	Cache:      cache.New(5*time.Minute, 10*time.Minute),
}

type Storage string

// GORM
func (s *Storage) Scan(value any) error {
	if v, ok := value.(string); ok {
		*s = Storage(Storage(v).BindSignature())
	}

	return nil
}

func (s Storage) Value() (driver.Value, error) {
	return s.UnBindSignature(), nil
}

// BindSignature
func (s Storage) BindSignature() string {
	path := s.UnBindSignature()

	if path != "" {
		path = "/" + strings.TrimPrefix(path, "/")

		if StorageOptions.HasSignature() {
			token := s.GetStorageToken()

			if token != "" {
				signature := HMacSha256([]byte(StorageOptions.SignatureKey), []byte(path))

				return fmt.Sprint(
					strings.TrimSuffix(StorageOptions.ServerURL, "/"),
					path,
					"?",
					token,
					",",
					signature,
				)
			}
		}

		return fmt.Sprint(
			strings.TrimSuffix(StorageOptions.ServerURL, "/"),
			path,
		)
	}

	return path
}

// UnBindSignature
func (s Storage) UnBindSignature() string {
	path := string(s)

	if path != "" {
		return strings.TrimPrefix(
			strings.Split(path, "?")[0],
			strings.TrimSuffix(StorageOptions.ServerURL, "/")+"/",
		)
	}

	return path
}

// UploadFile
func (s Storage) UploadFile(data []byte) string {
	token := s.GetStorageToken()

	if token != "" {
		return StorageOptions.Service.UploadFile(data, token, StorageOptions)
	}

	return ""
}

// GetStorageToken
func (s Storage) GetStorageToken() string {
	if v, ok := StorageOptions.Cache.Get("token"); ok {
		if token, ok := v.(string); ok {
			return token
		}
	}

	token := StorageOptions.Service.GetToken(StorageOptions)

	if token != "" {
		StorageOptions.Cache.Set("token", token, StorageOptions.Expired)
	}

	return ""
}

// String
func (s Storage) String() string {
	return s.UnBindSignature()
}

// HMacSha256
func HMacSha256(key []byte, data []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(data)

	return hex.EncodeToString(h.Sum(nil))
}
