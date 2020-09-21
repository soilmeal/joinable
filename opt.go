package joinable

// Option 은 Joinable 을 위한 옵션 클래스입니다.
type Option struct {
	ShouldCatchPanic bool
	Runnable         Runnable
}

// NewOption 은 Runnable 을 받아 새로운 Option 을 생성합니다.
func NewOption(runnable Runnable, shouldCatchPanic bool) *Option {
	return &Option{
		Runnable:         runnable,
		ShouldCatchPanic: shouldCatchPanic,
	}
}

// NewOptionWithFunc 은 func() 을 받아 새로운 Option 을 생성합니다.
func NewOptionWithFunc(runnable func(), shouldCatchPanic bool) *Option {
	return &Option{
		Runnable:         WrapToRunnable(runnable),
		ShouldCatchPanic: shouldCatchPanic,
	}
}
