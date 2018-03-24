package email

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
	"github.com/spf13/viper"
)

func SendNotification(emailAddress string, amount int) {

	m := gomail.NewMessage()
	m.SetHeader("From", "no-reply@blockbase.tech")
	m.SetHeader("To", emailAddress)
	m.SetAddressHeader("Cc", "no-reply@blockbase.tech", "PlanEX")
	m.SetHeader("Subject", "Deposit Success Alert")
	m.SetBody("text/html", "Hello, Your PlanEX's account has recharged "+string(amount)+" ETH.")

	d := gomail.NewDialer(viper.GetString("EmailServer"), viper.GetInt("EmailPort"),
		viper.GetString("EmailUser"), viper.GetString("EmailPassword"))
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
