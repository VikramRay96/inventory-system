package dto

// ErrorResponse Dto
type ErrorResponseDto struct {
	StatusCode StatusCode `json:"status_code"`
	Message    string     `json:"message"`
}

func (p *ErrorResponseDto) SetError(err StatusCode) {
	p.StatusCode = GetStatusDetails(err).StatusCode
	p.Message = GetStatusDetails(err).Message
}
