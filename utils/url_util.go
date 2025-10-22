package utils

import (
	"net/url"
)

// Encode URL encodes a string using UTF-8 encoding
// Similar to Java's URLEncoder.encode()
func URLEncode(value string) string {
	return url.QueryEscape(value)
}

// Decode URL decodes a string using UTF-8 encoding
// Similar to Java's URLDecoder.decode()
func URLDecode(value string) (string, error) {
	return url.QueryUnescape(value)
}

// DecodeSafe is a version that returns empty string on error (like Java version)
func URLDecodeSafe(value string) string {
	decoded, err := url.QueryUnescape(value)
	if err != nil {
		return ""
	}
	return decoded
}
