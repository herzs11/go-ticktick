package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/herzs11/go-ticktick/api/v1/types/project"
	"github.com/herzs11/go-ticktick/api/v1/types/tasks"
)

func (c *TickTickClient) getTaskByProjectIdAndTaskID(projectID, taskID string) (*tasks.Task, error) {
	req, err := http.NewRequest(
		"GET", fmt.Sprintf("%s%s/%s/task/%s", API_BASE_URL, project.PROJECT_ENDPOINT, projectID, taskID), nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var t tasks.Task
	err = json.NewDecoder(resp.Body).Decode(&t)
	if err != nil {
		return nil, err
	}
	if t.Id == "" {
		return nil, errors.New(fmt.Sprintf("Task with id %s not found with project '%s'", taskID, projectID))
	}
	return &t, nil

}

func (c *TickTickClient) GetTaskById(taskID string) (*tasks.Task, error) {
	return c.getTaskByProjectIdAndTaskID("inbox", taskID)
}

func (c *TickTickClient) GetTask(task *tasks.Task) error {
	if task.ProjectId == "" || strings.HasPrefix(task.ProjectId, "inbox") {
		t, err := c.GetTaskById(task.Id)
		if err != nil {
			return err
		}
		*task = *t
		return nil
	}
	t, err := c.getTaskByProjectIdAndTaskID(task.ProjectId, task.Id)
	if err != nil {
		return err
	}
	*task = *t
	return nil
}

func (c *TickTickClient) CompleteTask(task *tasks.Task) error {
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s/%s/task/%s/complete", API_BASE_URL, project.PROJECT_ENDPOINT, task.ProjectId, task.Id), nil,
	)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Task completion returned a non-200 status code: %d", resp.StatusCode))
	}
	task.Status = tasks.Completed
	task.CompletedTime = time.Now().Local()
	return nil
}

func (c *TickTickClient) DeleteTask(task *tasks.Task) error {
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s%s/%s/task/%s", API_BASE_URL, project.PROJECT_ENDPOINT, task.ProjectId, task.Id), nil,
	)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Task completion returned a non-200 status code: %d", resp.StatusCode))
	}
	return nil
}

func (c *TickTickClient) CreateTask(task *tasks.Task) error {
	err := validateTask(task)
	if err != nil {
		return err
	}

	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", API_BASE_URL, tasks.TASK_ENDPOINT), bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(task)
}

func (c *TickTickClient) UpdateTask(task *tasks.Task) error {
	err := validateUpdateTask(task)
	if err != nil {
		return err
	}

	data, err := json.Marshal(task)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST", fmt.Sprintf("%s%s/%s", API_BASE_URL, tasks.TASK_ENDPOINT, task.Id), bytes.NewBuffer(data),
	)

	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var respRaw interface{}
		err = json.NewDecoder(resp.Body).Decode(&respRaw)
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(task)
	return err
}
