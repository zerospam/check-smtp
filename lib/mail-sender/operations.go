package mail_sender

type Operation int

const (
	Timeout Operation = iota
	Connection
	Ehlo
	StartTls
	RcptTo
	MailFrom
	Data
	Quit
	SpfFail
)

func (op Operation) String() string {
	names := [...]string{
		"TIMEOUT",
		"CONNECTION",
		"EHLO",
		"STARTTLS",
		"RCPT_TO",
		"MAIL_FROM",
		"DATA",
		"QUIT",
		"SPF_FAIL"}

	if op < Timeout || op > SpfFail {
		return "Unknown"
	}

	return names[op]
}
