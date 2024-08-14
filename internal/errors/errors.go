package errors

import "errors"

func UnsupportedCodingTypeError(message string) error {
	return errors.New("Unsupported coding type. Message :" + message)
}
