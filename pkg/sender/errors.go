package sender

type NoWebsocketConnectionError struct {
}

func (NoWebsocketConnectionError) Error() string {
	return "no websocket connection"
}

var NoWebsocketConnectionErr NoWebsocketConnectionError = NoWebsocketConnectionError{}
