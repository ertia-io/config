package entities

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"strings"
	"time"

	"github.com/mikesmitty/edkey"
	"github.com/teris-io/shortid"
	"golang.org/x/crypto/ssh"
)

var (
	KeyStatusNew      = "NEW"
	KeyStatusAdapting = "ADAPTING"
	KeyStatusActive   = "ACTIVE"
	KeyStatusFailing  = "FAILING"
	KeyStatusDeleted  = "DELETED"
)

type SSHKey struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	Fingerprint string    `json:"fingerprint"`
	Error       string    `json:"error"`
	PrivateKey  string    `json:"privateKey"`
	PublicKey   string    `json:"publicKey"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

func (k *SSHKey) NeedsAdapting() bool {
	if k.Status == KeyStatusNew {
		return true
	}
	return false
}

// GenerateKeyPair creates a SSH key pair (Ed25519).
func GenerateKeyPair() (*SSHKey, error) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	// Public key instance.
	publicKeyInst, err := ssh.NewPublicKey(pubKey)
	if err != nil {
		return nil, err
	}

	id, err := shortid.Generate()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	// Serialize public key use in OpenSSH authorized_keys.
	publicKey := ssh.MarshalAuthorizedKey(publicKeyInst)

	// Encode private key in PEM format.
	privateKey := pem.EncodeToMemory(&pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: edkey.MarshalED25519PrivateKey(privKey),
	})

	return &SSHKey{
		ID:          id,
		Name:        id,
		PrivateKey:  string(privateKey),
		PublicKey:   strings.TrimRight(string(publicKey), "\n"),
		Fingerprint: "",
		Status:      KeyStatusNew,
		Created:     now,
		Updated:     now,
	}, nil
}
