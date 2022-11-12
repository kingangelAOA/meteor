package service

import (
	"web/db"
	"web/models"
)

func GetScripts() ([]models.BaseScript, error) {
	ss, err := db.GetScripts()
	if err != nil {
		return nil, err
	}
	var svs []models.BaseScript
	for _, s := range ss {
		svs = append(svs, models.BaseScript{
			Type: s.Type,
			Name: s.Name,
			Code: s.Code,
		})
	}
	return svs, nil
}

func GetScriptList() []models.ScriptListVO {
	return []models.ScriptListVO{}
}
