package eventbus

type EventBus interface {
	Publisher
	Consumer
	Close() error
}
