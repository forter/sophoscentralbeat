package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/go-yaml/yaml"
)

var (
	cipherKey                = []byte("0123456789012345")
	cipherKeyV2              = []byte("CCEF7CFA0DCB2237012FAE9EB09CCD70")
	clientsCipherKeyPath     = "/app/cmd/beats/cipherstore/"
	clientsCipherKeyFileName = "cipher_key.yml"
)

const (
	encV1 = 1
	encV2 = 2
)

//Encrypt function is used to encrypt the string
func Encrypt(message string) (encmess string, err error) {
	var mainCipherKey []byte
	if len(strings.TrimSpace(message)) == 0 {
		return "", errors.New("string is empty")
	}
	plainText := []byte(message)

	clientsCipherKey, err := GetClientsCipherKey()
	if err != nil {
		logp.Debug("No key with message : ", "%v", err) // no return since we want to use default encryption in such cases
	}
	if err == nil && clientsCipherKey != "" {
		mainCipherKey = []byte(clientsCipherKey)
	} else {
		mainCipherKey = cipherKeyV2
	}

	block, err := aes.NewCipher(mainCipherKey)
	if err != nil {
		return "", err
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//returns to base64 encoded string
	encmess = base64.StdEncoding.EncodeToString(cipherText)
	finalEnc := fmt.Sprintf("%d%s%s", encV2, "||", encmess)
	return finalEnc, nil
}

//EncryptV1 function is used to encrypt the string
func EncryptV1(message string) (encmess string, err error) {
	if len(strings.TrimSpace(message)) == 0 {
		return "", errors.New("string is empty")
	}
	plainText := []byte(message)

	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		return "", err
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//returns to base64 encoded string
	encmess = base64.URLEncoding.EncodeToString(cipherText)
	finalEnc := fmt.Sprintf("%d%s%s", encV1, "||", encmess)
	return finalEnc, nil
}

// CypherKeyStruct encapsulates cypher key data
type CypherKeyStruct struct {
	CypherKey string `yaml:"cypher_key"`
}

// GetClientsCipherKey is to get the cipher key of the client if any found
func GetClientsCipherKey() (string, error) {
	path := filepath.Join(clientsCipherKeyPath, url.QueryEscape(clientsCipherKeyFileName))
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	var cypherKeyVal CypherKeyStruct
	err = yaml.Unmarshal(data, &cypherKeyVal)
	if err != nil {
		return "", err
	}
	return cypherKeyVal.CypherKey, nil
}
