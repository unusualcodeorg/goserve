package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func generateRSAKeyPair() (*rsa.PrivateKey, error) {
	// Generate a new RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	// Validate private key
	if err := privateKey.Validate(); err != nil {
		return nil, err
	}

	return privateKey, nil
}

func savePrivateKeyToFile(privateKey *rsa.PrivateKey, filename string) error {
	// Encode private key to PEM format
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	// Create a new file for the private key
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write private key to the file
	if err := pem.Encode(file, pemBlock); err != nil {
		return err
	}

	return nil
}

func savePublicKeyToFile(publicKey *rsa.PublicKey, filename string) error {
	// Marshal public key to DER format
	derBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}

	// Encode public key to PEM format
	pemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derBytes,
	}

	// Create a new file for the public key
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write public key to the file
	if err := pem.Encode(file, pemBlock); err != nil {
		return err
	}

	return nil
}

func main() {
	// Generate RSA key pair
	privateKey, err := generateRSAKeyPair()
	if err != nil {
		fmt.Println("Error generating RSA key pair:", err)
		return
	}

	// Save private key to file
	err = savePrivateKeyToFile(privateKey, "keys/private.pem")
	if err != nil {
		fmt.Println("Error saving private key:", err)
		return
	}

	// Save public key to file
	err = savePublicKeyToFile(&privateKey.PublicKey, "keys/public.pem")
	if err != nil {
		fmt.Println("Error saving public key:", err)
		return
	}

	fmt.Println("RSA key pair generated and saved to keys/private.pem and keys/public.pem")
}
