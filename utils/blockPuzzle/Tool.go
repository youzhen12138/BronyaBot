package blockPuzzle

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/gtank/cryptopasta"
)

// CreateSign 生成签名
func CreateSign(args ...string) string {
	// 拼接字符串
	signStr := ""
	for _, arg := range args {
		signStr += arg
	}
	// 添加预定义密钥
	signStr += "3478cbbc33f84bd00d75d7dfa69e0daa"

	// 生成MD5签名
	hash := md5.Sum([]byte(signStr))
	return hex.EncodeToString(hash[:])
}

// Encrypt 使用 cryptopasta 加密
func Encrypt(plaintext string, key []byte) (string, error) {
	if len(key) != 32 {
		return "", errors.New("密钥长度必须为 32 字节")
	}

	ciphertext, err := cryptopasta.Encrypt([]byte(plaintext), (*[32]byte)(key))
	if err != nil {
		return "", fmt.Errorf("加密失败: %v", err)
	}

	return hex.EncodeToString(ciphertext), nil
}

// Decrypt 使用 cryptopasta 解密
func Decrypt(ciphertext string, key []byte) (string, error) {
	if len(key) != 32 {
		return "", errors.New("密钥长度必须为 32 字节")
	}

	cipherBytes, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("解码密文失败: %v", err)
	}

	plaintext, err := cryptopasta.Decrypt(cipherBytes, (*[32]byte)(key))
	if err != nil {
		return "", fmt.Errorf("解密失败: %v", err)
	}

	return string(plaintext), nil
}
