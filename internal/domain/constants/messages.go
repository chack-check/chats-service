package constants

type MessageTypes string

const (
	TextMessageType   MessageTypes = "text"
	EventMessageType  MessageTypes = "event"
	CallMessageType   MessageTypes = "call"
	VoiceMessageType  MessageTypes = "voice"
	CircleMessageType MessageTypes = "circle"
)
