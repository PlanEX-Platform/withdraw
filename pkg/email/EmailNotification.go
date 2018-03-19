package email

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
	"eth-withdraw/pkg/config"
)

func sendNotification(emailAddress string, amount int)  {

	m := gomail.NewMessage()
	m.SetHeader("From", "no-reply@planex.io")
	m.SetHeader("To", emailAddress)
	m.SetAddressHeader("Cc", "no-reply@planex.io", "PlanEX")
	m.SetHeader("Subject", "Deposit Success Alert")
	m.SetBody("text/html", "Hello, Your PlanEX's account has recharged " + string(amount) +" ETH.")


	d := gomail.NewDialer(config.CFG.EmailServer, config.CFG.EmailPort, config.CFG.EmailUser, config.CFG.EmailPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
