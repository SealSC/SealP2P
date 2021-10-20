package grsa

import (
	"crypto/rsa"
	"os"
	"encoding/pem"
	"crypto/x509"
	"io/ioutil"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"errors"
)

func PubSha1(key *rsa.PrivateKey) string {
	if key == nil {
		return "pub_key_nil"
	}
	hash := sha1.New()
	publicKey := x509.MarshalPKCS1PublicKey(&key.PublicKey)
	hash.Write(publicKey)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func SaveFile(filePrefix string, key *rsa.PrivateKey) error {
	prif, err := os.OpenFile(filePrefix, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	defer prif.Close()
	pubf, err := os.OpenFile(filePrefix+".pub", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	defer pubf.Close()
	if err = pem.Encode(pubf, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&key.PublicKey),
	}); err != nil {
		return err
	}
	if err = pem.Encode(prif, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}); err != nil {
		return err
	}
	return nil
}

func LoadFile(filename string) (*rsa.PrivateKey, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	p, _ := pem.Decode(file)
	return x509.ParsePKCS1PrivateKey(p.Bytes)
}

func RandKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}
func LoadPubKey(data []byte) (*rsa.PublicKey, error) {
	return x509.ParsePKCS1PublicKey(data)
}

func Decode(cipherText []byte, k *rsa.PrivateKey) ([]byte, error) {
	if cipherText == nil {
		return nil, errors.New("cipherText is nil")
	}
	return rsa.DecryptPKCS1v15(rand.Reader, k, cipherText)
}
func Encode(plainText []byte, k *rsa.PublicKey) ([]byte, error) {
	if plainText == nil {
		return nil, errors.New("plainText is nil")
	}
	return rsa.EncryptPKCS1v15(rand.Reader, k, plainText)
}
