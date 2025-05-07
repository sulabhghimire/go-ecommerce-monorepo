package helper

import "errors"

var (
	NOT_AUTHORIZED_ERROR = errors.New("you are not authorized to peform this action")
)