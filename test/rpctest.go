package main

import (
	"eth-withdraw/email"
	"eth-withdraw/config"
)

func init() {
	config.Load()
}

func main() {
	email.SendNotification("zzz@mail.ru", 12)
}