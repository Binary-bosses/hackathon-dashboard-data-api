package server

type CreateHackathonData struct {
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	StartTime     interface{} `json:"startTime"`
	EndTime       interface{} `json:"endTime"`
	HackathonPass string      `json:"hackathonPass"`
	Winner        string      `json:"winner,omitempty"`
}

type EditHackathonData struct {
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	StartTime     interface{} `json:"startTime"`
	EndTime       interface{} `json:"endTime"`
	HackathonPass string      `json:"hackathonPass"`
	Teams         []Team      `json:"teams"`
	Winner        string      `json:"winner,omitempty"`
}

type Team struct {
	Name    string   `json:"name"`
	Idea    string   `json:"idea,omitempty"`
	Members []string `json:"members"`
}

type HackathonData struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	StartTime   interface{} `json:"startTime"`
	EndTime     interface{} `json:"endTime"`
	Winner      string      `json:"winner,omitempty"`
	Teams       []HackTeam  `json:"teams"`
}

type HackTeam struct {
	Name string `json:"name"`
	Idea string `json:"idea"`
}

type HackathonHighLevelData struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	StartTime   interface{} `json:"startTime"`
	EndTime     interface{} `json:"endTime"`
	Winner      string      `json:"winner,omitempty"`
}

type HackathonEditPass struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

type Status struct {
	Status string `json:"status"`
}
