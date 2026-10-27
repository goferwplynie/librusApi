package librus

import (
	"bytes"
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
