package types

import "encoding/json"

type NullInt struct {
	Set   bool
	Value *int64
}

func (n *NullInt) UnmarshalJSON(data []byte) error {
	n.Set = true

	if string(data) == "null" {
		n.Value = nil
		return nil
	}

	var v int64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	n.Value = &v
	return nil
}
