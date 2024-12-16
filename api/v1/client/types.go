package client

import (
	"errors"
	"fmt"

	"go-ticktick/api/v1/types/project"
	"go-ticktick/api/v1/types/tasks"
)

type projectParams struct {
	Name     string `json:"name"`
	Color    string `json:"color,omitempty"`
	ViewMode string `json:"viewMode"`
	Kind     string `json:"kind"`
}

func validateProject(np *project.Project) error {
	if np.Name == "" {
		return errors.New("Must provide a name for the project")
	}
	if np.Color != "" && !validateRGBHex(np.Color) {
		return errors.New(fmt.Sprintf("%s is not a valid RGB hex code", np.Color))
	}
	if np.ViewMode.String() == "" {
		return errors.New(fmt.Sprintf("%d is not a valid project ViewMode", np.ViewMode))
	}
	if np.Kind.String() == "" {
		return errors.New(fmt.Sprintf("%d is not a valid project Kind", np.Kind))
	}
	return nil
}

func validateUpdateTask(t *tasks.Task) error {
	if t.Id == "" {
		return errors.New("task with an empty id cannot be updated")
	}
	return validateTask(t)
}

func validateTask(t *tasks.Task) error {
	if t.Title == "" {
		return errors.New("Task must have a title")
	}
	if t.ProjectId == "" {
		t.ProjectId = "inbox"
	}
	if t.Status.String() == "" {
		return errors.New("task has invalid Status")
	}
	if t.Priority.String() == "" {
		return errors.New("task has invalid priority")
	}
	if !t.DueDate.IsZero() && t.StartDate.IsZero() {
		t.StartDate = t.DueDate
	}
	for _, cl := range t.ChecklistItems {
		if err := validateChecklistItem(&cl); err != nil {
			return err
		}
	}
	return nil
}

func validateChecklistItem(cl *tasks.ChecklistItem) error {
	if cl.Title == "" {
		return errors.New("checklist item has empty title")
	}
	if cl.Status.String() == "" {
		return errors.New("checklist item has an invalid status")
	}
	return nil
}
