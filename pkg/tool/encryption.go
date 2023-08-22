package tool

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

var salt string

func Encrypt(op string) (ep string) {
	h := md5.New()
	h.Write([]byte(salt))
	ep = hex.EncodeToString(h.Sum([]byte(op)))
	return
}

func Md5ByBytes(b []byte) string {
	return fmt.Sprintf("%x", md5.Sum(b))
}

func NewSalt(s string) {
	salt = s
}
