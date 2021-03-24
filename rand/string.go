package rand

import (
	"crypto/rand"
	"encoding/base64"
	"lens-locked-go/config"
	"lens-locked-go/model"
)

func GenerateString(size uint) (string, *model.ApiError) {
	if size == 0 {
		size = config.ByteSliceSize
	}
	return generateString(size)
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
