package repository

import "fmt"

func notFoundErrorMessage(value string) string {
	message := fmt.Sprintf("%s not found", value)

	return message
}
