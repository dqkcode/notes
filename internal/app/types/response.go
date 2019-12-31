package types

type (
	Response struct {
		Code  string
		Data  interface{}
		Error string
	}
)

var (
	CodeSuccess = "0000"
	UserAlreadyExist = "0001"
	ErrorDB = "0002"
)
