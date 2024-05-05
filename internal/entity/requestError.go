package entity

type RequestError struct {
	StatusCode int
	Err        error
}
