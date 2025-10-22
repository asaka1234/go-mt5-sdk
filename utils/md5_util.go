package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// GetMD5 calculates MD5 hash of byte array and returns as 32-character hex string
func GetMD5(bytes []byte) string {
	hasher := md5.New()
	hasher.Write(bytes)
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

// GetMD5String calculates MD5 hash of a string with specified encoding
// Note: Go strings are UTF-8 by default, encoding parameter is kept for compatibility
func GetMD5String(value string, encoding string) string {
	// In Go, strings are UTF-8 by default, so we just convert to bytes
	// The encoding parameter is kept for API compatibility but not used
	return GetMD5([]byte(value))
}
