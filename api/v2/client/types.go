package client

type loginResponse struct {
	Token          string `json:"token"`
	UserId         string `json:"userId"`
	UserCode       string `json:"userCode"`
	Username       string `json:"Username"`
	TeamPro        bool   `json:"teamPro"`
	ProStartDate   string `json:"proStartDate"`
	ProEndDate     string `json:"proEndDate"`
	SubscribeType  string `json:"subscribeType"`
	SubscribeFreq  string `json:"subscribeFreq"`
	NeedSubscribe  bool   `json:"needSubscribe"`
	Freq           string `json:"freq"`
	InboxId        string `json:"inboxId"`
	TeamUser       bool   `json:"teamUser"`
	ActiveTeamUser bool   `json:"activeTeamUser"`
	FreeTrial      bool   `json:"freeTrial"`
	Pro            bool   `json:"pro"`
	Ds             bool   `json:"ds"`
}
