package validator

import (
	"fmt"
	"lens-locked-go/model"
)

func StringNotEmpty(name string, input string) *model.ApiError {
	if input == "" {
		message := fmt.Sprintf("%s must not be empty", name)

		return model.NewInternalServerApiError(message)
	}
	return nil
}
