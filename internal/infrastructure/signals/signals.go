package signals

type SignalsBus interface {
	Subscribe(topic string, callback Callback)
	Emit(topic string, args ...interface{}) error
	ProcessQueue() error
	AfterTransaction(topic string, args ...interface{}) error
	EventChannel() chan struct{}
}
