package models

type SharedLinkResponse struct {
	Id    string     `json:"id"`
	Roles []string   `json:"roles"`
	Link  LinkShared `json:"link"`
}
