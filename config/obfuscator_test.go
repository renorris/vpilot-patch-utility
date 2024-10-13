package config

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVPilotObfuscatorKey(t *testing.T) {
	input := hex.EncodeToString(ConfigObfuscatorKey)
	expectedResult := "9bf2bf12df6e3fb45c0db019e1982a109bf2bf12df6e3fb4"

	assert.Equal(t, expectedResult, input)
	return
}

func TestDeobfuscateFromBase64(t *testing.T) {
	input := "pZ9u441bE4a2NCGgqMxKNvwAIy0qEA+AwXB8c3sV90c="
	expectedResult := "http://status.vatsim.net/"
	key := ConfigObfuscatorKey

	plaintext, err := DeobfuscateFromBase64([]byte(input), key)
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, string(plaintext))
}

func TestObfuscateToBase64(t *testing.T) {
	input := "http://status.vatsim.net/"
	expectedResult := "pZ9u441bE4a2NCGgqMxKNvwAIy0qEA+AwXB8c3sV90c="
	key := ConfigObfuscatorKey

	ciphertextBase64, err := ObfuscateToBase64([]byte(input), key)
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, string(ciphertextBase64))
}
