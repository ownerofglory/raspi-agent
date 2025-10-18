package domain

import "io"

type VoiceAssistantRequest struct {
	Audio io.Reader
}

type VoiceAssistantResult struct {
	Audio io.Reader
}
