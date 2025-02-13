package models

import "errors"

var ErrorNoRowsAffected = errors.New("no rows affected")
var ErrorUserNotFound = errors.New("user not found")
var ErrorListNotFound = errors.New("list not found")
