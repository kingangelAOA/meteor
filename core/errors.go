package core

const (
	FormURLEncodeParamsMarshalToByteError = "form urlencode params to byte error"
)

type HttpError struct {
	Err string
}

func (he *HttpError) Error() string {
	return he.Err
}
