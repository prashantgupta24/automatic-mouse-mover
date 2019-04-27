package mousemover

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/prashantgupta24/activity-tracker/pkg/tracker"
)

type TestMover struct {
	suite.Suite
	activityTracker *tracker.Instance
	heartbeatCh chan *tracker.Heartbeat
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TestMover))
}
//Run once before all tests
func (suite *TestMover) SetupSuite() {
		heartbeatInterval := 60 
	workerInterval := 10

	suite.activityTracker = &tracker.Instance{
		HeartbeatInterval: heartbeatInterval,
		WorkerInterval:    workerInterval,
	}

	suite.heartbeatCh= make(chan *tracker.Heartbeat)
}

//Run once before each test
func (suite *TestMover) SetupTest() {
	instance = nil
}

func (suite *TestMover) TestSingleton() {
	t := suite.T()

	mouseMover1 := GetInstance()
	mouseMover1.run(suite.heartbeatCh, suite.activityTracker)

	time.Sleep(time.Millisecond * 500)

	mouseMover2 := GetInstance()
	assert.True(t, mouseMover2.isRunning(), "instance should have started")
}
func (suite *TestMover) TestAppStartAndStop() {
	t := suite.T()
	mouseMover := GetInstance()
	mouseMover.run(suite.heartbeatCh, suite.activityTracker)
	time.Sleep(time.Millisecond * 500) //wait for app to start
	assert.True(t, mouseMover.isRunning(), "app should have started")
	mouseMover.Quit()
	time.Sleep(time.Millisecond * 1000) //wait for app to stop
	assert.False(t, mouseMover.isRunning(), "app should have stopped")
}
