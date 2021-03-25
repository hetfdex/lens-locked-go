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

func New(hasherKey string) (*Hasher, *model.ApiError) {
	if validator.EmptyString(hasherKey) {
		return nil, model.NewInternalServerApiError("hasherKey must not be empty")
	}
	h := hmac.New(sha256.New, []byte(hasherKey))

	return &Hasher{
		hash: h,
	}, nil
}

func (h *Hasher) GenerateTokenHash(token string) (string, *model.ApiError) {
	if validator.EmptyString(token) {
		return "", model.NewInternalServerApiError("token must not be empty")
	}
	h.hash.Reset()

	_, err := h.hash.Write([]byte(token))

	if err != nil {
		return "", model.NewInternalServerApiError(err.Error())
	}
	b := h.hash.Sum(nil)

	return base64.URLEncoding.EncodeToString(b), nil
}
