package tools

import "encoding/base64"

func Base64ToString(input []byte) string {
	return base64.StdEncoding.EncodeToString(input)
}

func Base64ToBytes(str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(str)
}
