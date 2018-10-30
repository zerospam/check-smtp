package lib

import "github.com/zerospam/check-smtp/lib/smtp"

type SmtpError struct {
	Operation    smtp.Operation
	ErrorMessage string
}

type CheckResult struct {
	Request *TransportServer `json:"request"`
	Success bool             `json:"success"`
	Error   *SmtpError       `json:"error_message,omitempty"`
}

func NewSmtpError(Op smtp.Operation, err error) *SmtpError {
	return &SmtpError{Operation: Op, ErrorMessage: err.Error()}
}
