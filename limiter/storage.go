package limiter

type Storage interface {
	Increment(key string) (int, error)
	Reset(key string) error
	IsBlocked(key string) (bool, error)
	Block(key string, duration int) error
}
