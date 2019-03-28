package notify

import (
	log "github.com/sirupsen/logrus"

	gosxnotifier "github.com/deckarep/gosx-notifier"
)

//SendMessage to OSX notification centre
func SendMessage(msg string) {
	note := gosxnotifier.NewNotification(msg)
	err := note.Push()
	if err != nil {
		log.Error("could not send")
	}
}
