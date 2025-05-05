package preventidle

type IdleController struct {
	StopChan chan struct{}
	Running  bool
}

type SimulateAction struct {
	inputType  string
	actionName string
}

type LastInputInfo struct {
	cbSize uint32
	dwTime uint32
}

type KeyboardInput struct {
	wVk         uint16
	wScan       uint16
	dwFlags     uint32
	time        uint32
	dwExtraInfo uintptr
}

type MouseInput struct {
	dx          int32
	dy          int32
	mouseData   uint32
	dwFlags     uint32
	time        uint32
	dwExtraInfo uintptr
}

type input struct {
	Type uint32
	_    [4]byte // padding
	Ki   KeyboardInput
	Mi   MouseInput
}
