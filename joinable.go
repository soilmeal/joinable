package joinable

import (
	"errors"
	"fmt"
	"sync"

	"go.uber.org/atomic"
)

var idGeneratorForJoinable atomic.Uint64

// Joinable 은 goroutine wrapping 클래스입니다.
// 다른 언어에서 쓰이던 join 개념이 없는 것이 개인적으로 불편해서 만들었습니다.
type Joinable struct {
	id uint64

	err      error
	errMutex sync.Mutex

	opt *Option

	mutex *sync.Mutex
	cond  *sync.Cond
}

// NewJoinable 메소드는 새로운 Joinable 인스턴스를 생성합니다.
func NewJoinable(runnable Runnable) *Joinable {
	mutex := new(sync.Mutex)
	return &Joinable{
		id: idGeneratorForJoinable.Inc(),

		opt: NewOption(runnable, false),

		mutex: mutex,
		cond:  sync.NewCond(mutex),
	}
}

// NewJoinableWithFunc 메소드는 새로운 Joinable 인스턴스를 생성합니다.
// Runnable interface 를 구현하기 귀찮으신 분들을 위해서 준비했습니다.
func NewJoinableWithFunc(runnable func()) *Joinable {
	mutex := new(sync.Mutex)
	return &Joinable{
		id: idGeneratorForJoinable.Inc(),

		opt: NewOptionWithFunc(runnable, false),

		mutex: mutex,
		cond:  sync.NewCond(mutex),
	}
}

// NewJoinableWithOption 메소드는 새로운 Joinable 인스턴스를 생성합니다.
// Option 클래스에서 자세한 값을 확인해주세요.
func NewJoinableWithOption(opt *Option) *Joinable {
	mutex := new(sync.Mutex)
	return &Joinable{
		id: idGeneratorForJoinable.Inc(),

		opt: opt,

		mutex: mutex,
		cond:  sync.NewCond(mutex),
	}
}

// Start 메소드는 실제 goroutine 을 실행시킵니다.
func (joinable *Joinable) Start() {
	joinable.runGoroutine()
}

// HasError 메소드는 goroutine 실행 시 발생한 에러가 있는지를 확인합니다.
func (joinable *Joinable) HasError() bool {
	joinable.errMutex.Lock()
	defer joinable.errMutex.Unlock()

	return joinable.err != nil
}

// Error 메소드는 goroutine 실행 시 발생한 에러를 반환합니다.
func (joinable *Joinable) Error() error {
	joinable.errMutex.Lock()
	defer joinable.errMutex.Unlock()

	return joinable.err
}

// Join 메소드는 goroutine 실행 종료를 대기합니다.
// 생성시 인자 혹은 Option 객체에 설정한 실제 Runnable 구현이 정상적으로 끝나지 않는다면
// 무한 대기하게 되므로, 꼭 Runnable 구현을 확인해주시기 바랍니다.
func (joinable *Joinable) Join() {
	joinable.mutex.Lock()
	joinable.cond.Wait()
	joinable.mutex.Unlock()
}

// runGoroutine 메소드는 실제로 goroutine 을 실행합니다.
func (joinable *Joinable) runGoroutine() {
	go func() {
		joinable.clearError()

		if joinable.opt.ShouldCatchPanic {
			defer func() {
				err := recover()
				joinable.setError(err)
			}()
		}

		joinable.opt.Runnable.Run()

		joinable.mutex.Lock()
		joinable.cond.Signal()
		joinable.mutex.Unlock()
	}()
}

// clearError 메소드는 현재 가지고 있는 error 를 정리합니다.
func (joinable *Joinable) clearError() {
	joinable.errMutex.Lock()
	defer joinable.errMutex.Unlock()

	joinable.err = nil
}

// setError 는 인자로 받은 error 를 Joinable 의 err 에 설정합니다.
func (joinable *Joinable) setError(err interface{}) {
	joinable.errMutex.Lock()
	defer joinable.errMutex.Unlock()

	if err == nil {
		return
	}

	switch x := err.(type) {
	case string:
		joinable.err = errors.New(x)
	case error:
		joinable.err = x
	case fmt.Stringer:
		joinable.err = fmt.Errorf("occur error with this value - err = %s", x.String())
	default:
		joinable.err = errors.New("occur error with unknown value")
	}
}
