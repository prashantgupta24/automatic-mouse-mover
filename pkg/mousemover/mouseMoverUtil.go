package mousemover

import (
	"os"

	"github.com/go-vgo/robotgo"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func getLogger(m *MouseMover, doWriteToFile bool) *log.Logger {
	logger := log.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}

	if doWriteToFile {
		_, err := os.Stat(logDir)
		if err != nil {
			if os.IsNotExist(err) {
				err = os.Mkdir(logDir, os.ModePerm)
				if err != nil {
					log.Fatalf("error creating dir: %v", err)
				}
			}
		}

		logFile, err := os.OpenFile(logDir+"/"+logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		logger.SetOutput(logFile)
		m.logFile = logFile
	}

	return logger
}

func moveAndCheck(movePixel int, mouseMoveSuccessCh chan bool) {
	currentX, currentY := robotgo.GetMousePos()
	moveToX := currentX + movePixel
	moveToY := currentY + movePixel
	robotgo.MoveMouse(moveToX, moveToY)

	//check if mouse moved. Sometimes mac users need to give
	//extra permission for controlling the mouse
	movedX, movedY := robotgo.GetMousePos()
	if movedX == currentX && movedY == currentY {
		mouseMoveSuccessCh <- false
	} else {
		mouseMoveSuccessCh <- true
	}
}

func (m *MouseMover) isRunning() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.runningStatus
}
func (m *MouseMover) updateRunningStatus(isRunning bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.runningStatus = isRunning
}
