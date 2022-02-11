package entities

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/teris-io/shortid"
	"golang.org/x/crypto/ssh"
	"log"
	"time"
)

var(
	KeyStatusNew = "NEW"
	KeyStatusAdapting = "ADAPTING"
	KeyStatusActive = "ACTIVE"
	KeyStatusFailing = "FAILING"
	KeyStatusDeleted = "DELETED"
)

type SSHKey struct {
	ID string `json:"id""`
	ProviderID string `json:"providerId"`
	Name string `json:"name"`
	Status string `json:"status"`
	Fingerprint string `json:"fingerprint"`
	Error string `json:"error"`
	PrivateKey string `json:"privateKey"`
	PublicKey string `json:"publicKey"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

func(k *SSHKey) NeedsAdapting() bool{
	if(k.Status == KeyStatusNew){
		return true
	}
	return false
}

func GetPublicKeys() (*SSHKey,*rsa.PrivateKey, error){

	id, err := shortid.Generate()

	if(err!=nil){
		return nil,nil, err
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)

	if err != nil {
		return nil,nil, err
	}


	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil,nil, err
	}

/*
	err = ioutil.WriteFile(config.LubeKeysPath()+"/"+id,encodePrivateKeyToPEM(privateKey),0600)
	if(err!=nil){
		return nil, err
	}
*/


	pubKey, err := generatePublicKey(&privateKey.PublicKey)
	if(err!=nil){
		return nil, privateKey, err
	}
/*
	err = ioutil.WriteFile(config.LubeKeysPath()+"/"+id+".pub",pubKey,0600)
	if(err!=nil){
		return nil, err
	}

 */

	return &SSHKey{
			ID: id,
			Name: id,
			PublicKey: string(pubKey),
			Status: KeyStatusNew,

	},privateKey, nil

}

// encodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

// generatePublicKey take a rsa.PublicKey and return bytes suitable for writing to .pub file
// returns in the format "ssh-rsa ..."
func generatePublicKey(privatekey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privatekey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	log.Println("Public key generated")
	return pubKeyBytes, nil
}