package test_email

import (
	"crypto/sha1"
	"fmt"
	"github.com/rs/xid"
	"github.com/zerospam/check-smtp/lib/environment-vars"
	"io"
	"strings"
	"time"
)

type TestEmail struct {
	Body    string
	Headers map[string]string
}

func generateMessageId(subject string, from string, to string) string {
	hasherSha1 := sha1.New()

	io.WriteString(hasherSha1, subject)
	io.WriteString(hasherSha1, from)
	io.WriteString(hasherSha1, to)
	guid := xid.New()

	return fmt.Sprintf("%s-%x@%s", guid, sha1.Sum(nil), environmentvars.GetVars().SmtpCN)
}

func NewTestEmail(subject string, body string, from string, to string) *TestEmail {

	return &TestEmail{
		Body: body,
		Headers: map[string]string{
			"From":       fmt.Sprintf("Mail Tester <%s>", from),
			"To":         fmt.Sprintf("Mail Tester Receiver <%s>", to),
			"Subject":    subject,
			"Date":       time.Now().Format(time.RFC822Z),
			"Message-Id": fmt.Sprintf("<%s>", generateMessageId(subject, from, to)),

			"MIME-Version":              "1.0",
			"Content-Transfer-Encoding": "8bit",

			"Auto-Submitted": "auto-generated",
			"X-Mailer":       "SMTP Server Tester",
			"Content-Type":   "text/plain; charset=\"UTF-8\"",
		},
	}
}

func (t *TestEmail) String() string {
	var builder strings.Builder
	builder.Grow(len(t.Body))
	for header, value := range t.Headers {
		builder.WriteString(fmt.Sprintf("%s: %s\n", header, value))
	}
	builder.WriteString("\n")
	builder.WriteString(t.Body)
	return builder.String()
}
