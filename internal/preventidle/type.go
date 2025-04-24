package preventidle

type IdleController struct {
	StopChan chan struct{}
	Running  bool
}
