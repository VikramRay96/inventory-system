package dto

type ResponseDto struct {
	StatusCode StatusCode  `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}
