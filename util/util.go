package util

import (
	"encoding/json"
)

func UnmarshalTo[T any](data []byte) (dest *T, err error) {
    dest = new(T)
	err = json.Unmarshal(data, dest)
	if err != nil {
		return
	}
    return
}