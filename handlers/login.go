package handlers

type LoginRequest struct {
	Login string `json:"login"`
}

type IPRequest struct {
	Login string `json:"ip"`
}

type PasswordRequest struct {
	Login string `json:"password"`
}

type UnionRequest struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
	IP       string `json:"ip,omitempty"`
}

type RequestStruct struct {
	FiledName string
	JSONName  string
}

// var Name string
// var JS string
