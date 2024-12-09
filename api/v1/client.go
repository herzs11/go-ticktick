package v1

import (
	`context`
	`encoding/json`
	`errors`
	`fmt`
	`log`
	`net/http`
	`net/url`
	`time`
	
	"github.com/pkg/browser"
)

const (
	OAUTH_BASE_URL              = "https://ticktick.com/oauth"
	AUTHORIZATION_PAGE_ENDPOINT = "/authorize"
	ACCESS_TOKEN_ENDPOINT       = "/token"
	KEYRING_SERVICE             = "go-ticktick"
	API_BASE_URL                = "https://api.ticktick.com"
)

type oauthToken struct {
	AccessToken string  `json:"access_token"`
	TokenType   string  `json:"token_type"`
	ExpiresIn   float64 `json:"expires_in"`
	ExpiresTime int64   `json:"expires_time"`
	Scope       string  `json:"scope"`
}

func (t *oauthToken) validate() bool {
	if t.AccessToken == "" {
		return false
	}
	if time.Now().Unix() > t.ExpiresTime {
		return false
	}
	client := &http.Client{Timeout: time.Second * 5}
	req, err := http.NewRequest("GET", API_BASE_URL+"/open/v1/project", nil)
	if err != nil {
		return false
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.AccessToken))
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	return resp.StatusCode == 200
}

func newOauthTokenFromString(token string) oauthToken {
	return oauthToken{
		AccessToken: token,
		TokenType:   "bearer",
		ExpiresIn:   86400,
		ExpiresTime: time.Now().Add(time.Hour * 24).Unix(),
		Scope:       "tasks:read tasks:write",
	}
}

type Oauth2Client struct {
	ClientId          string
	ClientSecret      string
	RedirectURI       string
	authorizationCode string
	token             oauthToken
}

func createAuthorizationCodeListener(authCh chan string, serv *http.Server, path string) {
	http.HandleFunc(
		path, func(w http.ResponseWriter, r *http.Request) {
			// Parse the query parameters
			queryValues, err := url.ParseQuery(r.URL.RawQuery)
			if err != nil {
				http.Error(w, "Failed to parse query parameters", http.StatusBadRequest)
				return
			}
			
			// Extract the authorization code
			code := queryValues.Get("code")
			if code == "" {
				http.Error(w, "Authorization code not found", http.StatusBadRequest)
				return
			}
			
			// Send the authorization code to the main goroutine
			authCh <- code
			
			// Respond to the user
			fmt.Fprintf(w, "Successfully got authorization code from the redirected url, you may now close this window")
		},
	)
	log.Printf("Server listening on %s\n", serv.Addr)
	if err := serv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

func NewOauth2Client(clientID, clientSecret, redirectURI string) *Oauth2Client {
	return &Oauth2Client{
		ClientId:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
	}
}

func (oc *Oauth2Client) getAuthURL() string {
	params := map[string]string{
		"client_id":     oc.ClientId,
		"response_type": "code",
		"redirect_uri":  oc.RedirectURI,
		"scope":         "tasks:write tasks:read",
	}
	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}
	return fmt.Sprintf("%s%s?%s", OAUTH_BASE_URL, AUTHORIZATION_PAGE_ENDPOINT, values.Encode())
}

func (oc *Oauth2Client) openAuthURL() bool {
	url := oc.getAuthURL()
	err := browser.OpenURL(url)
	if err != nil {
		log.Printf("Error opening url in browser: %s\n", err)
		return false
	}
	return true
}

func (oc *Oauth2Client) makeRedirectServer() (*http.Server, string, error) {
	rdu, err := url.Parse(oc.RedirectURI)
	if err != nil {
		return nil, "", err
	}
	port := rdu.Port()
	if port == "" {
		return nil, "", errors.New("URL does not have a port")
	}
	if checkPort(port) {
	
	}
	// rdu.Path
	serv := &http.Server{Addr: ":" + port}
	if rdu.Path == "" {
		rdu.Path = "/"
	}
	return serv, rdu.Path, nil
	
}

func (oc *Oauth2Client) getAuthorizationCode() error {
	authCh := make(chan string)
	serv, path, err := oc.makeRedirectServer()
	if err != nil {
		return err
	}
	
	go createAuthorizationCodeListener(authCh, serv, path)
	
	res := oc.openAuthURL()
	if !res {
		return errors.New("Unable to open authorization redirect in browser")
	}
	
	oc.authorizationCode = <-authCh
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := serv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	return nil
}

func (oc *Oauth2Client) getOauthToken() error {
	params := map[string]string{
		"client_id":     oc.ClientId,
		"client_secret": oc.ClientSecret,
		"code":          oc.authorizationCode,
		"grant_type":    "authorization_code",
		"scope":         "tasks:write tasks:read",
		"redirect_uri":  oc.RedirectURI,
	}
	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}
	values.Encode()
	resp, err := http.PostForm(
		OAUTH_BASE_URL+ACCESS_TOKEN_ENDPOINT, values,
	)
	if err != nil {
		return err
	}
	
	defer resp.Body.Close()
	
	// var res map[string]interface{}
	var res oauthToken
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return err
	}
	expTS := time.Second * time.Duration(res.ExpiresIn)
	res.ExpiresTime = time.Now().Add(expTS).Unix()
	oc.token = res
	return nil
}

func (oc *Oauth2Client) Authenticate() error {
	token := getTokenFromKeyring(oc.ClientId)
	if token != nil {
		oc.token = *token
		return nil
	}
	
	token = getTokenFromEnvironment()
	if token != nil {
		oc.token = *token
		return nil
	}
	
	log.Println("Cannot get valid token from cache, authenticating...")
	err := oc.getAuthorizationCode()
	if err != nil {
		return err
	}
	err = oc.getOauthToken()
	if err != nil {
		return err
	}
	if !oc.token.validate() {
		return errors.New("Cannot validate token")
	}
	err = storeToken(oc.ClientId, oc.token)
	if err != nil {
		log.Println("Could not store token in keyring, save in environment variable 'TT_ACCESS_TOKEN'")
		log.Println("Token:", oc.token.AccessToken)
	}
	return nil
}
