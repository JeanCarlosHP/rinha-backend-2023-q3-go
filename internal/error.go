package internal

import "errors"

var ErrPeopleExists = errors.New("people alredy exists")
var ErrPeopleNotFound = errors.New("people not found")
var ErrInternalServer = errors.New("internal server error")
