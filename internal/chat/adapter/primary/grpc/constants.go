package grpcadapter

const (
	logMsgCleanupSessionFailure = "failed to cleanup session"
	logFieldRoom                = "room"
	logFieldUser                = "user"
	logFieldError               = "error"

	welcomeMessageFormat = "Bem-vindo %s!"
	noticeJoinedFormat   = "%s entrou na sala"
	noticeLeftFormat     = "%s saiu da sala"

	errMsgClientCanceled      = "client canceled stream"
	errMsgSessionExists       = "session already established"
	errMsgJoinPayloadRequired = "join payload required"
	errMsgChatPayloadRequired = "chat payload required"
	errMsgJoinRequired        = "join required before sending messages"
	errMsgLeavePayloadReq     = "leave payload required"
	errMsgNoActiveSession     = "no active session"
	errMsgInvalidPayload      = "invalid payload"
)

const (
	zeroUnixTimestamp = 0
)
