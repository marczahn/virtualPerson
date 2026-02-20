package output

import "fmt"

type SourceTag string

const (
	SourceBIO    SourceTag = "BIO"
	SourceDRIVES SourceTag = "DRIVES"
	SourceMIND   SourceTag = "MIND"
)

func FormatTaggedLine(tag SourceTag, message string) string {
	return fmt.Sprintf("[%s] %s", tag, message)
}
