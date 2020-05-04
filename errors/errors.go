package errors

//PswInvalidateErr password invalidate error
type PswInvalidateErr struct {
}

//Error implement error interface
func (e *PswInvalidateErr) Error() string {
	return "password is not right"
}

//PasswordErr password invalidate error
var PasswordErr = NewPswInvalidateErr()

//NewPswInvalidateErr return NewPswInvalidateErr
func NewPswInvalidateErr() *PswInvalidateErr {
	return &PswInvalidateErr{}
}
