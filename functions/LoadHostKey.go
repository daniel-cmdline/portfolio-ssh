package functions

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

// LoadHostKey carrega a chave do disco ou gera uma efêmera em memória.
// FIX: Começa com L maiúsculo para o main.go conseguir enxergar!
func LoadHostKey() (ssh.Signer, error) {
	hostKeyPath := os.Getenv("SSH_HOST_KEY_PATH")
	if hostKeyPath == "" {
		hostKeyPath = "id_rsa"
	}

	privateBytes, err := os.ReadFile(hostKeyPath)
	if err == nil {
		return ssh.ParsePrivateKey(privateBytes)
	}
	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to read host key %q: %w", hostKeyPath, err)
	}

	log.Printf("Host key %q not found; generating ephemeral SSH host key for this process.", hostKeyPath)
	return generateEphemeralHostKey()
}

// Fica com g minúsculo porque só a LoadHostKey precisa chamar ela aqui dentro
func generateEphemeralHostKey() (ssh.Signer, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ephemeral host key: %w", err)
	}

	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	return ssh.ParsePrivateKey(pemBytes)
}