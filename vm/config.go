package vm

type Config struct {
	StackPoolSize int
	StackCapacity int
}

func (cfg *Config) setDefaults() {
	if cfg.StackPoolSize == 0 {
		cfg.StackPoolSize = 10
	}
}
