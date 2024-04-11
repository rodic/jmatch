package jmatch

type UnexpectedEndOfInputErr struct{}

func (e UnexpectedEndOfInputErr) Error() string {
	return "invalid JSON. Unexpected end of JSON input"
}
