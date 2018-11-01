package test_email

import (
	"crypto/sha1"
	"fmt"
	"github.com/rs/xid"
	"io"
	"strings"
	"time"
)

type Email struct {
	From    string
	To      string
	Body    string
	Subject string
	Headers map[string]string
}

func (e *Email) generateMessageId(localName string) string {
	hasherSha1 := sha1.New()

	io.WriteString(hasherSha1, e.Subject)
	io.WriteString(hasherSha1, e.From)
	io.WriteString(hasherSha1, e.To)
	guid := xid.New()

	return fmt.Sprintf("%s-%x@%s", guid, sha1.Sum(nil), localName)
}

func (e *Email) PrepareHeaders(localName string) {

	e.Headers = map[string]string{
		"From":       fmt.Sprintf("Mail Server Tester <%s>", e.From),
		"To":         fmt.Sprintf("Mail Server Tester Receiver <%s>", e.To),
		"Subject":    e.Subject,
		"Date":       time.Now().Format(time.RFC822Z),
		"Message-Id": fmt.Sprintf("<%s>", e.generateMessageId(localName)),

		"MIME-Version":              "1.0",
		"Content-Transfer-Encoding": "8bit",

		"Auto-Submitted": "auto-generated",
		"X-Mailer":       "SMTP Server Tester",
		"Content-Type":   "text/plain; charset=\"UTF-8\"",
	}
}

func NewTestEmail(subject string, body string, from string, to string) *Email {

	return &Email{
		Body:    body,
		From:    from,
		To:      to,
		Subject: subject,
	}
}

func (e *Email) String() string {
	var builder strings.Builder
	builder.Grow(len(e.Body) + len(e.Headers)*10)
	for header, value := range e.Headers {
		builder.WriteString(fmt.Sprintf("%s: %s\n", header, value))
	}
	builder.WriteString("\n")
	builder.WriteString(e.Body)
	return builder.String()
}
