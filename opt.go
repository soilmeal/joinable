package joinable

// Option 은 Joinable 을 위한 옵션 클래스입니다.
type Option struct {
	ShouldRecoverPanic bool
	Runnable           Runnable
}

// NewOption 은 Runnable 을 받아 새로운 Option 을 생성합니다.
func NewOption(runnable Runnable, shouldRecoverPanic bool) *Option {
	return &Option{
		Runnable:           runnable,
		ShouldRecoverPanic: shouldRecoverPanic,
	}
}

// NewOptionWithFunc 은 func() 을 받아 새로운 Option 을 생성합니다.
func NewOptionWithFunc(runnable func(), shouldRecoverPanic bool) *Option {
	return &Option{
		Runnable:           WrapToRunnable(runnable),
		ShouldRecoverPanic: shouldRecoverPanic,
	}
}
