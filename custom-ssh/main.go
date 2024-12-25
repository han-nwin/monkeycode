package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Println("Error generating Ed25519 key:", err)
		return
	}

	fmt.Println("Public Key:", publicKey)
	fmt.Println("Private Key:", privateKey)

    //Encode private key in PKCS#8 format
    privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
    if err != nil {
        fmt.Println("Error encoding private key:", err)
        return
    }

    //Encode publicKey in PKIX format
    publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
    if err != nil {
        fmt.Println("Error encoding public key", err)
        return
    }

    //Create PEM block for private key
    privateKeyPEM := pem.EncodeToMemory(&pem.Block{
        Type: "ED25519 PRIVATE KEY",
        Bytes: privateKeyBytes,
    })

    //Create PEM block for public key
    publicKeyPEM := pem.EncodeToMemory(&pem.Block{
        Type: "ED25519 PUBLIC KEY",
        Bytes: publicKeyBytes,
    })
    
    //Save private key PEM to file
    if err := os.WriteFile(".ssh/id_ed25519", privateKeyPEM, 0600); err != nil {
        fmt.Println("Error saving private key", err)
        return
    }

    //Save public key PEM to file
    if err := os.WriteFile(".ssh/id_ed25519.pub", publicKeyPEM, 0600); err != nil {
        fmt.Println("Error saving public key", err)
        return
    }
}
