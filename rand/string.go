package rand

import (
	"crypto/rand"
	"encoding/base64"
	"lens-locked-go/model"
)

const byteSliceSize = 32

func GenerateTokenString() (string, *model.ApiError) {
	return generateString(byteSliceSize)
}

func generateString(size uint) (string, *model.ApiError) {
	b, err := generateBytes(size)

	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func generateBytes(size uint) ([]byte, *model.ApiError) {
	b := make([]byte, size)

	_, err := rand.Read(b)

	if err != nil {
		return nil, model.NewInternalServerApiError(err.Error())
	}
	return b, nil
}
