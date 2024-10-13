package userstring

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeLengthSingleByte(t *testing.T) {
	header, err := EncodeLength(51)
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x33}, header)
}

func TestEncodeLengthTwoBytes(t *testing.T) {
	header, err := EncodeLength(273)
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x81, 0x11}, header)
}

func TestEncodeLengthFourBytes(t *testing.T) {
	header, err := EncodeLength(17000)
	assert.Nil(t, err)
	assert.Equal(t, []byte{0b11000000, 0b0, 0b1000010, 0b01101000}, header)

	header, err = EncodeLength(536870911)
	assert.Nil(t, err)
	assert.Equal(t, []byte{0b11011111, 0xFF, 0xFF, 0xFF}, header)
}

func TestEncodeLengthErr(t *testing.T) {
	_, err := EncodeLength(-1)
	assert.NotNil(t, err)

	_, err = EncodeLength(536870912)
	assert.NotNil(t, err)
}

func TestDecodeLengthSingleByte(t *testing.T) {
	length, headerSize, err := DecodeLength([4]byte{0x33, 0x00, 0x00, 0x00})
	assert.Nil(t, err)
	assert.Equal(t, 51, length)
	assert.Equal(t, 1, headerSize)
}

func TestDecodeLengthTwoBytes(t *testing.T) {
	length, headerSize, err := DecodeLength([4]byte{0x81, 0x11, 0x00, 0x00})
	assert.Nil(t, err)
	assert.Equal(t, 273, length)
	assert.Equal(t, 2, headerSize)
}

func TestDecodeLengthFourBytes(t *testing.T) {
	length, headerSize, err := DecodeLength([4]byte{0b11000000, 0b0, 0b1000010, 0b01101000})
	assert.Nil(t, err)
	assert.Equal(t, 17000, length)
	assert.Equal(t, 4, headerSize)

	length, headerSize, err = DecodeLength([4]byte{0b11011111, 0xFF, 0xFF, 0xFF})
	assert.Nil(t, err)
	assert.Equal(t, 536870911, length)
	assert.Equal(t, 4, headerSize)
}
