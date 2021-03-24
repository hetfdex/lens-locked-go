package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
	"lens-locked-go/model"
)

type hasher struct {
	hash.Hash
}

func New(key string) *hasher {
	h := hmac.New(sha256.New, []byte(key))

	return &hasher{
		h,
	}
}

func (h *hasher) GenerateHash(input string) (string, *model.ApiError) {
	h.Reset()

	_, err := h.Write([]byte(input))

	if err != nil {
		return "", model.NewInternalServerApiError(err.Error())
	}
	b := h.Sum(nil)

	return base64.URLEncoding.EncodeToString(b), nil
}
