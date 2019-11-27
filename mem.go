package util

func TrimBytes(b []byte, ln int) []byte {
	if len(b) > ln {
		return b[:ln]
	} else {
		return b
	}
}
