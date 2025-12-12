package gox

type (
	ActE  func() error
	Act0  func()
	Act1  func(any)
	Act2  func(a any, b any)
	Act3  func(a any, b any, c any)
	Act4  func(a any, b any, c any, d any)
	Act5  func(a any, b any, c any, d any, e any)
	FuncE func() (any, error)
	Func0 func() any
	Func1 func(a any) any
	Func2 func(a any, b any) any
	Func3 func(a any, b any, c any) any
	Func4 func(a any, b any, c any, d any) any
	Func5 func(a any, b any, c any, d any, e any) any
	Func6 func(a any, b any, c any, d any, e any, f any) any
	Func7 func(a any, b any, c any, d any, e any, f any, g any) any
)

type (
	Judge0  func() bool
	JudgeE  func() (bool, error)
	Judge1  func(a any) bool
	Judge1E func(a any) (bool, error)
	Judge2  func(a any, b any) bool
	Judge2E func(a any, b any) (bool, error)
	Judge3  func(a any, b any, c any) bool
	Judge3E func(a any, b any, c any) (bool, error)
)
