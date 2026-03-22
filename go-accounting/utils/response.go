package utils

type SuccessResponse struct {
	Errno     int         `json:"Errno"`
	Errmsg    string      `json:"Errmsg"`
	ErrorCode string      `json:"ErrorCode"`
	Data      interface{} `json:"data"`
}

type ErrorResponse struct {
	Errno     int    `json:"Errno"`
	Errmsg    string `json:"Errmsg"`
	ErrorCode string `json:"ErrorCode"`
}

func SuccessResponseWithData(data interface{}, message ...string) SuccessResponse {
	msg := "success"
	if len(message) > 0 {
		msg = message[0]
	}
	return SuccessResponse{
		Errno:     0,
		Errmsg:    msg,
		ErrorCode: "OK",
		Data:      data,
	}
}

func ErrorResponseWithData(errno int, errmsg string, errorCode string) ErrorResponse {
	return ErrorResponse{
		Errno:     errno,
		Errmsg:    errmsg,
		ErrorCode: errorCode,
	}
}
