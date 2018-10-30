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
		"QUIT"}

	if op < Timeout || op > Quit {
		return "Unknown"
	}

	return names[op]
}
