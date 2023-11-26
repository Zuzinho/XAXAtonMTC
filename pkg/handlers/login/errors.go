package login

type InvalidSmsCodeError struct {
}

func (InvalidSmsCodeError) Error() string {
	return "invalid sms code value"
}

var InvalidSmsCodeErr InvalidSmsCodeError = InvalidSmsCodeError{}
