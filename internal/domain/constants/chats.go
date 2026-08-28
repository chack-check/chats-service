package constants

type ActionTypes string

const (
	WritingActionType         = "writing"
	AudioRecordingActionType  = "audio_recording"
	AudioSendingActionType    = "audio_sending"
	CircleRecordingActionType = "circle_recording"
	CircleSendingActionType   = "circle_sending"
	FilesSendingActionType    = "files_sending"
)

type ChatTypes string

var (
	UserChatType          ChatTypes = "user"
	GroupChatType         ChatTypes = "group"
	SavedMessagesChatType ChatTypes = "saved_messages"
)
