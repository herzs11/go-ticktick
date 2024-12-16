package tasks

import (
	"time"
)

const TASK_ENDPOINT = "/open/v1/task"

type Priority int
type Status int

const (
	None   Priority = 0
	Low             = 1
	Medium          = 3
	High            = 5
)

const (
	Normal Status = iota
	Completed
)

func (s Status) String() string {
	switch s {
	case Normal:
		return "Normal"
	case Completed:
		return "Completed"
	default:
		return ""
	}
}

func (p Priority) String() string {
	switch p {
	case None:
		return "None"
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	default:
		return ""
	}
}

type checklistItemJSON struct {
	Id            string `json:"id,omitempty"`
	Title         string `json:"title,omitempty"`
	Status        int    `json:"status,omitempty"`
	CompletedTime string `json:"completedTime,omitempty"`
	IsAllDay      bool   `json:"isAllDay,omitempty"`
	SortOrder     int    `json:"sortOrder,omitempty"`
	StartDate     int64  `json:"startDate,omitempty"`
	TimeZone      string `json:"timeZone,omitempty"`
}

type ChecklistItem struct {
	Id            string
	Title         string
	Status        Status
	CompletedTime time.Time
	IsAllDay      bool
	SortOrder     int
	StartDate     time.Time
	TimeZone      string
}

type taskJSON struct {
	Id             string          `json:"id,omitempty"`
	ProjectId      string          `json:"projectId,omitempty"`
	Title          string          `json:"title,omitempty"`
	IsAllDay       bool            `json:"isAllDay,omitempty"`
	CompletedTime  string          `json:"completedTime,omitempty"`
	Content        string          `json:"content,omitempty"`
	Desc           string          `json:"desc,omitempty"`
	DueDate        string          `json:"dueDate,omitempty"`
	ChecklistItems []ChecklistItem `json:"items,omitempty"`
	Priority       int             `json:"priority,omitempty"`
	Reminders      []string        `json:"reminders,omitempty"`
	Tags           []string        `json:"tags,omitempty"`
	RepeatFlag     string          `json:"repeatFlag,omitempty"`
	SortOrder      int64           `json:"sortOrder,omitempty"`
	StartDate      string          `json:"startDate,omitempty"`
	Status         int             `json:"status,omitempty"`
	TimeZone       string          `json:"timeZone,omitempty"`
}

type Task struct {
	Id             string
	ProjectId      string
	Title          string
	IsAllDay       bool
	CompletedTime  time.Time
	Content        string
	Desc           string
	DueDate        time.Time
	ChecklistItems []ChecklistItem
	Priority       Priority
	Reminders      []string
	RepeatFlag     string
	SortOrder      int64
	StartDate      time.Time
	Status         Status
	TimeZone       string
}
