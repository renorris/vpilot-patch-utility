package userstring

import "errors"

// DecodeLength decodes the length of a #US or #Blob string.
// https://ecma-international.org/wp-content/uploads/ECMA-335_6th_edition_june_2012.pdf
// II.24.2.4 #US and #Blob heaps
func DecodeLength(header [4]byte) (length int, headerSize int, err error) {
	if header[0]>>7 == 0 {
		// Length is the 7 LSBs of header[0]
		length = int(header[0])
		headerSize = 1
		return
	}

	if header[0]>>6 == 0b10 {
		// Length is (header[0] bbbbbb2 << 8 + header[1])
		length = (int(header[0]&0b00111111) << 8) + int(header[1])
		headerSize = 2
		return
	}

	if header[0]>>5 == 0b110 {
		// Length is (header[0] bbbbb2 << 24 + header[1] << 16 + header[2] << 8 + header[3])
		length = (int(header[0]&0b00011111) << 24) +
			(int(header[1]) << 16) +
			(int(header[2]) << 8) +
			(int(header[3]))
		headerSize = 4
		return
	}

	length = -1
	err = errors.New("invalid length header")
	return
}

// EncodeLength encodes the length of a #US or #Blob string.
// https://ecma-international.org/wp-content/uploads/ECMA-335_6th_edition_june_2012.pdf
// II.24.2.4 #US and #Blob heaps
func EncodeLength(length int) (header []byte, err error) {
	if length < 0 {
		err = errors.New("cannot encode negative length")
		return
	}

	// If the length fits into 7 bits, encode the header as a single byte
	// Encode a single byte if the length is < 128
	if length < 128 {
		header = []byte{byte(length)}
		return
	}

	// If the length fits into 14 bits, encode into 2 bytes
	// Set 0b10 header for the first byte to indicate that we're
	// using 2 bytes
	if length < 16384 {
		header = []byte{0b10000000 | byte(length>>8), byte(length)}
		return
	}

	// Max length is 2^29 bits
	// If the length fits into 29 bits, encode into 4 bytes
	// Set 0b110 header to indicate that we're using 4 bytes
	if length < 536870912 {
		header = []byte{0b11000000 | byte(length>>24), byte(length >> 16), byte(length >> 8), byte(length)}
		return
	}

	length = -1
	err = errors.New("cannot encode length at or above 536870912 (29 bits max)")
	return
}
