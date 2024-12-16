package client

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"path"
	"regexp"
	"strconv"

	"github.com/zalando/go-keyring"
)

const OAUTH2_FILENAME = ".gott_auth2"

func checkPort(port string) bool {
	// Attempt to listen on the port
	p, _ := strconv.Atoi(port)
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(p))
	if err != nil {
		return true
	}
	ln.Close()
	return false
}

func storeToken(clientID string, token oauthToken) error {
	err := storeTokenKeyring(clientID, token)
	if err == nil {
		return nil
	}
	log.Printf("Unable to store token in keyring service, saving to ~/.gott_oauth2")
	return storeTokenFile(token)
}

func storeTokenKeyring(clientID string, token oauthToken) error {
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}
	return keyring.Set(KEYRING_SERVICE, clientID, string(data))
}

func storeTokenFile(token oauthToken) error {
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	file, err := os.Create(path.Join(home, OAUTH2_FILENAME))
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, bytes.NewBuffer(data))
	return err
}

func getTokenFromKeyring(clientID string) (*oauthToken, error) {
	secret, err := keyring.Get(KEYRING_SERVICE, clientID)
	if err != nil {
		return nil, err
	}
	var retrievedToken oauthToken
	err = json.Unmarshal([]byte(secret), &retrievedToken)
	if err != nil {
		return nil, err
	}
	if !retrievedToken.validate() {
		return nil, errors.New("Unable to validate token from keyring")
	}
	return &retrievedToken, nil
}

func getTokenFromFile() (*oauthToken, error) {
	authToken := &oauthToken{}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(path.Join(home, OAUTH2_FILENAME))
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(file).Decode(authToken)
	if err != nil {
		return nil, err
	}
	if !authToken.validate() {
		return nil, errors.New("Unable to validate token from keyring")
	}
	return authToken, nil
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func validateRGBHex(s string) bool {
	matched, _ := regexp.MatchString(`^#[0-9a-fA-F]{6}$`, s)
	return matched
}
