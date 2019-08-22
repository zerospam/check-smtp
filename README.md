# Check SMTP

Mini HTTP service that takes a JSON with server information and check
 if it's accessible from the application and if it can receive emails.

## Env

| Key        | Requirement           | Explanation                                                                                   |
|------------|-----------------------|-----------------------------------------------------------------------------------------------|
| SHARED_KEY | mandatory              | Secret shared between main app and this one. (Needs to be sent in the header *Authorization*) |
| PORT       | optional (default 80)  | Port used for the application                                                                 |
| SMTP_CN | optional (default hostname)| Common Name to use for client certificate when doing a STARTTLS                                                                |
| SMTP_FROM_SPOOF | optional (default spoof@amazon.com)| Email to check for SPF checks (spoofing this email as MAIL FROM)                                              |
| SMTP_CONN_TIMEOUT | optional (30 seconds)| How long to wait for the SMTP server to answer|
| SMTP_OPT_TIMEOUT | optional (30 seconds)| How long to wait for the SMTP server to answer each command|
|TLS_MIN_VERSION| optional (TLS1.2) | Check TLS Table for acceptable values. This is used when doing a STARTTLS on server that supports it.|

### TLS

| TLS Version | EnvVar |
|-------------|--------|
| SSLv3.0     | SSL30  |
| TLSv1       | TLS10  |
| TLSv1.1     | TLS11  |
| TLS1.2      | TLS12  |
| TLS1.3      | TLS13  |

## Data
```json
{
  "from": "bounce@myserver.com",
  "body": "Hello World\n Are you doing well ?\n\nTester",
  "subject": "Hello World",
  "server": {
    "server": "example.com",
    "port": 25,
    "mx": false,
    "test_email": "test@example.com"
  }
}
```


| Key    | Explanation                                                                    |
|--------|--------------------------------------------------------------------------------|
| from | Address to use in the MAIL FROM|
| body | Body of the mail|
| subject | Subject of the mail|
| server.server | Server to check                                                                |
| server.port   | Port to use to attempt the connection                                          |
| server.mx     | Instead of resolving the IP, resolve the MX of the server first then check IPs |
| server.test_email     | Used as RCPT TO when doing SMTP checks |


## Response

### Success
```json
{
    "success": true,
    "hello_banner": "220 example.com ESMTP Postfix (Debian/GNU)",
    "tls_version": "VersionTLS12",
    "general_log": {
        "1/CONNECTION": "192.168.22.3:25",
        "2/EHLO": "tardis.example.com",
        "3/STARTTLS": "VersionTLS12",
        "4/MAIL FROM": "test@example.com",
        "5/RCPT TO": "me@example.com",
        "6/DATA": "\nHello World\nAre you doing well ?\n\nTester",
        "7/QUIT": ""
    },
    "spf_log": {
        "1/CONNECTION": "192.168.22.3:25",
        "2/EHLO": "tardis.zerospam.ca",
        "3/STARTTLS": "",
        "8/SPF-FAIL": "antoineaf@admincmd.com"
    }
}
```

### Error
```json
{
    "success": false,
    "hello_banner": "220 set.example.com ESMTP Postfix (Debian/GNU)",
    "tls_version": "VersionTLS12",
    "error_message": {
        "command": "STARTTLS",
        "error_msg": "x509: certificate is valid for example.com, not set.example.com"
    },
    "general_log": {
        "1/CONNECTION": "192.168.22.3:25",
        "2/EHLO": "tardis.example.com",
        "3/STARTTLS": ""       
    },
    "spf_log": {
    }
}
```
