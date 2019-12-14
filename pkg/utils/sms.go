package utils

import "github.com/sfreiberg/gotwilio"

var accoutSID string
var authToken string
var phone string
var twilioClient *gotwilio.Twilio

//InitTwilio init the twilio client with the pn sid and token
func InitTwilio(sid, auth, num string) {
	accoutSID = sid
	authToken = auth
	phone = num
	twilioClient = gotwilio.NewTwilioClient(accoutSID, authToken)
}

//SendSMS send sms via twilio client
func SendSMS(to string, msg string) {
	twilioClient.SendSMS(phone, to, msg, "", "")
}
