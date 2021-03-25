package rand

import (
	"crypto/rand"
	"encoding/base64"
	"lens-locked-go/model"
	"lens-locked-go/util"
)

func GenerateTokenString() (string, *model.ApiError) {
	return generateString(util.ByteSliceSize)
}

func generateString(size uint) (string, *model.ApiError) {
	b, err := generateBytes(size)

	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func generateBytes(size uint) ([]byte, *model.ApiError) {
	if size < 16 {
		return nil, model.NewInternalServerApiError("byte slice size must be at least 16")
	}
	b := make([]byte, size)

	_, err := rand.Read(b)

	if err != nil {
		return nil, model.NewInternalServerApiError(err.Error())
	}
	return b, nil
}
