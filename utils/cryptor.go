package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
)

var (
	SignSuffix = "3478cbbc33f84bd00d75d7dfa69e0daa"
	MoGuKEY    = "23DbtQHR2UMbH6mJ"
	CXKEY      = []byte("u2oh6Vu^HWe4_AES")
)

// PKCS7Padding 对数据进行填充，使其长度对齐到块大小的倍数
func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// AESCBCEncrypt 使用 AES CBC 模式加密
func AESCBCEncrypt(plaintext []byte) (string, error) {
	// 创建 AES 块加密器
	block, err := aes.NewCipher(CXKEY)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	// AES 块大小
	blockSize := block.BlockSize()

	// 填充明文
	plaintext = PKCS7Padding(plaintext, blockSize)

	// 创建 CBC 模式加密器
	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, CXKEY[:blockSize])
	mode.CryptBlocks(ciphertext, plaintext)

	// 使用 Base64 编码返回加密结果
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// ===================================GongxueYun
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

// AESECBPKCS5Padding 封装的 AES-ECB 加密器
type AESECBPKCS5Padding struct {
	key       []byte
	outFormat string
}

// NewAESECBPKCS5Padding 创建一个 AES-ECB 加密器实例
func NewAESECBPKCS5Padding(key string, outFormat string) (*AESECBPKCS5Padding, error) {
	if len(key) != 16 {
		return nil, errors.New("密钥长度必须为 16 字节")
	}

	if outFormat != "hex" && outFormat != "base64" {
		return nil, errors.New("输出格式必须是 'hex' 或 'base64'")
	}

	return &AESECBPKCS5Padding{
		key:       []byte(key),
		outFormat: outFormat,
	}, nil
}

// Encrypt 实现 AES-ECB 加密
func (a *AESECBPKCS5Padding) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", fmt.Errorf("创建 AES 加密器失败: %v", err)
	}

	// PKCS5 填充
	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	paddedText := append([]byte(plaintext), bytes.Repeat([]byte{byte(padding)}, padding)...)

	// ECB 加密（直接分块加密）
	ciphertext := make([]byte, len(paddedText))
	for i := 0; i < len(paddedText); i += aes.BlockSize {
		block.Encrypt(ciphertext[i:i+aes.BlockSize], paddedText[i:i+aes.BlockSize])
	}

	// 根据指定格式输出
	switch a.outFormat {
	case "hex":
		return hex.EncodeToString(ciphertext), nil
	case "base64":
		return base64.StdEncoding.EncodeToString(ciphertext), nil
	default:
		return "", errors.New("未知的输出格式")
	}
}

// Decrypt 实现 AES-ECB 解密
func (a *AESECBPKCS5Padding) Decrypt(ciphertext string) (string, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", fmt.Errorf("创建 AES 解密器失败: %v", err)
	}

	// 解码密文
	var decodedCiphertext []byte
	switch a.outFormat {
	case "hex":
		decodedCiphertext, err = hex.DecodeString(ciphertext)
		if err != nil {
			return "", fmt.Errorf("解码 hex 密文失败: %v", err)
		}
	case "base64":
		decodedCiphertext, err = base64.StdEncoding.DecodeString(ciphertext)
		if err != nil {
			return "", fmt.Errorf("解码 base64 密文失败: %v", err)
		}
	default:
		return "", errors.New("未知的输出格式")
	}

	// ECB 解密（直接分块解密）
	plaintext := make([]byte, len(decodedCiphertext))
	for i := 0; i < len(decodedCiphertext); i += aes.BlockSize {
		block.Decrypt(plaintext[i:i+aes.BlockSize], decodedCiphertext[i:i+aes.BlockSize])
	}

	// 去除 PKCS5 填充
	plaintext, err = pkcs5Unpadding(plaintext)
	if err != nil {
		return "", fmt.Errorf("去除 PKCS5 填充失败: %v", err)
	}

	return string(plaintext), nil
}

// pkcs5Unpadding 去除 PKCS5 填充
func pkcs5Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("数据长度为 0，无法去除填充")
	}
	padding := int(data[length-1])
	if padding > aes.BlockSize || padding > length {
		return nil, errors.New("填充长度不合法")
	}
	for _, pad := range data[length-padding:] {
		if pad != byte(padding) {
			return nil, errors.New("填充内容不合法")
		}
	}
	return data[:length-padding], nil
}
