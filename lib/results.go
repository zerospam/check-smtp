package lib

import "github.com/zerospam/check-smtp/lib/mail-sender"

type SmtpError struct {
	Operation    mail_sender.Operation
	ErrorMessage string
}

type CheckResult struct {
	Request *TransportServer `json:"request"`
	Success bool             `json:"success"`
	Error   *SmtpError       `json:"error_message,omitempty"`
}

func NewSmtpError(Op mail_sender.Operation, err error) *SmtpError {
	return &SmtpError{Operation: Op, ErrorMessage: err.Error()}
}
