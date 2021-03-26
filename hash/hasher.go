package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
	"lens-locked-go/model"
)

type Hasher struct {
	hash hash.Hash
}

func New(hasherKey string) (*Hasher, *model.ApiError) {
	if hasherKey == "" {
		return nil, model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("hasherKey"))
	}
	h := hmac.New(sha256.New, []byte(hasherKey))

	return &Hasher{
		hash: h,
	}, nil
}

func (h *Hasher) GenerateTokenHash(token string) (string, *model.ApiError) {
	if token == "" {
		return "", model.NewInternalServerApiError(model.MustNotBeEmptyErrorMessage("token"))
	}
	h.hash.Reset()

	_, err := h.hash.Write([]byte(token))

	if err != nil {
		return "", model.NewInternalServerApiError(err.Error())
	}
	b := h.hash.Sum(nil)

	return base64.URLEncoding.EncodeToString(b), nil
}
