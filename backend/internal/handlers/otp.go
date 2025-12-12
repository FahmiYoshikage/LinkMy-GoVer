package handlers

import (
	"crypto/rand"
	"fmt"
	"net/smtp"
	"os"
	"sync"
	"time"
)

// OTP storage - in production use Redis
type OTPEntry struct {
	Code      string
	Email     string
	ExpiresAt time.Time
}

var otpStore = struct {
	sync.RWMutex
	entries map[string]OTPEntry
}{entries: make(map[string]OTPEntry)}

// GenerateOTP creates a 6-digit OTP
func GenerateOTP() string {
	b := make([]byte, 3)
	rand.Read(b)
	return fmt.Sprintf("%06d", int(b[0])*int(b[1])*int(b[2])%1000000)
}

// StoreOTP saves OTP for email
func StoreOTP(email, otp string) {
	otpStore.Lock()
	defer otpStore.Unlock()
	otpStore.entries[email] = OTPEntry{
		Code:      otp,
		Email:     email,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
}

// VerifyOTP checks if OTP is valid
func VerifyOTP(email, otp string) bool {
	otpStore.RLock()
	defer otpStore.RUnlock()
	
	entry, exists := otpStore.entries[email]
	if !exists {
		return false
	}
	
	if time.Now().After(entry.ExpiresAt) {
		return false
	}
	
	return entry.Code == otp
}

// DeleteOTP removes OTP after use
func DeleteOTP(email string) {
	otpStore.Lock()
	defer otpStore.Unlock()
	delete(otpStore.entries, email)
}

// SendOTPEmail sends OTP via Gmail SMTP
func SendOTPEmail(toEmail, otp string) error {
	smtpHost := os.Getenv("MAIL_HOST")
	smtpPort := os.Getenv("MAIL_PORT")
	smtpUser := os.Getenv("MAIL_USERNAME")
	smtpPass := os.Getenv("MAIL_PASSWORD")
	mailFrom := os.Getenv("MAIL_FROM")
	
	if smtpHost == "" || smtpUser == "" || smtpPass == "" {
		return fmt.Errorf("SMTP not configured: host=%s, user=%s", smtpHost, smtpUser)
	}
	
	if mailFrom == "" {
		mailFrom = smtpUser
	}
	
	subject := "Kode Verifikasi LinkMy"
	body := fmt.Sprintf(`
Hai!

Kode verifikasi akun LinkMy kamu adalah:

%s

Kode ini berlaku selama 10 menit.

Jika kamu tidak meminta kode ini, abaikan email ini.

---
Tim LinkMy
`, otp)
	
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		mailFrom, toEmail, subject, body)
	
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	return smtp.SendMail(addr, auth, mailFrom, []string{toEmail}, []byte(msg))
}
