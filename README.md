# Check SMTP

Mini HTTP service that takes a JSON with server information and check
 if it's accessible from the application and if it can receive emails.

## Env

| Key        | Requirement           | Explanation                                                                                   |
|------------|-----------------------|-----------------------------------------------------------------------------------------------|
| SHARED_KEY | mandatory              | Secret shared between main app and this one. (Needs to be sent in the header *Authorization*) |
| PORT       | optional (default 80)  | Port used for the application                                                                 |
| SMTP_CN | optional (default hostname)| Common Name to use for client certificate when doing a STARTTLS                                                                |
| SMTP_FROM | optional (default local@hostname)| Email to use to do the MAIL FROM smtp command                                                               |
| SMTP_FROM_SPOOF | optional (default spoof@amazon.com)| Email to check for SPF checks (spoofing this email as MAIL FROM)                                              |
| SMTP_TIMEOUT | optional (30 seconds)| How long to wait for the SMTP server to answer|

## Data
```json
{
  "server": "example.com",
  "port": 25,
  "mx": false,
  "test_email": "test@example.com"
}
```


| Key    | Explanation                                                                    |
|--------|--------------------------------------------------------------------------------|
| server | Server to check                                                                |
| port   | Port to use to attempt the connection                                          |
| mx     | Instead of resolving the IP, resolve the MX of the server first then check IPs |
| test_email     | Used as RCPT TO when doing SMTP checks |
