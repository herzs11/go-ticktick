package v2

import (
	`net/http`
	
	`github.com/valyala/fasttemplate`
)

const (
	BASE_URL                = "https://api.ticktick.com/api/v2/"
	USER_AGENT_TEMPLATE_STR = `{"platform":"web","os":"OS X","device":"Firefox 95.0","name":"go-ticktick", "version":4531,
"id":"6490[token]","channel":"website","campaign":"","websocket":""}`
)

var (
	USER_AGENT_TEMPLATE = fasttemplate.New(USER_AGENT_TEMPLATE_STR, "[", "]")
)

type Client struct {
	username    string
	password    string
	AccessToken string
	Cookies     map[string]string
}

type loginParams struct {
	username string `json:"username"`
	password string `json:"password"`
}

func NewClient(username, password string) *Client {
	return &Client{
		username: username,
		password: password,
	}
}

func (c *Client) login() error {
	u := BASE_URL + "user/signin"
	httpClient := http.Client{}
	httpClient.Post(u, "application/json")
}
