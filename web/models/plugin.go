package models

type PluginVO struct {
	Index      int     `json:"index,omitempty"`
	ID         string  `json:"id,omitempty"`
	FileID     string  `json:"fileID,omitempty"`
	Name       string  `json:"name,omitempty"`
	Desc       string  `json:"desc,omitempty"`
	Inputs     []Input `json:"inputs,omitempty"`
	Code       string  `json:"code,omitempty"`
	Language   string  `json:"language,omitempty"`
	VersionNum int     `json:"versionNum,omitempty"`
	Status     string  `json:"status,omitempty"`
	CreateTime string  `json:"createTime,omitempty"`
	UpdateTime string  `json:"updateTime,omitempty"`
}

func (pv *PluginVO) GetData() map[string]string {
	result := map[string]string{}
	for _, input := range pv.Inputs {
		result[input.Name] = input.Value
	}
	return result
}

type Input struct {
	Name     string `json:"name,omitempty"`
	Required bool   `json:"required"`
	Desc     string `json:"desc,omitempty"`
	Value    string `json:"value,omitempty"`
}

func (i *Input) GetMultiNodeConfigByName(name string) string {
	if i.Name == name {
		return i.Value
	}
	return ""
}
