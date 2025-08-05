package schemas

import "encoding/json"

// WriteOnlyString a write only string is not serializable to json.
// in the SCIM RFC it has a mutability of writeOnly.
// This increases security to really ensure this is never sent to a client.
type WriteOnlyString string

func (s *WriteOnlyString) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}

func (s *WriteOnlyString) UnmarshalJSON(bytes []byte) error {
	var str string
	err := json.Unmarshal(bytes, &str)
	*s = WriteOnlyString(str)
	return err
}

func (s *WriteOnlyString) String() string {
	return string(*s)
}
