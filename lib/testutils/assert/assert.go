package assert

import (
	"testing"
)

func NoError(t *testing.T, err error, msg... string) {
	if err != nil {
		t.Fatalf("Unexpected error '%v'. %s", err, CombinedMsg(msg...))
	}
}

func CombinedMsg(messages... string) string {
	msg := ""
	if len(messages) > 0 {
		msg = messages[0]
	}

	if len(messages) > 1 {
		for _, message := range messages[1:] {
			msg += " " + message
		}
	}
	return msg
}
