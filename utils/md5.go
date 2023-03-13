package utils

import (
	"crypto/md5"
	"fmt"
)

func MD5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	sum := hash.Sum(nil)
	return fmt.Sprintf("%x\n", sum)
}

func MD5Salt(s string, salt string) string {
	return MD5(s + salt)
}
