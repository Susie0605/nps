package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"github.com/google/uuid"
	"math/rand"
	"strings"
	"time"
)

// en
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// de
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	err, origData = PKCS5UnPadding(origData)
	return origData, err
}

// Completion when the length is insufficient
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// Remove excess
func PKCS5UnPadding(origData []byte) (error, []byte) {
	length := len(origData)
	unpadding := int(origData[length-1])
	if (length - unpadding) < 0 {
		return errors.New("len error"), nil
	}
	return nil, origData[:(length - unpadding)]
}

// Generate 32-bit MD5 strings
func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// Generating Random Verification Key
func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GetVkey() string {
	// 生成UUID
	u, _ := uuid.NewRandom()
	// 将UUID转换为字符串
	uuidStr := u.String()
	uuidStr = strings.ReplaceAll(uuidStr, "-", "")
	// 截取前10位
	return uuidStr[:10]
}

func Base64Decoding(encodedString string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		// 处理解码错误
		return "", err
	}

	// 将解码后的字节数组转换为字符串
	decodedString := string(decodedBytes)

	// 如果startCmd开头不包含"nps "，则报错
	if len(decodedString) < 4 || !(decodedString[0:4] == "nps ") {
		return "", errors.New("快捷启动命令错误，请检查")
	}

	return decodedString[4:], nil

}
