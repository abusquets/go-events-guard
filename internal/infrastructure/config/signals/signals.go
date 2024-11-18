package signals

type SignalKey string

const (
	UserCreated SignalKey = "user.created"
	UserUpdated SignalKey = "user.updated"
)
