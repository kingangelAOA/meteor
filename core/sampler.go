package core


type Sampler interface {
	Request(SamplerConfig) Result
}

type SamplerConfig interface {
	GetConfig() interface{}
	Close()
}
