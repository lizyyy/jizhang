package utils

// Response 统一响应结构
type Response struct {
	Errno     int         `json:"Errno"`
	Errmsg    string      `json:"Errmsg"`
	ErrorCode string      `json:"ErrorCode"`
	Data      interface{} `json:"data"`
}

// SuccessResponse 成功响应
func SuccessResponse(data interface{}, message string) Response {
	if message == "" {
		message = "success"
	}
	return Response{
		Errno:     0,
		Errmsg:    message,
		ErrorCode: "OK",
		Data:      data,
	}
}

// ErrorResponse 错误响应
func ErrorResponse(errno int, errmsg string, errorCode string) Response {
	return Response{
		Errno:     errno,
		Errmsg:    errmsg,
		ErrorCode: errorCode,
		Data:      nil,
	}
}
