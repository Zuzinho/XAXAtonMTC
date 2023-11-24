package sms

import (
	"github.com/dex35/smsru"
)

type Sender struct {
	smsApiID string
}

func NewSender(smsApiID string) *Sender {
	return &Sender{
		smsApiID: smsApiID,
	}
}

func (sender *Sender) SendSms(phoneNumber, text string) error {
	client := smsru.CreateClient(sender.smsApiID)

	sms := smsru.CreateSMS(phoneNumber, text)

	_, err := client.SmsSend(sms)

	return err
}
