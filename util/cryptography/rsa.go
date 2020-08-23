package cryptography

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
)


func RSAEncrypt(origData []byte, modulus string, exponent int64) string {
	bigOrigData := big.NewInt(0).SetBytes(origData)
	bigModulus, _ := big.NewInt(0).SetString(modulus, 16)
	bigRs := bigOrigData.Exp(bigOrigData, big.NewInt(exponent), bigModulus)
	return fmt.Sprintf("%0256x", bigRs)
}

// BytesToPublicKey bytes to public key
func BytesToPublicKey(pub []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, err
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		log.Fatal("converting RSA public key failed")
	}
	return key, nil
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, pub, msg)
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
