package mail_sender

type Operation int

const (
	TIMEOUT Operation = iota
	CONN
	EHLO
	STARTTLS
	RCPT_TO
	MAIL_FROM
	DATA
	QUIT
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

	if op < TIMEOUT || op > QUIT {
		return "Unknown"
	}

	return names[op]
}
