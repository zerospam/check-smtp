package smtp_commands

type Commands int

const (
	Timeout Commands = iota
	Connection
	Ehlo
	StartTls
	RcptTo
	MailFrom
	Data
	Quit
	SpfFail
)

func (op Commands) String() string {
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
