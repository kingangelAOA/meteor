package models

type NodeVO struct {
	Index         int     `json:"index,omitempty"`
	ID            string  `json:"id,omitempty"`
	Describe      string  `json:"describe,omitempty"`
	Name          string  `json:"name,omitempty"`
	Type          string  `json:"type,omitempty"`
	Ico           string  `json:"ico,omitempty"`
	PluginID      string  `json:"pluginID,omitempty"`
	DefaultInputs []Input `json:"defaultInputs,omitempty"`
	CreateTime    string  `json:"createTime,omitempty"`
	UpdateTIme    string  `json:"updateTime,omitempty" `
}
