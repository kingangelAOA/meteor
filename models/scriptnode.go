package models

type ScriptNode map[string][]BaseScript

type BaseScript struct {
	ID   string
	Name string
	Code string
}
