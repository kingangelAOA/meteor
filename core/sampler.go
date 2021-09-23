package core


type Sampler interface {
	Request() Result
}

type SamplerConfig interface {
	GetConfig() interface{}
	Close()
}
