package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

// Sign signs message with private key using RSA-SHA1.
func Sign(privateKey *rsa.PrivateKey, message []byte) (signature []byte, err error) {
	hasher := sha1.New()
	hasher.Write(message)
	return rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, hasher.Sum(nil))
}

// SignBase64 just like Sign except returning base64 encoded string instead of raw bytes.
func SignBase64(privateKey *rsa.PrivateKey, message []byte) (signature string, err error) {
	hasher := sha1.New()
	hasher.Write(message)
	var rawSig []byte
	if rawSig, err = rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, hasher.Sum(nil)); err != nil {
		return
	} else {
		return base64.StdEncoding.EncodeToString(rawSig), nil
	}
}

// Verify checks message-signature pair using given public key.
// A valid signature is indicated by returning a nil error.
func Verify(publicKey *rsa.PublicKey, message []byte, signature []byte) (err error) {
	hasher := sha1.New()
	hasher.Write(message)
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA1, hasher.Sum(nil), signature)

}

// VerifyBase64 just like Verify except taking b64encoded signature instead of raw bytes.
func VerifyBase64(publicKey *rsa.PublicKey, message []byte, b64Signature string) (err error) {
	hasher := sha1.New()
	hasher.Write(message)

	signature, err := base64.StdEncoding.DecodeString(b64Signature)
	if err != nil {
		return
	}
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA1, hasher.Sum(nil), signature)
}

// ParsePrivateKey will load private key from given string in PKCS1 format
func ParsePrivateKey(keyStr []byte) (privateKey *rsa.PrivateKey, err error) {
	block, _ := pem.Decode(keyStr)
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// ParsePublicKey will load private key from given string in PKIX format
func ParsePublicKey(keyStr []byte) (publicKey *rsa.PublicKey, err error) {
	block, _ := pem.Decode(keyStr)
	publicInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if pk, ok := publicInterface.(*rsa.PublicKey); !ok {
		return nil, fmt.Errorf("parse failed: not a valid rsa public key")
	} else {
		return pk, nil
	}
}

// LoadPrivateKey will load private key from file
func LoadPrivateKey(path string) (privateKey *rsa.PrivateKey, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	return ParsePrivateKey(data)
}

// LoadPublicKey will load public key from file
func LoadPublicKey(path string) (public *rsa.PublicKey, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(path); err != nil {
		return
	} else {
		return ParsePublicKey(data)
	}
}