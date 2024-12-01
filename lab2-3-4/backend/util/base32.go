package util

import "encoding/base32"

func ConvertStringToByteSlice(str string) ([]byte, error) {
	bytes, err := base32.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func ConvertByteSliceToString(bytes []byte) string {
	return base32.StdEncoding.EncodeToString(bytes)
}
