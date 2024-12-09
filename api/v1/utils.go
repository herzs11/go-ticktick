package v1

import (
	`encoding/json`
	`log`
	`net`
	`os`
	`strconv`
	
	`github.com/zalando/go-keyring`
)

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
	return keyring.Set(KEYRING_SERVICE, clientID, string(data))
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

func getTokenFromEnvironment() *oauthToken {
	authToken := os.Getenv("TT_ACCESS_TOKEN")
	token := newOauthTokenFromString(authToken)
	if !token.validate() {
		log.Println("Could not validate token from keyring")
		return nil
	}
	return &token
}
