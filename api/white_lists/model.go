package white_lists

type createWhiteListRequest struct {
	Address string `json:"address"`
}

type updateWhiteListRequest struct {
	Address *string `json:"address"`
}

type createWhiteListResponse struct {
	Id int `json:"id"`
}

type whiteListResult struct {
	Items []*WhiteListModel `json:"items"`
}

type successResult struct {
	Success bool `json:"success"`
}
