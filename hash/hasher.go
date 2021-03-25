package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
	"lens-locked-go/model"
	"lens-locked-go/validator"
)

type Hasher struct {
	hash hash.Hash
}

func New(key string) *Hasher {
	h := hmac.New(sha256.New, []byte(key))

	return &Hasher{
		hash: h,
	}
}

func (h *Hasher) GenerateHash(input string) (string, *model.ApiError) {
	apiErr := validator.StringNotEmpty("input", input)

	if apiErr != nil {
		return "", apiErr
	}
	h.hash.Reset()

	_, err := h.hash.Write([]byte(input))

	if err != nil {
		return "", model.NewInternalServerApiError(err.Error())
	}
	b := h.hash.Sum(nil)

	return base64.URLEncoding.EncodeToString(b), nil
}
