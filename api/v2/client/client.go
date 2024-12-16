package client

import (
	`bytes`
	`encoding/json`
	`fmt`
	`io`
	`net/http`
	`net/http/cookiejar`
	`os`
	
	`github.com/valyala/fasttemplate`
)

const (
	USER_AGENT            = "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:129.0) Gecko/20100101 Firefox/129.0"
	BASE_URL              = "https://api.ticktick.com/api/v2/"
	X_DEVICE_TEMPLATE_STR = `{"platform":"web","os":"OS X","device":"Firefox 95.0","name":"go-ticktick", "version":4531,"id":"6491[token]","channel":"website","campaign":"","websocket":""}`
)

var (
	X_DEVICE_TEMPLATE = fasttemplate.New(X_DEVICE_TEMPLATE_STR, "[", "]")
)

type Client struct {
	username    string
	password    string
	AccessToken string
	header      *http.Header
	userId      string
	httpClient  *http.Client
	InboxId     string
}

type loginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewClient(username, password string) *Client {
	h := &http.Header{}
	h.Add("User-Agent", USER_AGENT)
	randToken, _ := randomHex(10)
	xDev := X_DEVICE_TEMPLATE.ExecuteString(
		map[string]interface{}{
			"token": randToken,
		},
	)
	h.Add("x-device", xDev)
	jar, _ := cookiejar.New(nil)
	return &Client{
		username: username,
		password: password,
		header:   h,
		httpClient: &http.Client{
			Jar: jar,
		},
	}
}

func (c *Client) login() error {
	u := BASE_URL + "user/signon"
	lp := loginParams{Username: c.username, Password: c.password}
	params, err := json.Marshal(lp)
	if err != nil {
		return fmt.Errorf("Failed to marshal login params: %s", err.Error())
	}
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(params))
	if err != nil {
		return fmt.Errorf("Failed to create request object: %s", err.Error())
	}
	req.URL.Query().Add("wc", "true")
	req.URL.Query().Add("remember", "true")
	req.Header = *c.header
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Referer", "https://ticktick.com")
	req.Header.Add("Origin", "https://ticktick.com")
	resp, err := c.httpClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("Login request failed: %s", err.Error())
	}
	var responseBody loginResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return fmt.Errorf("Got unexpected response: %s", err)
	}
	
	return nil
}

func (c *Client) getState() error {
	u := BASE_URL + "batch/check/0"
	req, _ := http.NewRequest("GET", u, nil)
	req.Header = *c.header
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	file, _ := os.Create("testout.json")
	defer file.Close()
	
	io.Copy(file, resp.Body)
	
	return nil
}
