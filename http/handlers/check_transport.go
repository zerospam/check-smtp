package handlers

import (
	"encoding/json"
	"github.com/zerospam/check-smtp/lib"
	"github.com/zerospam/check-smtp/lib/environment-vars"
	"github.com/zerospam/check-smtp/lib/mail-sender/smtp-commands"
	"github.com/zerospam/check-smtp/test-email"
	"log"
	"net/http"
	"strings"
)

func getRequestIp(req *http.Request) string {
	if header := req.Header.Get("X-Forwarded-For"); header != "" {
		exploded := strings.Split(header, ",")
		return strings.Trim(exploded[len(exploded)-1], " ")
	}

	return req.RemoteAddr
}

func generateResult(smtpError *lib.SmtpError, banner string, tlsVersion string, generalLog smtp_commands.CommandLog, spfLog smtp_commands.CommandLog) lib.CheckResult {
	success := smtpError == nil
	return lib.CheckResult{
		Error:       smtpError,
		Success:     success,
		TlsVersion:  tlsVersion,
		HelloBanner: banner,
		GeneralLog:  generalLog,
		SPFLog:      spfLog,
	}
}

func testServer(server *lib.TransportServer, email *test_email.Email) lib.CheckResult {
	client, err := environmentvars.GetVars().NewSmtpClient(server)
	if err != nil {
		return generateResult(err, "", "", make(smtp_commands.CommandLog), make(smtp_commands.CommandLog))
	}

	err = client.SendTestEmail(email)

	if err != nil {
		banner, tlsVersion := client.GetHelloBanner()
		return generateResult(err, banner, tlsVersion, client.GetCommandLog(), make(smtp_commands.CommandLog))
	}

	generalLog := client.GetCommandLog()
	//new client to do the spoofing
	//Can't reuse previous client as it closed the connection
	client, err = environmentvars.GetVars().NewSmtpClient(server)
	if err != nil {
		banner, tlsVersion := client.GetHelloBanner()
		return generateResult(err, banner, tlsVersion, generalLog, client.GetCommandLog())
	}

	err = client.SpoofingTest(environmentvars.GetVars().SmtpMailSpoof.Address)

	if err != nil {
		if err.Command == smtp_commands.RcptTo || err.Command == smtp_commands.MailFrom {
			err.Command = smtp_commands.SpfFail
		}
		banner, tlsVersion := client.GetHelloBanner()
		return generateResult(err, banner, tlsVersion, generalLog, client.GetCommandLog())
	}

	banner, tlsVersion := client.GetHelloBanner()
	return generateResult(err, banner, tlsVersion, generalLog, client.GetCommandLog())

}

func CheckTransport(w http.ResponseWriter, req *http.Request) {
	var testEmailRequest test_email.TestEmailRequest

	if req.Method != "POST" {
		http.Error(w, "Only POST accepted.", 405)
		return
	}

	if req.Header.Get("Authorization") != environmentvars.GetVars().SharedKey {
		http.Error(w, "Wrong Key sent.", 402)
		log.Printf("[%s] - %s - %v\n", req.RemoteAddr, req.Method, "REJECT")
		return
	}

	if req.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	defer req.Body.Close()
	json.NewDecoder(req.Body).Decode(&testEmailRequest)

	w.Header().Add("Content-Type", "application/json")
	email := testEmailRequest.ToTestEmail()
	server := testEmailRequest.Server

	result := testServer(server, email)

	json.NewEncoder(w).Encode(result)

	log.Printf("[%s] - %s (%s:%d) - %v\n", getRequestIp(req), req.Method, testEmailRequest.Server.Server, testEmailRequest.Server.Port, result.Success)

}
