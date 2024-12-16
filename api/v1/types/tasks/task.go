package tasks

import (
	"encoding/json"
	"time"
)

func (c *ChecklistItem) toJSON() checklistItemJSON {
	return checklistItemJSON{
		Id:            c.Id,
		Title:         c.Title,
		Status:        int(c.Status),
		CompletedTime: convertLocalTime(c.CompletedTime),
		IsAllDay:      c.IsAllDay,
		SortOrder:     c.SortOrder,
		StartDate:     c.StartDate.UTC().Unix(),
		TimeZone:      c.TimeZone,
	}
}

func (c *ChecklistItem) MarshalJSON() ([]byte, error) {
	cj := c.toJSON()
	return json.Marshal(&cj)
}

func (c *ChecklistItem) UnmarshalJSON(data []byte) error {
	cj := checklistItemJSON{}
	err := json.Unmarshal(data, &cj)
	if err != nil {
		return err
	}

	c.Id = cj.Id
	c.Title = cj.Title
	c.Status = Status(cj.Status)
	c.CompletedTime = convertUTCString(cj.CompletedTime)
	c.IsAllDay = cj.IsAllDay
	c.SortOrder = cj.SortOrder
	c.StartDate = time.Unix(cj.StartDate, 0).Local()
	c.TimeZone = cj.TimeZone
	return nil
}

func (t *Task) toJSON() taskJSON {
	tj := taskJSON{
		Id:             t.Id,
		ProjectId:      t.ProjectId,
		Title:          t.Title,
		IsAllDay:       t.IsAllDay,
		CompletedTime:  convertLocalTime(t.CompletedTime),
		Content:        t.Content,
		Desc:           t.Desc,
		DueDate:        convertLocalTime(t.DueDate),
		ChecklistItems: t.ChecklistItems,
		Priority:       int(t.Priority),
		Reminders:      t.Reminders,
		RepeatFlag:     t.RepeatFlag,
		SortOrder:      t.SortOrder,
		StartDate:      convertLocalTime(t.StartDate),
		Status:         int(t.Status),
		TimeZone:       t.TimeZone,
	}
	return tj
}

func (t *Task) MarshalJSON() ([]byte, error) {
	tj := t.toJSON()
	return json.Marshal(&tj)
}

func (t *Task) UnmarshalJSON(data []byte) error {
	tj := taskJSON{}
	err := json.Unmarshal(data, &tj)
	if err != nil {
		return err
	}

	t.Id = tj.Id
	t.ProjectId = tj.ProjectId
	t.Title = tj.Title
	t.IsAllDay = tj.IsAllDay
	t.CompletedTime = convertUTCString(tj.CompletedTime)
	t.Content = tj.Content
	t.Desc = tj.Desc
	t.DueDate = convertUTCString(tj.DueDate)
	t.ChecklistItems = tj.ChecklistItems
	t.Priority = Priority(tj.Priority)
	t.Reminders = tj.Reminders
	t.RepeatFlag = tj.RepeatFlag
	t.SortOrder = tj.SortOrder
	t.StartDate = convertUTCString(tj.StartDate)
	if tj.Status == 2 {
		t.Status = Completed
	} else {
		t.Status = Normal
	}

	return nil
}
