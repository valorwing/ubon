package ubon

import "ubon/internal/serialize"

func MarshalUBON(input any) ([]byte, error) {

	return serialize.Serialize(input)
}

func UnmarshalUBON(data []byte, target *any) error {

	return nil
}
