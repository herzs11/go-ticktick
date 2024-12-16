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
	var m projectJSON
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}
	po.Id = m.Id
	po.Name = m.Name
	po.Color = m.Color
	po.Kind = kindFromString(m.Kind)
	po.ViewMode = viewModeFromString(m.ViewMode)
	po.GroupId = m.GroupId
	po.Closed = m.Closed
	po.Tasks = m.Tasks
	return nil
}
