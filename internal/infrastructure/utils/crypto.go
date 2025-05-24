package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"github.com/gogf/gf/v2/frame/g"
)

// AESEncrypt AES加密
func AESEncrypt(plaintext []byte) (string, error) {
	key := []byte(g.Cfg().MustGet("aes.key").String())
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 填充
	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	plaintext = append(plaintext, padtext...)

	// 加密
	ciphertext := make([]byte, len(plaintext))
	blockMode := cipher.NewCBCEncrypter(block, key[:aes.BlockSize])
	blockMode.CryptBlocks(ciphertext, plaintext)

	// Base64编码
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AESDecrypt AES解密
func AESDecrypt(ciphertext string) ([]byte, error) {
	key := []byte(g.Cfg().MustGet("aes.key").String())
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Base64解码
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	// 解密
	plaintext := make([]byte, len(ciphertextBytes))
	blockMode := cipher.NewCBCDecrypter(block, key[:aes.BlockSize])
	blockMode.CryptBlocks(plaintext, ciphertextBytes)

	// 去除填充
	padding := int(plaintext[len(plaintext)-1])
	return plaintext[:len(plaintext)-padding], nil
} 