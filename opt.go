package joinable

// Option 은 Joinable 을 위한 옵션 클래스입니다.
type Option struct {
	// ShouldRecoverPanic 은 panic 이 발생했을 때 recover() 를 자동으로 호출할 지를 결정합니다.
	// 만약 직접 호출하기를 원하신다면, false 로 설정해주시기 바랍니다.
	ShouldRecoverPanic bool

	// Joinable 에 전달될 실제 Runnable 입니다.
	// Runnable 을 구현하는 struct 를 직접 작성하기 귀찮다면
	// func() 를 전달하는 방법도 제공하니 참고하시기 바랍니다.
	Runnable Runnable
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
