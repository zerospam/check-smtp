package lib

import (
	"fmt"
	"github.com/zerospam/check-smtp/lib/mail-sender/smtp-commands"
)

type SmtpError struct {
	Command      smtp_commands.Commands `json:"command"`
	ErrorMessage string                 `json:"error_msg"`
}

type CheckResult struct {
	Request *TransportServer `json:"request"`
	Success bool             `json:"success"`
	Error   *SmtpError       `json:"error_message,omitempty"`
}

func NewSmtpError(Op smtp_commands.Commands, err error) *SmtpError {
	return &SmtpError{Command: Op, ErrorMessage: err.Error()}
}

func (se *SmtpError) String() string {
	return fmt.Sprintf("%s: %s", se.Command, se.ErrorMessage)
}
