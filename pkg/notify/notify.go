package notify

import (
	"fmt"

	gosxnotifier "github.com/deckarep/gosx-notifier"
)

//SendMessage to OSX notification centre
func SendMessage(msg string) {
	note := gosxnotifier.NewNotification(msg)
	err := note.Push()
	if err != nil {
		fmt.Println("could not send notification")
	}
}
