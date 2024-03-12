package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strings"
)

func Base64Encode(src []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(src))
}

func Base64Decode(src []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(src))
}

func CheckBase64Password(pwd, pwdFromDB string) bool {
	decode, _ := Base64Decode([]byte(pwdFromDB))
	return strings.EqualFold(pwd, string(decode))
}

func Sha256Encode(src []byte) string {
	sum256 := sha256.Sum256(src)
	return hex.EncodeToString(sum256[0:])
}

func CheckSha256Password(pwd, pwdFromDB string) bool {
	encode := Sha256Encode([]byte(pwd))
	return encode == pwdFromDB
}
