package smtp

type Operation int

const (
	EHLO Operation = iota
	STARTTLS
	RCPT_TO
	MAIL_FROM
	DATA
	QUIT
)

func (op Operation) String() string {
	names := [...]string{
		"EHLO",
		"STARTTLS",
		"RCPT_TO",
		"MAIL_FROM",
		"DATA",
		"QUIT"}

	if op < EHLO || op > QUIT {
		return "Unknown"
	}

	return names[op]
}
