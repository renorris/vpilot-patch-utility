package config

import (
	"crypto/cipher"
	"crypto/des"
)

// Decrypts ciphertext using Triple DES in ECB mode
func tripleDESDecrypt(ciphertext []byte, key []byte) (plaintext []byte, err error) {
	// Create a new DES cipher
	var block cipher.Block
	if block, err = des.NewTripleDESCipher(key); err != nil {
		return nil, err
	}

	// Initialize byte array for the plaintext
	plaintext = make([]byte, len(ciphertext))

	// decrypt in ECB mode
	for i := 0; i < len(ciphertext); i += block.BlockSize() {
		block.Decrypt(plaintext[i:i+block.BlockSize()], ciphertext[i:i+block.BlockSize()])
	}

	return pkcs7strip(plaintext, block.BlockSize())
}

// Encrypts plaintext using Triple DES in ECB mode
func tripleDESEncrypt(plaintext []byte, key []byte) (ciphertext []byte, err error) {
	// Create a new DES cipher
	var block cipher.Block
	if block, err = des.NewTripleDESCipher(key); err != nil {
		return nil, err
	}

	// Pad the plaintext to be a multiple of the block size
	paddedPlaintext, err := pkcs7pad(plaintext, block.BlockSize())
	if err != nil {
		return nil, err
	}

	// Initialize byte array for the ciphertext
	ciphertext = make([]byte, len(paddedPlaintext))

	// encrypt in ECB mode
	for i := 0; i < len(paddedPlaintext); i += block.BlockSize() {
		block.Encrypt(ciphertext[i:i+block.BlockSize()], paddedPlaintext[i:i+block.BlockSize()])
	}

	return ciphertext, nil
}
