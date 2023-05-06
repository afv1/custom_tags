package confidential

type Handler func(input string) string

type Config struct {
	Tags     []string
	Mappers  map[string]Handler
	Handlers []Handler
}

type StructMask interface {
	Proceed(input any) any
}

var StructMasker StructMask

func InitStructMask(cfg *Config) {
	StructMasker = newSM()
}
