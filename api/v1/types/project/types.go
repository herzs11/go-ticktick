package project

import (
	"encoding/json"
	
	"github.com/herzs11/go-ticktick/api/v1/types/tasks"
)

type ViewMode int
type Kind int

const (
	List ViewMode = iota
	Kanban
	Timeline
)

const (
	Task Kind = iota
	Note
)

func (v ViewMode) String() string {
	switch v {
	case List:
		return "list"
	case Kanban:
		return "kanban"
	case Timeline:
		return "timeline"
	default:
		return ""
	}
}

func (k Kind) String() string {
	switch k {
	case Task:
		return "TASK"
	case Note:
		return "NOTE"
	default:
		return ""
	}
}

func kindFromString(s string) Kind {
	switch s {
	case "TASK":
		return Task
	case "NOTE":
		return Note
	default:
		return Task
	}
}

func viewModeFromString(s string) ViewMode {
	switch s {
	case "list":
		return List
	case "kanban":
		return Kanban
	case "timeline":
		return Timeline
	default:
		return List
	}
}

type Project struct {
	Id       string
	Name     string
	Color    string
	ViewMode ViewMode
	Kind     Kind
	GroupId  string
	Closed   bool
	Tasks    []tasks.Task
}

type projectJSON struct {
	Id       string       `json:"id,omitempty"`
	Name     string       `json:"name"`
	Color    string       `json:"color,omitempty"`
	ViewMode string       `json:"viewMode"`
	Kind     string       `json:"kind"`
	GroupId  string       `json:"groupId,omitempty"`
	Closed   bool         `json:"closed,omitempty"`
	Tasks    []tasks.Task `json:"tasks"`
}

type projectTaskJSON struct {
	Project projectJSON  `json:"project"`
	Tasks   []tasks.Task `json:"tasks"`
}

func (po *Project) MarshalJSON() ([]byte, error) {
	m := &projectJSON{
		Id:       po.Id,
		Name:     po.Name,
		Color:    po.Color,
		ViewMode: po.ViewMode.String(),
		Kind:     po.Kind.String(),
		Tasks:    po.Tasks,
	}
	return json.Marshal(m)
}

func (po *Project) UnmarshalJSON(data []byte) error {
	var (
		m1 projectJSON
		m2 projectTaskJSON
	)
	err := json.Unmarshal(data, &m1)
	if err != nil {
		return err
	}
	if m1.Id != "" {
		po.Id = m1.Id
		po.Name = m1.Name
		po.Color = m1.Color
		po.Kind = kindFromString(m1.Kind)
		po.ViewMode = viewModeFromString(m1.ViewMode)
		po.GroupId = m1.GroupId
		po.Closed = m1.Closed
		po.Tasks = m1.Tasks
		return nil
	}
	err = json.Unmarshal(data, &m2)
	if err != nil {
		return err
	}
	po.Id = m2.Project.Id
	po.Name = m2.Project.Name
	po.Color = m2.Project.Color
	po.Kind = kindFromString(m2.Project.Kind)
	po.ViewMode = viewModeFromString(m2.Project.ViewMode)
	po.GroupId = m2.Project.GroupId
	po.Closed = m2.Project.Closed
	po.Tasks = m2.Tasks
	return nil
}
