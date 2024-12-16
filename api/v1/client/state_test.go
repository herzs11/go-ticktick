package client

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestNewTickTickState(t *testing.T) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Log("Could not load .env file")
	}
	ts, err := NewTickTickState()
	if err != nil {
		t.Fatal(err)
	}

	ts.WriteToJSON()
}
