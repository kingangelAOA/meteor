package models

type UserRes struct {
	Token string `json:"token"`
}

type UserInfoRes struct {
	Roles        []string `json:"roles"`
	Introduction string   `json:"introduction"`
	Avatar       string   `json:"avatar"`
	Name         string   `json:"name"`
}
