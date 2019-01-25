package errors

//PswInvalidateErr password invalidate error
type PswInvalidateErr struct {
}

//Error implement error interface
func (e *PswInvalidateErr) Error() string {
	return "password is not right"
}

//NewPswInvalidateErr return NewPswInvalidateErr
func NewPswInvalidateErr() *PswInvalidateErr {
	return &PswInvalidateErr{}
}
