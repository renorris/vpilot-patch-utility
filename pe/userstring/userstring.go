package userstring

import (
	"bytes"
	"encoding/binary"
	"errors"
	"golang.org/x/text/encoding/unicode"
	"io"
	"os"
	"unicode/utf16"
)

// ReadUserString reads a UTF-16 string from the #US heap at the specified
// file offset, then returns the UTF-8 representation of that string.
func ReadUserString(file *os.File, fileOffset int64) (str string, err error) {
	if _, err = file.Seek(fileOffset, io.SeekStart); err != nil {
		return
	}

	lengthHeader := make([]byte, 4)
	if _, err = io.ReadFull(file, lengthHeader); err != nil {
		return
	}

	var dataLength, headerSize int
	if dataLength, headerSize, err = DecodeLength([4]byte(lengthHeader)); err != nil {
		return
	}

	if dataLength%2 != 1 {
		err = errors.New("user string data length should be odd")
		return
	}

	// Seek to the byte right after the length header
	if _, err = file.Seek(fileOffset+int64(headerSize), io.SeekStart); err != nil {
		return
	}

	strData := make([]byte, dataLength)
	if _, err = io.ReadFull(file, strData); err != nil {
		return
	}

	// The last byte is a terminal byte. Ignore it.
	utf16Str := make([]uint16, (len(strData)-1)/2)
	for i := 0; i < len(strData)-1; i += 2 {
		utf16Str[i/2] = binary.LittleEndian.Uint16(strData[i : i+2])
	}

	// Decode UTF-16 to runes
	runes := utf16.Decode(utf16Str)

	// Convert runes to UTF-8 string
	utf8String := string(runes)

	str = utf8String
	return
}

// WriteUserString writes a string on the #US heap at the specified
// file offset using the provided UTF-8 encoded string `str`.
func WriteUserString(file *os.File, fileOffset int64, str string) (err error) {
	// Convert UTF-8 characters into UTF-16 string
	encoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()

	var utf16Bytes []byte
	if utf16Bytes, err = encoder.Bytes([]byte(str)); err != nil {
		return
	}

	// Check if the last byte needs to be set:
	//
	// https://ecma-international.org/wp-content/uploads/ECMA-335_6th_edition_june_2012.pdf
	// II.24.2.4 #US and #Blob heaps
	//
	// "Strings in the #US (user string) heap are encoded using 16-bit Unicode encodings. The count on each
	// string is the number of bytes (not characters) in the string. Furthermore, there is an additional terminal
	// byte (so all byte counts are odd, not even). This final byte holds the value 1 if and only if any UTF16
	// character within the string has any bit set in its top byte, or its low byte is any of the following: 0x01–
	// 0x08, 0x0E–0x1F, 0x27, 0x2D, 0x7F. Otherwise, it holds 0. The 1 signifies Unicode characters that
	// require handling beyond that normally provided for 8-bit encoding sets."

	setTerminalBit := false
	for i := 1; i < len(utf16Bytes); i += 2 {
		if utf16Bytes[i] > 0 {
			setTerminalBit = true
			break
		}
		switch utf16Bytes[i-1] {
		case 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
			0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15,
			0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D,
			0x1E, 0x1F, 0x27, 0x2D, 0x7F:
			setTerminalBit = true
			break
		}
	}

	if setTerminalBit {
		utf16Bytes = append(utf16Bytes, []byte{0x01}...)
	} else {
		utf16Bytes = append(utf16Bytes, []byte{0x00}...)
	}

	// Encode the length of utf16Bytes
	var header []byte
	if header, err = EncodeLength(len(utf16Bytes)); err != nil {
		return
	}

	// Seek to the file offset
	if _, err = file.Seek(fileOffset, io.SeekStart); err != nil {
		return
	}

	// Write the header
	if _, err = io.Copy(file, bytes.NewReader(header)); err != nil {
		return err
	}

	// Write the utf16 string bytes
	if _, err = io.Copy(file, bytes.NewReader(utf16Bytes)); err != nil {
		return err
	}

	return nil
}
