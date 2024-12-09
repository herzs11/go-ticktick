package v1

import (
	`encoding/json`
	`fmt`
	`os`
	`testing`
	`time`
	
	`github.com/joho/godotenv`
	`github.com/zalando/go-keyring`
)

func TestOauth2GetToken(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Log("Could not load .env file")
	}
	testClientID := os.Getenv("TT_CLIENT_ID")
	if testClientID == "" {
		t.Fatal("Client id must be set in the TT_CLIENT_ID environment variable")
	}
	testClientSecret := os.Getenv("TT_CLIENT_SECRET")
	testRedirectUri := os.Getenv("TT_REDIRECT_URI")
	c := NewOauth2Client(testClientID, testClientSecret, testRedirectUri)
	err = c.getAuthorizationCode()
	if err != nil {
		t.Fatal("Error getting authorization code: ", err)
	}
	if c.authorizationCode == "" {
		t.Fatal("Could not get authorization code")
	}
	fmt.Println("Got authorization code from listener server", c.authorizationCode)
	if checkPort("8080") {
		t.Fatal("Port still in use")
	}
	
	err = c.getOauthToken()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(c.token.AccessToken)
}

func TestTokenStore(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Log("Could not load .env file")
	}
	testClientID := os.Getenv("TT_CLIENT_ID")
	if testClientID == "" {
		t.Fatal("Client id must be set in the TT_CLIENT_ID environment variable")
	}
	token := oauthToken{
		AccessToken: "f563026e-443d-46b1-89c7-a9f0c0bf907f",
		TokenType:   "bearer",
		ExpiresIn:   15125617,
		Scope:       "tasks:read tasks:write",
	}
	expTS := time.Second * time.Duration(token.ExpiresIn)
	token.ExpiresTime = time.Now().Add(expTS).Unix()
	data, err := json.Marshal(token)
	if err != nil {
		t.Fatal(err)
	}
	
	service := "go-ticktick"
	err = keyring.Set(service, testClientID, string(data))
	if err != nil {
		t.Fatal(err)
	}
	secret, err := keyring.Get(service, testClientID)
	if err != nil {
		t.Fatal(err)
	}
	retrievedData := []byte(secret)
	var retrievedToken oauthToken
	err = json.Unmarshal(retrievedData, &retrievedToken)
	if err != nil {
		t.Fatalf("Unable to unmarshal data: %s", err.Error())
	}
	
	if token.AccessToken != retrievedToken.AccessToken {
		t.Fatalf("Expected token %s, got %s", token.AccessToken, retrievedToken.AccessToken)
	}
	if token.ExpiresTime != retrievedToken.ExpiresTime {
		t.Fatalf("Expected expTime %d, got %d", token.ExpiresTime, retrievedToken.ExpiresTime)
	}
}

func TestOauth2Client_Authenticate(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Log("Could not load .env file")
	}
	testClientID := os.Getenv("TT_CLIENT_ID")
	if testClientID == "" {
		t.Fatal("Client id must be set in the TT_CLIENT_ID environment variable")
	}
	testClientSecret := os.Getenv("TT_CLIENT_SECRET")
	testRedirectUri := os.Getenv("TT_REDIRECT_URI")
	c := NewOauth2Client(testClientID, testClientSecret, testRedirectUri)
	err = c.Authenticate()
	if err != nil {
		t.Fatalf("Error authenticating oauth client: %s", err.Error())
	}
	
	token := getTokenFromKeyring(c.ClientId)
	if token.AccessToken != c.token.AccessToken {
		t.Fatalf("Expected token %s from keyring, got %s", c.token.AccessToken, token.AccessToken)
	}
}
