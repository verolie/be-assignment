package model

type DataStatusResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type StatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
