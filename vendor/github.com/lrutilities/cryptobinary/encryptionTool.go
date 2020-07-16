package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var (
	cipherKey   = []byte("0123456789012345")
	cipherKeyV2 = []byte("CCEF7CFA0DCB2237012FAE9EB09CCD70")
)

const (
	encV1 = 1
	encV2 = 2
)

//Encrypt function is used to encrypt the string
func Encrypt(message string) (encmess string, err error) {
	plainText := []byte(message)

	block, err := aes.NewCipher(cipherKeyV2)
	if err != nil {
		return
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

// Decrypt function is used to decrypt the string
func Decrypt(securemess string) (decodedmess string, err error) {

	if len(strings.TrimSpace(securemess)) == 0 {
		return "", errors.New("string is empty")
	}
	decodedStr := strings.Split(securemess, "||")
	if len(decodedStr) == 2 {
		ver, err := strconv.Atoi(decodedStr[0])
		if err != nil {
			return "", err
		}
		switch ver {
		case encV1:
			decodedmess, err = decrypt1(decodedStr[1])
			if err != nil {
				return "", err
			}
		case encV2:
			decodedmess, err = decrypt2(decodedStr[1])
			if err != nil {
				return "", err
			}
		default:
			return "", errors.New("invalid encryption")
		}
	}

	return decodedmess, nil
}
func main() {
	arg := os.Args[1]
	argDecrypt := os.Args[2]
	file := os.Args[3]
	if argDecrypt == "encrypt" {
		if file == "1" {
			err := encryptFile(arg)
			if err != nil {
				fmt.Printf("error:%v", err)
			}
		} else {
			encrptedVal, err := Encrypt(arg)
			if err != nil {
				fmt.Printf("error:%v", err)
			}
			fmt.Println(encrptedVal)
		}
	} else {
		if file == "1" {
			err := decryptFile(arg)
			if err != nil {
				fmt.Printf("error:%v", err)
			}
		} else {
			decrptedVal, err := Decrypt(arg)
			if err != nil {
				fmt.Printf("error:%v", err)
			}
			fmt.Println(decrptedVal)
		}
	}
}

func decryptFile(filePath string) error {
	c, err := ioutil.ReadFile(filePath) // just pass the file name
	if err != nil {
		return fmt.Errorf("fail to decrypt credentials: %v", err)
	}

	decryptedContent, err := Decrypt(string(c))
	if err != nil {
		return errors.New("error decrypting Content")
	}
	fmt.Println(decryptedContent)
	return nil
}

func encryptFile(filePath string) error {

	c, err := ioutil.ReadFile(filePath) // just pass the file name
	if err != nil {
		return fmt.Errorf("fail to encrypted credentials: %v", err)
	}

	encryptedContent, err := Encrypt(string(c))
	if err != nil {
		return errors.New("error encrypting Content")
	}
	err = ioutil.WriteFile(filePath+".encrypt", []byte(encryptedContent), 0)
	if err != nil {
		return fmt.Errorf("fail to encrypt credentials: %v", err)
	}
	fmt.Println(encryptedContent)
	return nil
}

func decrypt1(securemess string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(securemess)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("ciphertext block size is too short")
		return "", err
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	decodedmess := string(cipherText)
	return decodedmess, nil
}

func decrypt2(securemess string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(securemess)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(cipherKeyV2)
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("ciphertext block size is too short")
		return "", err
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	decodedmess := string(cipherText)
	return decodedmess, nil
}
