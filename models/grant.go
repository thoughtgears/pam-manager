package models

type Grant struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Requester     string   `json:"requester"`
	Duration      int64    `json:"duration"`
	Justification string   `json:"justification"`
	State         string   `json:"state"`
	Roles         []string `json:"roles"`
}
