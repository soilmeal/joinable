package joinable

// Runnable 은 Joinable 이 실행 가능한 작업을 나타내는 interface 입니다.
type Runnable interface {
	Run()
}

// WrapToRunnable 은 Runnable 을 구현하는 struct 를
// 만드는 것보다, 간단하게 func() 을 사용하고 싶어하시는 분들을 위한
// 함수입니다.
func WrapToRunnable(runnable func()) Runnable {
	return &wrapper{
		runnable: runnable,
	}
}

// func() 을 맴버로 가지고 Run() 호출 시 실행하는 struct 입니다.
type wrapper struct {
	runnable func()
}

// 실제 func() 을 실행합니다.
func (w *wrapper) Run() {
	if w.runnable == nil {
		return
	}

	w.runnable()
}
