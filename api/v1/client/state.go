package client

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"go-ticktick/api/v1/types/project"
)

type TickTickState struct {
	*TickTickClient
	Projects []*project.Project
	*sync.Mutex
}

func NewTickTickState() (*TickTickState, error) {
	clientID := os.Getenv("TT_CLIENT_ID")
	if clientID == "" {
		return nil, errors.New("TT_CLIENT_ID environment variable is not set")
	}
	clientSecret := os.Getenv("TT_CLIENT_SECRET")
	if clientSecret == "" {
		return nil, errors.New("TT_CLIENT_SECRET environment variable is not set")
	}
	redirectUri := os.Getenv("TT_REDIRECT_URI")
	if redirectUri == "" {
		return nil, errors.New("TT_REDIRECT_URI environment variable is not set")
	}
	c := NewTickTickClient(clientID, clientSecret, redirectUri)
	err := c.Authenticate()
	if err != nil {
		return nil, err
	}
	ts := &TickTickState{
		TickTickClient: c,
		Mutex:          &sync.Mutex{},
	}
	err = ts.GetAll()
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func (ts *TickTickState) GetAll() error {
	ts.Lock()
	defer ts.Unlock()
	projs, err := ts.TickTickClient.GetAllProjects(true)
	if err != nil {
		return err
	}
	inbox, err := ts.GetInbox()
	if err != nil {
		return err
	}
	projs = append(projs, inbox)
	ts.Projects = projs
	return nil
}

func (ts *TickTickState) WriteToJSON() error {
	data, err := json.Marshal(ts.Projects)
	if err != nil {
		return err
	}
	file, err := os.Create("state.json")
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	return err
}
