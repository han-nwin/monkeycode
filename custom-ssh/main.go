package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"

	"golang.org/x/crypto/ssh"
)

func main() {
    //==== Key Pair Generation ====

    //Step 1: Generate key pair
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("Error generating Ed25519 key:", err)
	}

	fmt.Println("Public Key:", publicKey)
	fmt.Println("Private Key:", privateKey)

    //Step 2: Selialize private key in OpenSSH format
	comment := "user@hostname" // Optional comment for the private key
	privateKeyBytes, err := ssh.MarshalPrivateKey(privateKey, comment)
	if err != nil {
		log.Fatalf("Error marshaling private key: %v", err)
	}
    //Create PEM block for private key
    privateKeyPEM := pem.EncodeToMemory(privateKeyBytes)

    //Step 3: Save private key PEM to file
    if err := os.WriteFile(".ssh/id_ed25519", privateKeyPEM, 0600); err != nil {
        log.Fatalf("Error saving private key", err)
    }


	// Step 4: Convert public key to OpenSSH format
	sshPublicKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		log.Fatalf("Error converting public key to OpenSSH format: %v", err)
	}
    sshPublicKeyBytes := ssh.MarshalAuthorizedKey(sshPublicKey)

    //Step 5: Save public key to file in OpenSSH format
    if err := os.WriteFile(".ssh/id_ed25519.pub", sshPublicKeyBytes, 0644); err != nil {
        log.Fatalf("Error saving public key", err)
    }
    
    //=== Establish ssh server ===
// Step 6: Start SSH server
	private, err := ssh.NewSignerFromKey(privateKey)
	if err != nil {
		log.Fatalf("Error creating SSH signer: %v", err)
	}

	config := &ssh.ServerConfig{
		NoClientAuth: true,
	}
	config.AddHostKey(private)

	listener, err := net.Listen("tcp", "127.0.0.1:2222")
	if err != nil {
		log.Fatalf("Failed to listen for connections: %v", err)
	}
	defer listener.Close()

	fmt.Println("SSH server listening on 127.0.0.1:2222")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn, config)
	}
}

func handleConnection(conn net.Conn, config *ssh.ServerConfig) {
	defer conn.Close()

	// Perform SSH handshake
	sshConn, channels, requests, err := ssh.NewServerConn(conn, config)
	if err != nil {
		log.Printf("Failed to handshake: %v", err)
		return
	}
	defer sshConn.Close()
	log.Printf("New SSH connection from %s", sshConn.RemoteAddr())

	// Discard global requests
	go ssh.DiscardRequests(requests)

	// Handle channels
	for newChannel := range channels {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unsupported channel type")
			continue
		}

		channel, _, err := newChannel.Accept()
		if err != nil {
			log.Printf("Could not accept channel: %v", err)
			continue
		}
		go handleChannel(channel)
	}
}

func handleChannel(channel ssh.Channel) {
	defer channel.Close()

	// Automatically start a shell
	cmd := "/usr/bin/bash"

	// Set the channel as stdin/stdout/stderr for the shell process
	shell := exec.Command(cmd)
	shell.Stdin = channel
	shell.Stdout = channel
	shell.Stderr = channel

	// Run the shell
	if err := shell.Run(); err != nil {
		log.Printf("Error running shell: %v", err)
	}
}
