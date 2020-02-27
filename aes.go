/////////////////////////////////////////////////////////////////////////////////
//
// Copyright (C) 2019-2020, Unectio Inc, All Right Reserved.
//
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
/////////////////////////////////////////////////////////////////////////////////

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
