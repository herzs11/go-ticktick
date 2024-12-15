package client

import (
	`log`
	`os`
	`testing`
	
	`github.com/joho/godotenv`
)

func TestLogin(t *testing.T) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatal(err)
	}
	ttPass := os.Getenv("TT_PASS")
	if ttPass == "" {
		log.Fatal("TT_PASS environment variable not set")
	}
	ttUser := os.Getenv("TT_USER")
	if ttUser == "" {
		log.Fatal("TT_USER environment variable not set")
	}
	c := NewClient(ttUser, ttPass)
	err = c.login()
	if err != nil {
		log.Fatal(err)
	}
	err = c.getState()
	if err != nil {
		log.Fatal(err)
	}
}
