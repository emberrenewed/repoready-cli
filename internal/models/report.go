package models

type SystemReport struct {
	OS           string
	Architecture string
	Tools        []ToolStatus
}
