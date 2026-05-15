package models

type ToolStatus struct {
	Name                string
	Command             string
	Version             string
	Installed           bool
	Required            bool
	InstallInstructions string
}
