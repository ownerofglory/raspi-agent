package domain

import "io"

type RecordingResult interface {
	SaveTo(writer io.Writer) error
}
