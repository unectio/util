package util

import (
	"io"
	"fmt"
	"bytes"
	"errors"
	"crypto/aes"
	"crypto/rand"
	"crypto/cipher"
)

const (
	keyTrim	= 16
)

func pad(buf []byte) []byte {
	padsz := aes.BlockSize - len(buf)%aes.BlockSize
	return append(buf, bytes.Repeat([]byte{byte(padsz)}, padsz)...)
}

func unpad(buf []byte) []byte {
	l := len(buf)
	padsz := int(buf[l - 1])
	if padsz > l {
		return nil
	}
	return buf[:(l - padsz)]
}

func Encrypt(key []byte, ptext []byte) ([]byte, error) {
	pmsg := pad(ptext)

	if len(key) > keyTrim {
		key = key[:keyTrim]
	}

	cip, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("Error setting up AES: %s", err.Error())
	}

	cmsg := make([]byte, aes.BlockSize + len(pmsg))
	nonce := cmsg[:aes.BlockSize]

	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, fmt.Errorf("Error generating nonce: %s", err.Error())
	}

	e := cipher.NewCFBEncrypter(cip, nonce)
	e.XORKeyStream(cmsg[aes.BlockSize:], []byte(pmsg))

	return cmsg, nil
}

func Decrypt(key []byte, msg []byte) ([]byte, error) {
	if len(msg) % aes.BlockSize != 0 {
		return nil, errors.New("Cipher text trimmed")
	}

	if len(key) > keyTrim {
		key = key[:keyTrim]
	}

	aesc, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("Error setting up AES: %s", err.Error())
	}

	nonce := msg[:aes.BlockSize]
	cmsg := msg[aes.BlockSize:]

	d := cipher.NewCFBDecrypter(aesc, nonce)
	d.XORKeyStream(cmsg, cmsg)

	xmsg := unpad(cmsg)
	if xmsg == nil {
		return nil, errors.New("Decoded message unpad error")
	}

	return xmsg, nil
}
