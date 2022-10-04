package mousemover

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/prashantgupta24/activity-tracker/pkg/activity"
	"github.com/prashantgupta24/activity-tracker/pkg/tracker"
)

type TestMover struct {
	suite.Suite
	activityTracker *tracker.Instance
	heartbeatCh     chan *tracker.Heartbeat
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

	suite.heartbeatCh = make(chan *tracker.Heartbeat)
}

//Run once before each test
func (suite *TestMover) SetupTest() {
	instance = nil
}
func (suite *TestMover) TestAppStart() {
	t := suite.T()
	mouseMover := GetInstance()
	mouseMover.run(suite.heartbeatCh, suite.activityTracker)
	time.Sleep(time.Millisecond * 500) //wait for app to start
	assert.True(t, mouseMover.state.isRunning(), "app should have started")
}
func (suite *TestMover) TestSingleton() {
	t := suite.T()

	mouseMover1 := GetInstance()
	mouseMover1.run(suite.heartbeatCh, suite.activityTracker)

	time.Sleep(time.Millisecond * 500)

	mouseMover2 := GetInstance()
	assert.True(t, mouseMover2.state.isRunning(), "instance should have started")
}

func (suite *TestMover) TestLogFile() {
	t := suite.T()
	mouseMover := GetInstance()
	logFileName := "test1"

	getLogger(mouseMover, true, logFileName)

	filePath := logDir + "/" + logFileName
	assert.FileExists(t, filePath, "log file should exist")
	os.RemoveAll(filePath)
}
func (suite *TestMover) TestSystemSleepAndWake() {
	t := suite.T()
	mouseMover := GetInstance()

	state := &state{
		override: &override{
			valueToReturn: true,
		},
	}
	mouseMover.state = state
	heartbeatCh := make(chan *tracker.Heartbeat)

	mouseMover.run(heartbeatCh, suite.activityTracker)
	time.Sleep(time.Millisecond * 500) //wait for app to start
	assert.True(t, mouseMover.state.isRunning(), "instance should have started")
	assert.False(t, mouseMover.state.isSystemSleeping(), "machine should not be sleeping")

	//fake a machine-sleep activity
	machineSleepActivityMap := make(map[activity.Type][]time.Time)
	var sleepTimeArray []time.Time
	sleepTimeArray = append(sleepTimeArray, time.Now())
	machineSleepActivityMap[activity.MachineSleep] = sleepTimeArray
	heartbeatCh <- &tracker.Heartbeat{
		WasAnyActivity: true,
		ActivityMap:    machineSleepActivityMap,
	}
	time.Sleep(time.Millisecond * 500) //wait for it to be registered
	assert.True(t, mouseMover.state.isSystemSleeping(), "machine should be sleeping now")

	//assert app is sleeping
	heartbeatCh <- &tracker.Heartbeat{
		WasAnyActivity: false,
	}

	time.Sleep(time.Millisecond * 500) //wait for it to be registered
	assert.True(t, time.Time.IsZero(state.getLastMouseMovedTime()), "should be default but is ", state.getLastMouseMovedTime())
	assert.Equal(t, state.getDidNotMoveCount(), 0, "should be 0")

	//fake a machine-wake activity
	machineWakeActivityMap := make(map[activity.Type][]time.Time)
	var wakeTimeArray []time.Time
	wakeTimeArray = append(wakeTimeArray, time.Now())
	machineWakeActivityMap[activity.MachineWake] = wakeTimeArray
	heartbeatCh <- &tracker.Heartbeat{
		WasAnyActivity: true,
		ActivityMap:    machineWakeActivityMap,
	}

	time.Sleep(time.Millisecond * 500) //wait for it to be registered
	assert.False(t, mouseMover.state.isSystemSleeping(), "machine should be awake now")
}

func (suite *TestMover) TestMouseMoveSuccess() {
	t := suite.T()
	mouseMover := GetInstance()

	state := &state{
		override: &override{
			valueToReturn: true,
		},
	}
	mouseMover.state = state
	heartbeatCh := make(chan *tracker.Heartbeat)

	mouseMover.run(heartbeatCh, suite.activityTracker)
	time.Sleep(time.Millisecond * 500) //wait for app to start
	assert.True(t, state.isRunning(), "instance should have started")
	assert.False(t, state.isSystemSleeping(), "machine should not be sleeping")
	assert.True(t, time.Time.IsZero(state.getLastMouseMovedTime()), "should be default")
	assert.Equal(t, state.getDidNotMoveCount(), 0, "should be 0")

	heartbeatCh <- &tracker.Heartbeat{
		WasAnyActivity: false,
	}

	time.Sleep(time.Millisecond * 500) //wait for it to be registered
	assert.False(t, time.Time.IsZero(state.getLastMouseMovedTime()), "should be default but is ", state.getLastMouseMovedTime())
}

func (suite *TestMover) TestMouseMoveFailure() {
	t := suite.T()
	mouseMover := GetInstance()

	state := &state{
		override: &override{
			valueToReturn: false,
		},
	}
	mouseMover.state = state
	heartbeatCh := make(chan *tracker.Heartbeat)

	mouseMover.run(heartbeatCh, suite.activityTracker)
	time.Sleep(time.Millisecond * 500) //wait for app to start
	assert.True(t, state.isRunning(), "instance should have started")
	assert.False(t, state.isSystemSleeping(), "machine should not be sleeping")
	assert.True(t, time.Time.IsZero(state.getLastMouseMovedTime()), "should be default")
	assert.Equal(t, state.getDidNotMoveCount(), 0, "should be 0")
	assert.True(t, state.getLastErrorTime().IsZero(), "should be default")

	heartbeatCh <- &tracker.Heartbeat{
		WasAnyActivity: false,
	}

	time.Sleep(time.Millisecond * 500) //wait for it to be registered
	assert.True(t, time.Time.IsZero(state.getLastMouseMovedTime()), "should be default but is ", state.getLastMouseMovedTime())
	assert.NotEqual(t, state.getDidNotMoveCount(), 0, "should not be 0")
	assert.NotEqual(t, state.getLastErrorTime(), 0, "should not be 0")
}
