package util

import (
	"fmt"
	"regexp"
)

const host = "localhost"
const port = 5432
const user = "postgres"
const password = "Abcde12345!"
const dbname = "lenslocked_dev"

var Dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

const Address = "localhost:8080"

const Pepper = "6Sk65RHhGW7S4qnVPV7m"

const ByteSliceSize = 32

const HasherKey = "yzzmGPkAA9FTmbtzz9jB"

const CookieName = "login_token"

const BaseTag = "base"
const BaseFilename = "view/base.gohtml"
const AlertFilename = "view/alert.gohtml"

const ContentTypeKey = "Content-Type"
const ContentTypeValue = "text/html"

var EmailRegex = regexp.MustCompile(`^[a-z0-9_.+-]+@[a-z0-9-]+\.[a-z0-9-.]+$`)

const InvalidPasswordErrorMessage = "invalid password"
const InvalidEmailErrorMessage = "invalid email address"
const EmailInUseErrorMessage = "email address is already in use"
const NoUserUpdateNeededErrorMessage = "no user update needed"
const ByteSliceSizeErrorMessage = "byte slice size must be at least 16"

func MustNotBeEmptyErrorMessage(value string) string {
	message := fmt.Sprintf("%s must not be empty", value)

	return message
}
