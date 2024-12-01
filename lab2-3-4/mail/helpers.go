package mail

import (
	"fmt"
	"strings"
)

func GenerateOtpEmail(otp string) string {
	htmlTemplate := `
<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f9f9f9;
            margin: 0;
            padding: 20px;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #ffffff;
            border-radius: 8px;
            padding: 20px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
        }
        .header {
            text-align: center;
            background-color: #007bff;
            color: #ffffff;
            padding: 10px 0;
            border-radius: 8px 8px 0 0;
        }
        .content {
            margin: 20px 0;
            text-align: center;
        }
        .otp {
            font-size: 24px;
            font-weight: bold;
            color: #007bff;
        }
        .footer {
            text-align: center;
            font-size: 12px;
            color: #777777;
            margin-top: 20px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Verification Code</h1>
        </div>
        <div class="content">
            <p>Here is youur verification code:</p>
            <p class="otp">%s</p>

        </div>
        <div class="footer">
            <p>If you didn’t request this code, please ignore this email.</p>
        </div>
    </div>
</body>
</html>
`
	return fmt.Sprintf(htmlTemplate, strings.TrimSpace(otp))
}

func GenerateWelcomeEmail(email string) string {
	htmlTemplate := `
<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f9f9f9;
            margin: 0;
            padding: 20px;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #ffffff;
            border-radius: 8px;
            padding: 20px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
        }
        .header {
            text-align: center;
            background-color: #28a745;
            color: #ffffff;
            padding: 10px 0;
            border-radius: 8px 8px 0 0;
        }
        .content {
            margin: 20px 0;
            text-align: center;
        }
        .welcome {
            font-size: 20px;
            font-weight: bold;
            color: #333333;
        }
        .footer {
            text-align: center;
            font-size: 12px;
            color: #777777;
            margin-top: 20px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome!</h1>
        </div>
        <div class="content">
            <p class="welcome">Hi %s,</p>
        </div>
        <div class="footer">
            <p>Thank you for joining us!</p>
        </div>
    </div>
</body>
</html>
`
	return fmt.Sprintf(htmlTemplate, strings.TrimSpace(email))
}
