package constants

type SystemFiletype string

func (e SystemFiletype) IsValid() bool {
	switch e {
	case AvatarFiletype, FileInChatFiletype, VoiceFiletype, CircleFiletype:
		return true
	}

	return false
}

func (e SystemFiletype) String() string {
	return string(e)
}

const (
	AvatarFiletype     SystemFiletype = "avatar"
	FileInChatFiletype SystemFiletype = "file_in_chat"
	VoiceFiletype      SystemFiletype = "voice"
	CircleFiletype     SystemFiletype = "circle"
)
