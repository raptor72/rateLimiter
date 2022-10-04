package handlers

type LoginRequest struct {
	Login string `json:"login"`
}

type IpRequest struct {
	Login string `json:"ip"`
}

type PasswordRequest struct {
	Login string `json:"password"`
}

type UnionRequest struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
	Ip       string `json:"ip,omitempty"`
}

type RequestStruct struct {
	FiledName string
	JsonName  string
}

var Name string
var JS string

type TTT struct {
	Name string `json: wdwd`
}
