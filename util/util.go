/*
	Copyright © 2023 Gabriel Pozo

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/
package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func IsValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func IsValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(regex, email)
	return match
}

func FileNameWithoutExtSliceNotation(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

func RemoveFile(fileName string) error {
	return os.Remove(fileName)
}

func GetNameFilesConfig() []string {
	entries, err := os.ReadDir("config")
	if err != nil {
		log.Fatalf("Can't files from config folder %s", err)
	}
	var names []string
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			log.Fatalf("Error whiles reading files %s", err)
		}
		names = append(names, info.Name())
	}

	return names
}

func ConfigFolderCreate() {
	path := "config"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func Encrypt(plaintext, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertextByte := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	chipertext := base64.StdEncoding.EncodeToString(ciphertextByte)
	return chipertext, nil
}

func Decrypt(cipherText, key string) (string, error) {
	// prepare cipher
	keyByte := []byte(key)
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	//// process ciphertext
	ciphertextByte, _ := base64.StdEncoding.DecodeString(cipherText)
	nonce, ciphertextByteClean := ciphertextByte[:nonceSize], ciphertextByte[nonceSize:]
	plaintextByte, err := gcm.Open(nil, nonce, ciphertextByteClean, nil)
	if err != nil {
		return "", err
	}
	return string(plaintextByte), err
}
