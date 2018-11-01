package test_email

import "github.com/zerospam/check-smtp/lib"

type TestEmailRequest struct {
	From    string               `json:"from"`
	Body    string               `json:"body"`
	Subject string               `json:"subject"`
	Server  *lib.TransportServer `json:"server"`
}

func (t *TestEmailRequest) ToTestEmail() *Email {
	return NewTestEmail(t.Subject, t.Body, t.From, t.Server.TestEmail)
}
