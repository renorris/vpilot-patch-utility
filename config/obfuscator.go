package config

import (
	"crypto/md5"
	"encoding/base64"
)

var ConfigObfuscatorKey = configObfuscatorKey()

func configObfuscatorKey() []byte {
	const guid = "5575ac09-f2de-4a1e-808b-e3398e17f8bf"
	sum := md5.Sum([]byte(guid))

	key := make([]byte, 24)
	copy(key[:16], sum[:])
	copy(key[16:], sum[:8])
	return key
}

func encrypt(plaintext []byte, key []byte) (ciphertext []byte, err error) {
	return tripleDESEncrypt(plaintext, key)
}

func decrypt(ciphertext []byte, key []byte) (plaintext []byte, err error) {
	return tripleDESDecrypt(ciphertext, key)
}

func ObfuscateToBase64(plaintext []byte, key []byte) (ciphertextBase64 []byte, err error) {
	var ciphertext []byte
	if ciphertext, err = encrypt(plaintext, key); err != nil {
		return
	}

	ciphertextBase64 = []byte(base64.StdEncoding.EncodeToString(ciphertext))
	return
}

func DeobfuscateFromBase64(ciphertextBase64 []byte, key []byte) (plaintext []byte, err error) {
	var ciphertext []byte
	if ciphertext, err = base64.StdEncoding.DecodeString(string(ciphertextBase64)); err != nil {
		return
	}

	if plaintext, err = decrypt(ciphertext, key); err != nil {
		return
	}

	return
}
