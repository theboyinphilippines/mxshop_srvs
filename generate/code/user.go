package code

const (
	// ErrUserNotFound -   404:  User not found
	ErrUserNotFound int = iota + 100401
	// ErrUserInvalidParam - 401 :  User invalid parameter
	ErrUserInvalidParam
)

type Response struct {
	codeName int
	httpCode int
	msg      string
}

func register(codeName int, httpCode int, msg string) *Response {
	return &Response{
		codeName: codeName,
		httpCode: httpCode,
		msg:      msg,
	}

}
