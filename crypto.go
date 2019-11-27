package util

import (
	"crypto/rand"
	"crypto/hmac"
	"encoding/hex"
	"crypto/sha256"
)

func Sha256(s string) string {
	hash := sha256.New()
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}

var rsRunes = []rune("_1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomBytes(l int) ([]byte, error) {
	ri := make([]byte, l)

	_, err := rand.Read(ri)
	if err != nil {
		return nil, err
	}

	return ri, nil
}

func RandomString(l int) (string, error) {
	ri := make([]byte, l)

	_, err := rand.Read(ri)
	if err != nil {
		return "", err
	}

	rs := make([]rune, l)
	for i, j := range ri {
		rs[i] = rsRunes[int(j) % len(rsRunes)]
	}

	return string(rs), nil
}


func HashHMAC(sec []byte, data ...[]byte) string {
	hash := hmac.New(sha256.New, sec)
	for _, dat := range data {
		hash.Write(dat)
	}
	return hex.EncodeToString(hash.Sum(nil))
}
