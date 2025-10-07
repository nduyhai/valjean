package ports

type Worker interface {
	Submit(func())
	Shutdown()
}
