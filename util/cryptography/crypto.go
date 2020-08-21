package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func AesEncryptCBCWithIv(origData []byte, key []byte, iv []byte) (encrypted []byte) {
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()                 // 获取秘钥块的长度
	origData = pkcs5Padding(origData, blockSize)   // 补全码
	blockMode := cipher.NewCBCEncrypter(block, iv) // 加密模式
	encrypted = make([]byte, len(origData))        // 创建数组
	blockMode.CryptBlocks(encrypted, origData)     // 加密
	return encrypted
}

func RSAEncryptV2(origData []byte, publicKey *rsa.PublicKey) []byte {
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, origData)
	if err != nil {
		fmt.Println("rsa.EncryptPKCS1v15:", err)
		return encrypted
	}
	return encrypted
}

func ParsePublicKey(publicKey []byte) (*rsa.PublicKey, error) {
	pemBlock, _ := pem.Decode(publicKey)
	if pemBlock == nil {
		fmt.Println("pem.Decode error")
		return nil, fmt.Errorf("pem.Decode error")
	}
	pubKey, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if err != nil {
		fmt.Println("x509.ParsePKCS1PublicKey:", err)
		return nil, err
	}
	return pubKey.(*rsa.PublicKey), nil
}
