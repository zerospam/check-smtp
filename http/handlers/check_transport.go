package handlers

import (
	"encoding/json"
	"github.com/zerospam/check-smtp/lib"
	"github.com/zerospam/check-smtp/lib/environment-vars"
	"github.com/zerospam/check-smtp/lib/mail-sender"
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

func generateResult(server *lib.TransportServer, smtpError *lib.SmtpError) lib.CheckResult {
	success := smtpError == nil
	return lib.CheckResult{
		Request: server,
		Error:   smtpError,
		Success: success,
	}
}

func testServer(server *lib.TransportServer, email *test_email.Email) lib.CheckResult {
	client, err := mail_sender.NewClient(server, environmentvars.GetVars().SmtpCN, environmentvars.GetVars().SmtpConnectionTimeout, environmentvars.GetVars().SmtpOperationTimeout)
	if err != nil {
		return generateResult(server, err)
	}

	err = client.SendTestEmail(email)

	if err != nil {
		return generateResult(server, err)
	}

	//new client to do the spoofing
	//Can't reuse previous client as it closed the connection
	client, err = mail_sender.NewClient(server, environmentvars.GetVars().SmtpCN, environmentvars.GetVars().SmtpConnectionTimeout, environmentvars.GetVars().SmtpOperationTimeout)
	if err != nil {
		return generateResult(server, err)
	}

	err = client.SpoofingTest(environmentvars.GetVars().SmtpMailSpoof.Address)

	if err != nil {
		return generateResult(server, err)
	}

	return generateResult(server, nil)

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

	json.NewDecoder(req.Body).Decode(&testEmailRequest)

	defer req.Body.Close()

	w.Header().Add("Content-Type", "application/json")
	email := testEmailRequest.ToTestEmail()
	server := testEmailRequest.Server

	result := testServer(server, email)

	json.NewEncoder(w).Encode(result)

	log.Printf("[%s] - %s (%s:%d) - %v\n", getRequestIp(req), req.Method, testEmailRequest.Server.Server, testEmailRequest.Server.Port, result.Success)

}
