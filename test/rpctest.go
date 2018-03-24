package main

import (
	"withdraw/email"
	"withdraw/config"
)

func init() {
	config.Load()
}

func main() {
	email.SendNotification("zzz@mail.ru", 12)
}