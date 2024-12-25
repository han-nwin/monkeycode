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

		// Perform SSH handshake
		sshConn, chans, reqs, err := ssh.NewServerConn(conn, config)
		if err != nil {
			log.Printf("Failed to handshake: %v", err)
			conn.Close()
			continue
		}

		// Handle incoming channels (e.g., terminal sessions)
		go handleConnection(sshConn, chans, reqs)
	}
}

// Handle SSH connection and session setup
func handleConnection(conn ssh.Conn, chans <-chan ssh.NewChannel, reqs <-chan *ssh.Request) {
	for newChannel := range chans {
		// Accept the channel of type "session"
		if newChannel.ChannelType() == "session" {
			channel, _, err := newChannel.Accept()
			if err != nil {
				log.Printf("could not accept channel: %v", err)
				continue
			}

			// Handle requests for this channel
			go func() {
				for req := range reqs {
					go handleRequest(channel, req)
				}
			}()
		}
	}
}

// Handle incoming SSH requests
func handleRequest(channel ssh.Channel, req *ssh.Request) {
	switch req.Type {
	case "exec":
		// Handle exec command for terminal interaction
		cmd := exec.Command("bash")
		cmd.Env = append(os.Environ(), "TERM=xterm-256color")
		cmd.Stdin = channel
		cmd.Stdout = channel
		cmd.Stderr = channel

		if err := cmd.Start(); err != nil {
			log.Printf("failed to start bash: %v", err)
			channel.Write([]byte("Failed to start bash\n"))
			return
		}

		// Wait for the command to finish
		if err := cmd.Wait(); err != nil {
			log.Printf("failed to wait for bash process: %v", err)
		}
	}
}

