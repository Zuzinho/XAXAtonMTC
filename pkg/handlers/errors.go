package handlers

type InvalidRequestParamsError struct {
}

func (InvalidRequestParamsError) Error() string {
	return "invalid params 'user_id' or 'access' value"
}

type NoUserIDParamError struct {
}

func (NoUserIDParamError) Error() string {
	return "no param 'user_id'"
}

var (
	InvalidRequestParamsErr = InvalidRequestParamsError{}
	NoUserIDParamErr        = NoUserIDParamError{}
)
