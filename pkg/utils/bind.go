package utils

import "encoding/json"

func BindFromStruct(from any, to any) error {
	temporaryVariable, _ := json.Marshal(from)
	if err := json.Unmarshal(temporaryVariable, &to); err != nil {
		return err
	}

	return nil
}
