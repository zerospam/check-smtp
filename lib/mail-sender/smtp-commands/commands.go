package smtp_commands

type Commands string

const (
	Timeout    Commands = "TIMEOUT"
	Connection          = "CONNECTION"
	Ehlo                = "EHLO"
	StartTls            = "STARTTLS"
	RcptTo              = "RCPT TO"
	MailFrom            = "MAIL FROM"
	Data                = "DATA"
	Quit                = "QUIT"
	SpfFail             = "SPF-FAIL"
)

type CommandLog map[Commands]string
