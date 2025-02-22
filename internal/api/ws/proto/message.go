package proto

import "google.golang.org/protobuf/proto"

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (m *Message) UnmarshalBinary(data []byte) error {
	return proto.Unmarshal(data, m)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (m *Message) MarshalBinary() ([]byte, error) {
	return proto.Marshal(m)
}
