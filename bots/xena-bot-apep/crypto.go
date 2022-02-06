package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

// generatePrivateKey creates a RSA private key.
func generatePrivateKey() *rsa.PrivateKey {
	secret, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		fmt.Println(err)
	}
	return secret
}

// publicKeyToPEM converts RSA public key object to PEM encoded string.
func publicKeyToPEM(publicKey *rsa.PublicKey) string {
	spkiDER, _ := x509.MarshalPKIXPublicKey(publicKey)
	spkiPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: spkiDER,
		},
	)
	return string(spkiPEM)
}

// privateKeyToPEM converts RSA private key object to PEM encoded string.
func privateKeyToPEM(privateKey *rsa.PrivateKey) string {
	spkiDER, _ := x509.MarshalPKCS8PrivateKey(privateKey)
	spkiPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: spkiDER,
		},
	)
	return string(spkiPEM)
}

// importPEMPrivateKey converts a PEM encoded private key into the rsa.PublicKey object.
func importPEMPrivateKey(spkiPEM string) *rsa.PrivateKey {
	body, _ := pem.Decode([]byte(spkiPEM))
	privateKey, _ := x509.ParsePKCS1PrivateKey(body.Bytes)
	return privateKey
}

// importPEMPublicKey converts a PEM encoded public key into the rsa.PublicKey object.
func importPEMPublicKey(spkiPEM string) *rsa.PublicKey {
	body, _ := pem.Decode([]byte(spkiPEM))
	publicKey, _ := x509.ParsePKIXPublicKey(body.Bytes)
	if publicKey, ok := publicKey.(*rsa.PublicKey); ok {
		return publicKey
	}
	return nil
}

func encryptRSAOAEP(secretMessage string, key rsa.PublicKey) string {
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, &key, []byte(secretMessage), label)
	if err != nil {
		fmt.Println(err)
	}
	return base64.StdEncoding.EncodeToString(ciphertext)
}

func decryptRSAOAEP(cipherText string, privKey rsa.PrivateKey) string {
	ct, _ := base64.StdEncoding.DecodeString(cipherText)
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, &privKey, ct, label)
	if err != nil {
		fmt.Println(err)
	}
	return string(plaintext)
}