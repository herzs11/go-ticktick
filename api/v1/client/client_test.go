package client

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/herzs11/go-ticktick/api/v1/types/project"
	"github.com/herzs11/go-ticktick/api/v1/types/tasks"
	"github.com/joho/godotenv"
	"github.com/zalando/go-keyring"
)

func getAuthenticatedClient(t *testing.T) *TickTickClient {
	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Log("Could not load .env file")
	}
	testClientID := os.Getenv("TT_CLIENT_ID")
	if testClientID == "" {
		t.Fatal("TickTickClient id must be set in the TT_CLIENT_ID environment variable")
	}
	testClientSecret := os.Getenv("TT_CLIENT_SECRET")
	testRedirectUri := os.Getenv("TT_REDIRECT_URI")
	c := NewTickTickClient(testClientID, testClientSecret, testRedirectUri)
	err = c.Authenticate()
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestOauth2GetToken(t *testing.T) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Log("Could not load .env file")
	}
	testClientID := os.Getenv("TT_CLIENT_ID")
	if testClientID == "" {
		t.Fatal("TickTickClient id must be set in the TT_CLIENT_ID environment variable")
	}
	testClientSecret := os.Getenv("TT_CLIENT_SECRET")
	testRedirectUri := os.Getenv("TT_REDIRECT_URI")
	c := NewTickTickClient(testClientID, testClientSecret, testRedirectUri)
	err = c.getAuthorizationCode()
	if err != nil {
		t.Fatal("Error getting authorization code: ", err)
	}
	if c.authorizationCode == "" {
		t.Fatal("Could not get authorization code")
	}
	fmt.Println("Got authorization code from listener server", c.authorizationCode)
	if checkPort("8080") {
		t.Fatal("Port still in use")
	}

	err = c.getOauthToken()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(c.token.AccessToken)
}

func TestTokenFileStore(t *testing.T) {
	c := getAuthenticatedClient(t)
	expTS := time.Second * time.Duration(c.token.ExpiresIn)
	c.token.ExpiresTime = time.Now().Add(expTS).Unix()
	err := storeTokenFile(c.token)
	if err != nil {
		log.Fatal(err)
	}

	retrievedToken, err := getTokenFromFile()
	if err != nil {
		log.Fatal(err)
	}

	if retrievedToken.AccessToken != c.token.AccessToken {
		log.Fatal("Retrieved token does not equal created token")
	}

	if retrievedToken.ExpiresTime != c.token.ExpiresTime {
		log.Fatal("Retrieved expiry time != created expires time")
	}
}

func TestTokenKeyringStore(t *testing.T) {
	c := getAuthenticatedClient(t)

	data, err := json.Marshal(&c.token)
	if err != nil {
		t.Fatal(err)
	}

	service := "go-ticktick-test"
	err = keyring.Set(service, c.ClientId, string(data))
	if err != nil {
		t.Fatal(err)
	}
	secret, err := keyring.Get(service, c.ClientId)
	if err != nil {
		t.Fatal(err)
	}
	retrievedData := []byte(secret)
	var retrievedToken oauthToken
	err = json.Unmarshal(retrievedData, &retrievedToken)
	if err != nil {
		t.Fatalf("Unable to unmarshal data: %s", err.Error())
	}

	if c.token.AccessToken != retrievedToken.AccessToken {
		t.Fatalf("Expected token %s, got %s", c.token.AccessToken, retrievedToken.AccessToken)
	}
	if c.token.ExpiresTime != retrievedToken.ExpiresTime {
		t.Fatalf("Expected expTime %d, got %d", c.token.ExpiresTime, retrievedToken.ExpiresTime)
	}
}

func TestOauth2Client_Authenticate(t *testing.T) {
	c := getAuthenticatedClient(t)

	token, err := getTokenFromKeyring(c.ClientId)
	if err != nil {
		t.Fatalf("Unable to get token from keyring: %s", err.Error())
	}
	if token.AccessToken != c.token.AccessToken {
		t.Fatalf("Expected token %s from keyring, got %s", c.token.AccessToken, token.AccessToken)
	}
}

func TestCreateNewProject(t *testing.T) {
	c := getAuthenticatedClient(t)
	tests := [][]string{
		[]string{"Go Test5", "#00FF00"},
		[]string{"Go Test6", "#FFFF00"},
		[]string{"Go Test7", "#00FFFF"},
		[]string{"Go Test8", "#0000FF"},
	}
	for _, s := range tests {

		pIn := project.Project{Name: s[0], Color: s[1]}
		err := c.CreateNewProject(&pIn)
		if err != nil {
			t.Fatal(err)
		}
		if pIn.Name != s[0] {
			t.Fatalf("Expected %s, got %s", s[0], pIn.Name)
		}
		if pIn.Color != s[1] {
			t.Fatalf("Expected %s got %s", s[1], pIn.Color)
		}
		fmt.Printf("%+v", pIn)
	}
}

func TestGetAllProjects(t *testing.T) {
	c := getAuthenticatedClient(t)

	projs, err := c.GetAllProjects(false)
	if err != nil {
		t.Fatal(err)
	}
	if len(projs) != 1 {
		t.Fatalf("Expected 4 projects, got %d", len(projs))
	}
	for _, p := range projs {
		fmt.Printf("Id: %s\nName: %s\nColor%s\n\n", p.Id, p.Name, p.Color)
	}
}

// TODO: Delete request on project returns 500, sent an email to support
func TestDeleteProjectById(t *testing.T) {
	c := getAuthenticatedClient(t)
	p := project.Project{Name: "TestCreate"}
	err := c.CreateNewProject(&p)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	err = c.DeleteProjectById(p.Id)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(5 * time.Second)

	p2, err := c.GetProjectById(p.Id, false)
	if err != nil {
		t.Fatal(err)
	}
	if p2 != nil {
		t.Fatal("Expected unsuccessful pull of project")
	}
}

func TestGetProjectById(t *testing.T) {
	c := getAuthenticatedClient(t)

	p := project.Project{Name: "TestCreationNew"}
	err := c.CreateNewProject(&p)
	pRes, err := c.GetProjectById(p.Id, false)
	if err != nil {
		t.Fatal(err)
	}

	if pRes.Id != p.Id {
		t.Fatalf("Expected id %s, got %s", p.Id, pRes.Id)
	}
	if pRes.Color != p.Color {
		t.Fatalf("Expected color %s, got %s", p.Color, pRes.Color)
	}
	c.DeleteProjectById(p.Id)
}

func TestUpdateProject(t *testing.T) {
	c := getAuthenticatedClient(t)

	p := project.Project{Name: "TestProjectUpdate"}
	err := c.CreateNewProject(&p)
	if err != nil {
		t.Fatal(err)
	}

	p.Name = "UpdatedProjectName"
	p.ViewMode = project.Kanban
	err = c.UpdateProject(&p)
	if err != nil {
		t.Fatal(err)
	}

	projRes, err := c.GetProjectById(p.Id, false)
	if err != nil {
		t.Fatal(err)
	}

	if p.Name != projRes.Name {
		t.Fatalf("Expected %s, got %s", p.Name, projRes.Name)
	}
	if p.ViewMode != projRes.ViewMode {
		t.Fatalf("Expected %s, got %s", p.ViewMode.String(), projRes.ViewMode.String())
	}
}

func TestGetProjectWithTasks(t *testing.T) {
	c := getAuthenticatedClient(t)
	proj := project.Project{Name: "testProject"}
	err := c.CreateNewProject(&proj)
	if err != nil {
		t.Fatal(err)
	}
	t1 := tasks.Task{
		Title:     "newTask",
		ProjectId: proj.Id,
	}
	t2 := tasks.Task{
		Title:     "newTask2",
		ProjectId: proj.Id,
	}
	err = c.CreateTask(&t1)
	if err != nil {
		t.Fatal(err)
	}
	proj.Tasks = append(proj.Tasks, t1)
	err = c.CreateTask(&t2)
	if err != nil {
		t.Fatal(err)
	}
	proj.Tasks = append(proj.Tasks, t2)

	projT, err := c.GetProjectById(proj.Id, true)
	if err != nil {
		t.Fatal(err)
	}

	if len(projT.Tasks) != 2 {
		t.Fatalf("Expected 2 related tasks, got %d", len(projT.Tasks))
	}

}

func TestGetInbox(t *testing.T) {
	c := getAuthenticatedClient(t)
	inbox, err := c.GetInbox()
	if err != nil {
		t.Fatal(err)
	}
	if inbox.Id != "inbox" {
		t.Fatalf("Expected inbox, got %s", inbox.Id)
	}
	if inbox.Name != "Inbox" {
		t.Fatalf("Expected Inbox, got %s", inbox.Name)
	}
	if len(inbox.Tasks) == 0 {
		t.Fatal("Expected tasks to be populated")
	}
}

func TestGetTaskById(t *testing.T) {
	c := getAuthenticatedClient(t)

	t1 := tasks.Task{
		Title:    "testTask",
		Priority: tasks.High,
		ChecklistItems: []tasks.ChecklistItem{
			tasks.ChecklistItem{
				Title:  "subItem1",
				Status: 0,
			},
			tasks.ChecklistItem{
				Title:  "subItem2",
				Status: 0,
			},
		},
	}
	err := c.CreateTask(&t1)
	if err != nil {
		t.Fatal(err)
	}

	task, err := c.GetTaskById(t1.Id)
	if err != nil {
		t.Fatal(err)
	}
	if len(task.ChecklistItems) != 2 {
		t.Fatalf("Task does not have 2 subtasks attached, got %d", len(task.ChecklistItems))
	}
	if task.Title != "testTask" {
		t.Fatalf("Expected task to be titled %s, got %s", "testTask", task.Title)
	}
	if task.Priority.String() != "High" {
		t.Fatalf("Expected task priority to be High, got %s", task.Priority.String())
	}
	if err := c.DeleteTask(task); err != nil {
		t.Fatal(err)
	}
}

func TestGetTask(t *testing.T) {
	c := getAuthenticatedClient(t)
	t1 := tasks.Task{
		Title: "testTask",
	}
	err := c.CreateTask(&t1)
	if err != nil {
		t.Fatal(err)
	}

	task := tasks.Task{
		Id: t1.Id,
	}
	err = c.GetTask(&task)
	if err != nil {
		t.Fatal(err)
	}

	if task.Id != t1.Id {
		t.Fatalf("Expected id %s, got id %s", t1.Id, task.Id)
	}
	if task.Title != "testTask" {
		t.Fatalf("Expected %s, got %s", "testTask", task.Title)
	}

	err = c.DeleteTask(&task)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCompleteTask(t *testing.T) {
	c := getAuthenticatedClient(t)
	task := tasks.Task{
		Title:  "TestTask",
		Status: tasks.Normal,
	}
	err := c.CreateTask(&task)
	if err != nil {
		t.Fatal(err)
	}

	err = c.CompleteTask(&task)
	if err != nil {
		t.Fatal(err)
	}

	if task.Status.String() != "Completed" {
		t.Fatal("Failed to update status on object")
	}

	task2, err := c.GetTaskById(task.Id)
	if err != nil {
		t.Fatal(err)
	}
	if task2.Status.String() != "Completed" {
		t.Fatal("Task is not completed on server side")
	}
	err = c.DeleteTask(task2)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteTask(t *testing.T) {
	c := getAuthenticatedClient(t)
	inbox, err := c.GetInbox()
	if err != nil {
		t.Fatal(err)
	}
	task := tasks.Task{
		Title:    "TestTask",
		Priority: tasks.High,
	}
	err = c.CreateTask(&task)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Created task with id %s", task.Id)
	time.Sleep(1 * time.Second)

	err = c.DeleteTask(&task)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)

	inbox, err = c.GetInbox()
	if err != nil {
		t.Fatal(err)
	}

	for _, tsk := range inbox.Tasks {
		if tsk.Id == task.Id {
			t.Fatalf("Found test task with id %s", tsk.Id)
		}
	}
}

func TestCreateTask(t *testing.T) {
	c := getAuthenticatedClient(t)
	now := time.Now().Local()
	task := tasks.Task{
		Title:    "TEST TASK 1",
		DueDate:  now,
		Priority: tasks.High,
		TimeZone: "America/New_York",
	}

	err := c.CreateTask(&task)
	if err != nil {
		t.Fatalf("Error creating task: %s", err.Error())
	}

	if task.Id == "" {
		t.Fatal("Task missing Id after creation")
	}
	if task.Title != "TEST TASK 1" {
		t.Fatalf("Expected title '%s', got '%s'", "TEST TASK 1", task.Title)
	}
	if task.DueDate.YearDay() != now.YearDay() || task.DueDate.Hour() != now.Hour() || task.DueDate.Minute() != now.Minute() || task.DueDate.Second() != now.Second() {
		t.Fatalf("Expected due date %s, got %s", now.String(), task.DueDate.String())
	}
	err = c.DeleteTask(&task)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateTask(t *testing.T) {
	c := getAuthenticatedClient(t)
	now := time.Now().Local()
	task := tasks.Task{
		Title:    "TEST TASK 2",
		DueDate:  now,
		Priority: tasks.High,
		TimeZone: "America/New_York",
	}

	err := c.CreateTask(&task)
	if err != nil {
		t.Fatal(err)
	}

	task.Priority = tasks.Medium
	err = c.UpdateTask(&task)
	if err != nil {
		t.Fatal(err)
	}

	t1, err := c.GetTaskById(task.Id)
	if t1.Priority.String() != "Medium" {
		t.Fatalf("Expected Priority %s, got %s", task.Priority.String(), t1.Priority.String())
	}
}
