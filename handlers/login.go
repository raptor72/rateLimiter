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
