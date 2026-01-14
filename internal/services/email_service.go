package services

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/khusa-mahal/backend/internal/models"
)

type EmailService struct {
	smtpHost string
	smtpPort string
	username string
	password string
}

func NewEmailService() *EmailService {
	return &EmailService{
		smtpHost: os.Getenv("SMTP_HOST"),
		smtpPort: os.Getenv("SMTP_PORT"),
		username: os.Getenv("SMTP_USERNAME"),
		password: os.Getenv("SMTP_PASSWORD"),
	}
}

func (s *EmailService) SendOTP(to, code string) error {
	from := s.username
	subject := "Your Verification Code - Khusa Mahal"
	body := fmt.Sprintf("Your verification code is: %s\n\nThis code will expire in 10 minutes.", code)

	message := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", to, subject, body))

	auth := smtp.PlainAuth("", s.username, s.password, s.smtpHost)

	// In development/testing if no creds, just log it
	if s.username == "" || s.password == "" {
		fmt.Printf(" [MOCK EMAIL] To: %s | OTP: %s\n", to, code)
		return nil
	}

	err := smtp.SendMail(s.smtpHost+":"+s.smtpPort, auth, from, []string{to}, message)
	if err != nil {
		// Log detailed error for debugging
		fmt.Printf("Failed to send email: %v\n", err)
		return err
	}

	return nil
}

func (s *EmailService) SendOrderConfirmationEmail(to string, orderID string, items []models.OrderDetailsItem, shippingAddress models.Address, total float64) error {
	from := s.username
	subject := "Order Confirmation - Khusa Mahal"

	fmt.Printf("Attempting to send email to: %s for Order: %s\n", to, orderID)

	// Build product rows with images
	rowStr := ""
	for _, item := range items {
		rowStr += fmt.Sprintf(`
			<tr style="border-bottom: 1px solid #eee;">
				<td style="padding: 15px;">
					<div style="display: flex; align-items: center; gap: 15px;">
						<img src="%s" alt="%s" style="width: 80px; height: 80px; object-fit: cover; border-radius: 8px; border: 1px solid #ddd;">
						<div>
							<div style="font-weight: bold; margin-bottom: 5px;">%s</div>
							<div style="font-size: 12px; color: #666;">Size: %s | Color: %s</div>
						</div>
					</div>
				</td>
				<td style="padding: 15px; text-align: center;">%d</td>
				<td style="padding: 15px; text-align: right; font-weight: bold;">PKR %.2f</td>
			</tr>
		`, item.Image, item.Name, item.Name, item.Size, item.Color, item.Quantity, item.Price*float64(item.Quantity))
	}

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
<style>
	body { font-family: Arial, sans-serif; color: #333; margin: 0; padding: 0; background-color: #f4f4f4; }
	.container { width: 100%%; max-width: 600px; margin: 20px auto; background: white; border-radius: 10px; overflow: hidden; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
	.header { background: linear-gradient(135deg, #800000 0%%, #a00000 100%%); color: white; padding: 30px; text-align: center; }
	.header h1 { margin: 0; font-size: 28px; }
	.content { padding: 30px; }
	.order-id { background: #f9f9f9; padding: 15px; border-radius: 8px; margin-bottom: 20px; }
	table { width: 100%%; border-collapse: collapse; margin: 20px 0; }
	th { background: #f9f9f9; text-align: left; padding: 12px; font-size: 14px; color: #666; text-transform: uppercase; }
	.address-box { background: #f9f9f9; padding: 20px; border-radius: 8px; margin: 20px 0; }
	.address-box h3 { margin-top: 0; color: #800000; }
	.total-row { border-top: 2px solid #800000; margin-top: 20px; padding-top: 15px; text-align: right; }
	.total-row .amount { font-size: 24px; font-weight: bold; color: #800000; }
	.footer { background: #f9f9f9; padding: 20px; text-align: center; font-size: 12px; color: #777; }
	.footer a { color: #800000; text-decoration: none; }
</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>✓ Order Confirmed!</h1>
			<p style="margin: 10px 0 0 0; opacity: 0.9;">Thank you for shopping with Khusa Mahal</p>
		</div>
		
		<div class="content">
			<div class="order-id">
				<strong>Order ID:</strong> %s
			</div>
			
			<h3 style="color: #800000; margin-top: 30px;">Order Details</h3>
			<table>
				<thead>
					<tr>
						<th>Product</th>
						<th style="text-align: center;">Qty</th>
						<th style="text-align: right;">Total</th>
					</tr>
				</thead>
				<tbody>
					%s
				</tbody>
			</table>

			<div class="total-row">
				<div style="margin-bottom: 10px; font-size: 16px;">Total Amount</div>
				<div class="amount">PKR %.2f</div>
			</div>

			<div class="address-box">
				<h3>Shipping Address</h3>
				<p style="margin: 5px 0; line-height: 1.6;">
					%s<br>
					%s, %s %s<br>
					%s
				</p>
			</div>

			<p style="margin-top: 30px; color: #666; font-size: 14px;">
				Your order is being processed and will be shipped soon. You will receive a tracking number once your order ships.
			</p>
		</div>

		<div class="footer">
			<p><strong>Khusa Mahal</strong> - Traditional Elegance</p>
			<p>Questions? Contact us at <a href="mailto:support@khusamahal.com">support@khusamahal.com</a></p>
		</div>
	</div>
</body>
</html>
`, orderID, rowStr, total, shippingAddress.Street, shippingAddress.City, shippingAddress.State, shippingAddress.ZipCode, shippingAddress.Country)

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n%s\r\n%s", to, subject, mime, htmlBody))

	auth := smtp.PlainAuth("", s.username, s.password, s.smtpHost)

	if s.username == "" || s.password == "" {
		fmt.Printf(" [MOCK EMAIL] To: %s | Order Confirmation for %s\n", to, orderID)
		return nil
	}

	err := smtp.SendMail(s.smtpHost+":"+s.smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		fmt.Printf("Failed to send order email: %v\n", err)
		return err
	}
	return nil
}
