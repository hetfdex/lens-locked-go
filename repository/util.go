package repository

import "fmt"

const noFieldToQueryErrorMessage = "no field to query"

func notFoundErrorMessage(value string) string {
	message := fmt.Sprintf("%s not found", value)

	return message
}
