package util

const(
	SUCCESS = 0
	SUCCESS_MSG = "ok"
)

type Response struct {
	ErrNo  int64 `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	Data   interface{} `json:"data"`
}

func SuccessResponse(data interface{}) Response {
	return Response{ErrNo:SUCCESS,ErrMsg:SUCCESS_MSG,Data:data}
}

func FailResponse(errNo int64, errMsg string, data interface{}) Response {
	return Response{errNo,errMsg,data}
}

