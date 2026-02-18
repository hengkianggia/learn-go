package email

import "fmt"

// GetOTPTemplate returns a beautifully designed HTML email template for OTP verification.
func GetOTPTemplate(otp string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Your OTP Code</title>
    <style>
        body {
            font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif;
            background-color: #f4f6f8;
            margin: 0;
            padding: 0;
            -webkit-font-smoothing: antialiased;
        }
        .container {
            max-width: 600px;
            margin: 40px auto;
            background-color: #ffffff;
            border-radius: 12px;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
            overflow: hidden;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            padding: 30px;
            text-align: center;
            color: #ffffff;
        }
        .header h1 {
            margin: 0;
            font-size: 28px;
            font-weight: 700;
            letter-spacing: 1px;
        }
        .content {
            padding: 40px 30px;
            text-align: center;
            color: #333333;
        }
        .greeting {
            font-size: 18px;
            margin-bottom: 20px;
            color: #555555;
        }
        .otp-box {
            background-color: #f0f4f8;
            border: 2px dashed #667eea;
            border-radius: 8px;
            padding: 20px;
            margin: 30px 0;
            display: inline-block;
        }
        .otp-code {
            font-size: 36px;
            font-weight: 800;
            color: #2d3748;
            letter-spacing: 8px;
            margin: 0;
        }
        .instruction {
            font-size: 16px;
            line-height: 1.6;
            color: #4a5568;
            margin-bottom: 30px;
        }
        .expiry {
            font-size: 14px;
            color: #e53e3e;
            font-weight: 600;
        }
        .footer {
            background-color: #f9fafb;
            padding: 20px;
            text-align: center;
            font-size: 12px;
            color: #a0aec0;
            border-top: 1px solid #edf2f7;
        }
        .footer p {
            margin: 5px 0;
        }
        @media only screen and (max-width: 600px) {
            .container {
                margin: 20px;
                width: auto;
            }
            .content {
                padding: 20px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>LearnGo App</h1>
        </div>
        <div class="content">
            <p class="greeting">Hello!</p>
            <p class="instruction">We received a request to verify your email address. Please use the following One-Time Password (OTP) to complete your registration:</p>
            
            <div class="otp-box">
                <h2 class="otp-code">%s</h2>
            </div>
            
            <p class="expiry">This code will expire in 5 minutes.</p>
            <p class="instruction" style="font-size: 14px; color: #718096; margin-top: 30px;">
                If you did not request this code, please ignore this email.
            </p>
        </div>
        <div class="footer">
            <p>&copy; 2026 LearnGo App. All rights reserved.</p>
            <p>This is an automated message, please do not reply.</p>
        </div>
    </div>
</body>
</html>
	`, otp)
}
