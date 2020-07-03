package dto

import "time"

type AdminInfoOutput struct {
	ID   int `json:"id"`
	Name string `json:"name"`
	LoginTime time.Time   `json:"loginTime"`
	Avatar string     `json:"avatar"`
	Introduction  string       `json:"introduction"`
	Roles  []string      `json:"roles"`
}
