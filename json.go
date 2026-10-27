package librus

import (
	"bytes"
	"encoding/json"
	"strconv"
)

// flexInt decodes integers that Librus sends either as JSON numbers or as
// strings (e.g. "Id": "4133586"). null and "" decode to 0.
type flexInt int

func (f *flexInt) UnmarshalJSON(b []byte) error {
	b = bytes.Trim(b, `"`)
	if len(b) == 0 || string(b) == "null" {
		*f = 0
		return nil
	}
	n, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}
	*f = flexInt(n)
	return nil
}

// ref is the {"Id": ..., "Url": ...} object Librus uses to point at other resources.
type ref struct {
	ID flexInt `json:"Id"`
}

// flexString decodes identifiers that Librus sends either as strings or as
// JSON numbers (e.g. attendance ids: 5 or "t209235"). null decodes to "".
type flexString string

func (f *flexString) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		*f = ""
		return nil
	}
	if len(b) > 0 && b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		*f = flexString(s)
		return nil
	}
	*f = flexString(b)
	return nil
}
