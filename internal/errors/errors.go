package errors

import (
	"errors"
)

func UnsupportedCodingTypeError(message string) error {
	return errors.New("Unsupported coding type. Message: " + message)
}

func InvalidHeader(message string) error {
	return errors.New("Invalid header. Message: " + message)
}

func InvalidOutDataPointer(message string) error {
	return errors.New("Invalid output data pointer. Message: " + message)
}
