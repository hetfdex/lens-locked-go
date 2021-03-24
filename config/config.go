package config

import "fmt"

const host = "localhost"
const port = 5432
const user = "postgres"
const password = "Abcde12345!"
const dbname = "lenslocked_dev"

var Dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

const Pepper = "6Sk65RHhGW7S4qnVPV7m"

const ByteSliceSize = 32

const HasherKey = "yzzmGPkAA9FTmbtzz9jB"

const CookieName = "login_token"
