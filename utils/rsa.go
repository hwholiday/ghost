package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"hash"
)

func RsaNewKey(bits ...int) ([]byte, []byte, error) {
	var bit int
	if len(bits) > 0 {
		bit = bits[0]
	} else {
		bit = 1024
	}
	privateKey, err := rsa.GenerateKey(rand.Reader, bit)
	if err != nil {
		return nil, nil, err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	private := pem.EncodeToMemory(block)
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		panic(err)
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	public := pem.EncodeToMemory(block)
	return private, public, nil
}

// RsaEncrypt RSA加密
// plainText 要加密的数据
// publicKeyBytes 公钥
func RsaEncrypt(plainText, publicKeyBytes []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKeyBytes)
	publicKeyAny, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	publicKey := publicKeyAny.(*rsa.PublicKey)
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err != nil {
		return nil, err
	}
	return cipherText, nil
}

// RsaDecrypt RSA解密
// cipherText 需要解密的byte数据
// privateKeyBytes 私钥
func RsaDecrypt(cipherText, privateKeyBytes []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKeyBytes)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	return plainText, err
}

// RsaSign RSA签名
// signContent 签名的byte数据
// privateKeyBytes 私钥
// h Hash算法，默认为SHA256
func RsaSign(privateKeyBytes, signContent []byte, h ...crypto.Hash) ([]byte, error) {
	var shaNew hash.Hash
	var cryptoNew crypto.Hash
	if len(h) > 0 {
		shaNew = h[0].New()
		cryptoNew = h[0]
	} else {
		shaNew = sha256.New()
		cryptoNew = crypto.SHA256
	}
	shaNew.Write([]byte(signContent))
	hashed := shaNew.Sum(nil)
	block, _ := pem.Decode(privateKeyBytes)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.SignPKCS1v15(rand.Reader, privateKey, cryptoNew, hashed)
}

// RsaVerifySign 签名校验
// signContent 签名的byte数据
// publicKeyBytes 公钥
// sign 签名
// h Hash算法，默认为SHA256
func RsaVerifySign(publicKeyBytes, signContent, sign []byte, h ...crypto.Hash) error {
	var shaNew hash.Hash
	var cryptoNew crypto.Hash
	if len(h) > 0 {
		shaNew = h[0].New()
		cryptoNew = h[0]
	} else {
		shaNew = sha256.New()
		cryptoNew = crypto.SHA256
	}
	shaNew.Write(signContent)
	block, _ := pem.Decode(publicKeyBytes)
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(publicKey.(*rsa.PublicKey), cryptoNew, shaNew.Sum(nil), sign)
}
