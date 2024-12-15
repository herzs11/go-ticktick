package client

import (
	`bytes`
	`encoding/json`
	`io`
	`log`
	`net`
	`os`
	`strconv`
	
	`github.com/zalando/go-keyring`
)

const OAUTH2_FILENAME = "~/.gott_auth2"

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
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}
	err = keyring.Set(KEYRING_SERVICE, clientID, string(data))
	if err != nil {
		log.Printf("Unable to store token in keyring service, saving to ~/.gott_oauth2")
	}
	file, err := os.Create(OAUTH2_FILENAME)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, bytes.NewBuffer(data))
	return err
}

func getTokenFromKeyring(clientID string) *oauthToken {
	secret, err := keyring.Get(KEYRING_SERVICE, clientID)
	if err != nil {
		log.Printf("Error getting token from keyring: %s\n", err.Error())
		return nil
	}
	var retrievedToken oauthToken
	err = json.Unmarshal([]byte(secret), &retrievedToken)
	if err != nil {
		log.Printf("Error unmarshalling into token: %s\n", err.Error())
		return nil
	}
	if !retrievedToken.validate() {
		log.Println("Could not validate token from keyring")
		return nil
	}
	return &retrievedToken
}

func getTokenFromFile() *oauthToken {
	authToken := &oauthToken{}
	file, err := os.Open(OAUTH2_FILENAME)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Unable to get oauthtoken from file: ~/.gott_oauth2 file does not exist")
		} else {
			log.Printf("Unable to open file")
		}
		return nil
	}
	err = json.NewDecoder(file).Decode(authToken)
	if err != nil {
		log.Printf("Unable to parse json from file")
		return nil
	}
	if !authToken.validate() {
		log.Printf("Unable to validate token, login again")
		return nil
	}
	return authToken
}

func getTokenFromEnvironment() *oauthToken {
	authToken := os.Getenv("TT_ACCESS_TOKEN")
	token := newOauthTokenFromString(authToken)
	if !token.validate() {
		log.Println("Could not validate token from keyring")
		return nil
	}
	return &token
}
