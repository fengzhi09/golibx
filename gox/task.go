package gox

import (
	"runtime"
	"sync"
	"time"
)

type TaskArgs = any
type TaskID string
type TaskState = string
type StepFunc func(taskId TaskID, args TaskArgs) error

const (
	StateNone      TaskState = ""
	StateWaiting   TaskState = "waiting"
	StateCompleted TaskState = "completed"
	StateError     TaskState = "failed"
	StateExpired   TaskState = "expired"
)

func InitTaskRunner(num int) {
	for i := 0; i < num; i++ {
		go startRunner()
	}
}

func StartTask(args TaskArgs, steps []StepFunc, expireTime time.Duration) TaskID {
	return addAsyncTask(newAsyncTask(args, steps, expireTime))
}

func (taskID TaskID) Exists() bool {
	taskStateMutex.Lock()
	defer taskStateMutex.Unlock()
	_, exists := taskState[taskID]
	return exists
}

func (taskID TaskID) Remove() {

}
func (taskID TaskID) IsWaiting() bool {
	return getState(taskID) == StateWaiting
}

func (taskID TaskID) IsExpired() bool {
	return getState(taskID) == StateExpired
}

func (taskID TaskID) IsCompleted() bool {
	return getState(taskID) == StateCompleted
}

func (taskID TaskID) IsFailed() bool {
	return getState(taskID) == StateError
}

var (
	taskStateMutex sync.Mutex
	taskPool       sync.Pool
	taskChan       = make(chan *asyncTask, runtime.NumCPU())
	taskState      = make(map[TaskID]TaskState)
)

type asyncTask struct {
	id        TaskID
	args      TaskArgs
	steps     []StepFunc
	expiredAt int64
}

func setState(taskId TaskID, state TaskState) {
	taskStateMutex.Lock()
	defer taskStateMutex.Unlock()
	taskState[taskId] = state
}

func getState(taskId TaskID) TaskState {
	taskStateMutex.Lock()
	defer taskStateMutex.Unlock()

	resultState, exists := taskState[taskId]
	return IfElse(exists, StateNone, resultState).(TaskState)
}

func newAsyncTask(args TaskArgs, steps []StepFunc, expireTime time.Duration) *asyncTask {
	var expiredAt int64
	if expireTime > 0 {
		expiredAt = time.Now().Add(expireTime).UnixNano()
	} else {
		expiredAt = -1
	}

	t := taskPool.Get()
	if t == nil {
		return &asyncTask{newTaskId(), args, steps, expiredAt}
	} else {
		task := t.(*asyncTask)
		(*task).args = args
		(*task).steps = steps
		(*task).id = newTaskId()
		(*task).expiredAt = expiredAt
		return task
	}

}

func newTaskId() TaskID {
	return TaskID(NewOIDHex())
}

func addAsyncTask(task *asyncTask) TaskID {
	go func() {
		taskChan <- task
	}()
	taskId := (*task).id
	setState(taskId, StateWaiting)
	return taskId
}

func startRunner() {
	var id TaskID
	var err error
	for {
		task := <-taskChan
		id = (*task).id
		if ((*task).expiredAt > 0 && time.Now().UnixNano() < (*task).expiredAt) || (*task).expiredAt < 0 {
			for _, f := range (*task).steps {
				err = f((*task).id, (*task).args)
			}
			if err != nil {
				setState(id, StateError)
				taskPool.Put(task)
			} else {
				setState(id, StateCompleted)
				taskPool.Put(task)
			}
		} else {
			setState(id, StateExpired)
			taskPool.Put(task)
		}
	}
}
