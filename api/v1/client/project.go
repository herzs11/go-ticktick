package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/herzs11/go-ticktick/api/v1/types/project"
)

func (c *TickTickClient) CreateNewProjectFromName(name string) (*project.Project, error) {
	p := project.Project{Name: name}
	err := c.CreateNewProject(&p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (c *TickTickClient) CreateNewProject(proj *project.Project) error {
	err := validateProject(proj)
	if err != nil {
		return err
	}
	projParams := &projectParams{
		Name:     proj.Name,
		Color:    proj.Color,
		ViewMode: proj.ViewMode.String(),
		Kind:     proj.ViewMode.String(),
	}

	data, err := json.Marshal(projParams)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", API_BASE_URL+project.PROJECT_ENDPOINT, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(proj)
}

func (c *TickTickClient) GetAllProjects(includeTasks bool) ([]*project.Project, error) {
	req, err := http.NewRequest("GET", API_BASE_URL+project.PROJECT_ENDPOINT, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var projs []*project.Project
	err = json.NewDecoder(resp.Body).Decode(&projs)
	if !includeTasks {
		return projs, err
	}
	var pTasks []*project.Project
	for _, p := range projs {
		pt, err := c.GetProjectById(p.Id, true)
		if err != nil {
			return nil, err
		}
		pTasks = append(pTasks, pt)
	}
	return pTasks, nil
}

func (c *TickTickClient) GetInbox() (*project.Project, error) {
	p, err := c.GetProjectById("inbox", true)
	if err != nil {
		return nil, err
	}
	p.Id = "inbox"
	p.Name = "Inbox"
	return p, nil
}

func (c *TickTickClient) GetProjectById(id string, includeTasks bool) (*project.Project, error) {
	if strings.HasPrefix(id, "inbox") && !includeTasks {
		return nil, errors.New("must return tasks with inbox")
	}
	u := fmt.Sprintf("%s%s/%s", API_BASE_URL, project.PROJECT_ENDPOINT, id)
	if includeTasks {
		u = u + "/data"
	}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var p project.Project
	err = json.NewDecoder(resp.Body).Decode(&p)
	return &p, err
}

func (c *TickTickClient) DeleteProjectById(id string) error {
	if strings.HasPrefix(id, "inbox") {
		return errors.New("cannot delete inbox project")
	}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s%s/%s", API_BASE_URL, project.PROJECT_ENDPOINT, id), nil)
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Got status code %d", resp.StatusCode))
	}
	return nil
}

func (c *TickTickClient) UpdateProject(proj *project.Project) error {
	if strings.HasPrefix(
		proj.Id,
		"inbox",
	) {
		return errors.New("cannot update inbox project")
	}
	err := validateProject(proj)
	if err != nil {
		return err
	}
	projParams := &projectParams{
		Name:     proj.Name,
		Color:    proj.Color,
		ViewMode: proj.ViewMode.String(),
		Kind:     proj.ViewMode.String(),
	}

	data, err := json.Marshal(projParams)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST", fmt.Sprintf("%s%s/%s", API_BASE_URL, project.PROJECT_ENDPOINT, proj.Id), bytes.NewBuffer(data),
	)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(proj)
}
