package entity

type DeviceState struct {
	Id    string
	State []State
}

type State struct {
	Tag   string
	Value string
}
