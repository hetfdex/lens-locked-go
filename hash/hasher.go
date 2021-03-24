package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
	"lens-locked-go/model"
)

type Hasher struct {
	h hash.Hash
}

func New(key string) *Hasher {
	h := hmac.New(sha256.New, []byte(key))

	return &Hasher{
		h,
	}
}

func (h *Hasher) GenerateHash(input string) (string, *model.ApiError) {
	if input == "" {
		return "", model.NewInternalServerApiError("string must not be empty")
	}
	h.h.Reset()

	_, err := h.h.Write([]byte(input))

	if err != nil {
		return "", model.NewInternalServerApiError(err.Error())
	}
	b := h.h.Sum(nil)

	return base64.URLEncoding.EncodeToString(b), nil
}
