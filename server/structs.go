package server

type CreateHackathonData struct {
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	StartTime     interface{} `json:"startTime"`
	EndTime       interface{} `json:"endTime"`
	Teams         []Team      `json:"teams"`
	HackathonPass string      `json:"hackathonPass"`
}

type Team struct {
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

type HackathonData struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	StartTime   interface{} `json:"startTime"`
	EndTime     interface{} `json:"endTime"`
	Teams       []string    `json:"teams"`
}

type HackathonEditPass struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

type Status struct {
	Status string `json:"status"`
}
