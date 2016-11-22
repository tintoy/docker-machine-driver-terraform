package main

import (
	"crypto/rand"
	"encoding/base64"
)

// The number of bytes in a single-use password.
const oneTimePasswordLength = 16

// Generate a cryptographically-random single use password.
func (driver *Driver) generateOneTimePassword() (password string, err error) {
	data := make([]byte, oneTimePasswordLength)
	_, err = rand.Read(data)
	if err != nil {
		return
	}

	encodedData := make([]byte,
		base64.StdEncoding.EncodedLen(oneTimePasswordLength),
	)
	base64.StdEncoding.Encode(encodedData, data)

	password = string(encodedData)

	return
}
