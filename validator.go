package mailify

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"strings"
	"time"
)

// getHostname gets the fully qualified domain name for HELO command
// GetHostname attempts to retrieve the fully qualified domain name (FQDN) of the current host.
// It first tries to get the hostname using os.Hostname(). If that fails, it returns a fallback
// hostname "verifier.local". If successful, it then attempts to resolve the IP addresses
// associated with the hostname using net.LookupIP(). If that fails, it returns the hostname.
// If successful, it performs a reverse DNS lookup on the first IPv4 address found using
// net.LookupAddr(). If that succeeds and returns at least one name, it returns the first name
// with the trailing dot removed. If all attempts fail, it returns the hostname with ".local" appended.
func (c *Client) GetHostname() (string, error) {
	// Try to get the hostname
	hostname, err := os.Hostname()
	if err != nil {
		return "verifier.local", fmt.Errorf("failed to get hostname: %v", err) // fallback hostname
	}

	// Try to get the FQDN
	addrs, err := net.LookupIP(hostname)
	if err != nil {
		return hostname, fmt.Errorf("failed to lookup IP for hostname %s: %v", hostname, err)
	}

	// Try to get the reverse DNS lookup
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			names, err := net.LookupAddr(ipv4.String())
			if err == nil && len(names) > 0 {
				return strings.TrimSuffix(names[0], "."), nil
			}
		}
	}

	return hostname + ".local", nil
}

// TryConnectingSMTP attempts to establish an SMTP connection and validate an email address.
// It performs the following steps:
// 1. Creates a new validation result indicating the domain has MX records.
// 2. Creates a new dialer with a timeout.
// 3. Formats the address based on IP version (IPv4 or IPv6).
// 4. Handles connection based on the port (SMTPS or plain/STARTTLS).
// 5. Creates an SMTP client.
// 6. Performs HELO/EHLO command.
// 7. Initiates STARTTLS if available and not already using TLS.
// 8. Sends MAIL FROM command.
// 9. Sends RCPT TO command.
// 10. Interprets the response to determine if the email address is valid.
//
// Parameters:
// - smtpDetails: Details of the SMTP server (IP address, port, server name).
// - senderEmail: The email address of the sender.
// - recipientEmail: The email address of the recipient to be validated.
// - localName: The local name to use in the HELO/EHLO command.
// - useTLS: A boolean indicating whether to use TLS.
//
// Returns:
// - A pointer to a ValidationResult struct containing the validation outcome.
// - An error if any step in the process fails.
func (c *Client) TryConnectingSMTP(smtpDetails *SMTPDetails, recipientEmail, localName string, useTLS bool) (*ValidationResult, error) {

	// Create a new validation result. If we are here, we know the domain has MX records.
	result := &ValidationResult{
		IsValid: false,
		HasMX:   true,
	}

	// Create a new dialer with a timeout
	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}

	// Format address based on IP version
	var address string
	ip := net.ParseIP(smtpDetails.IPAddress)
	if ip.To4() != nil {
		// IPv4
		address = fmt.Sprintf("%s:%s", smtpDetails.IPAddress, smtpDetails.Port)
	} else {
		// IPv6 - wrap in square brackets
		address = fmt.Sprintf("[%s]:%s", smtpDetails.IPAddress, smtpDetails.Port)
	}

	// fmt.Printf("Trying to connect to %s\n", address)

	var conn net.Conn
	var err error

	// Handle connection based on port
	switch smtpDetails.Port {
	case "465": // SMTPS
		conn, err = tls.DialWithDialer(dialer, "tcp", address, &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         smtpDetails.Server,
		})
	default: // Plain or STARTTLS
		conn, err = dialer.Dial("tcp", address)
	}

	if err != nil {
		return result, fmt.Errorf("connection failed: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, smtpDetails.Server)
	if err != nil {
		return result, fmt.Errorf("SMTP client creation failed: %v", err)
	}
	defer client.Close()

	// HELO/EHLO
	if err = client.Hello(localName); err != nil {
		return result, fmt.Errorf("HELO failed: %v", err)
	}

	// STARTTLS if available and not already TLS
	if smtpDetails.Port != "465" && useTLS {
		if ok, _ := client.Extension("STARTTLS"); ok {
			config := &tls.Config{
				InsecureSkipVerify: true,
				ServerName:         smtpDetails.Server,
			}
			if err = client.StartTLS(config); err != nil {
				// fmt.Printf("STARTTLS failed: %v\n", err)
				fmt.Printf("STARTTLS failed: %v\n", err)
			}
		}
	}

	// MAIL FROM
	if err = client.Mail(c.SenderEmail); err != nil {
		return result, fmt.Errorf("MAIL FROM failed: %v", err)
	}

	// RCPT TO
	err = client.Rcpt(recipientEmail)
	client.Quit()

	if err != nil {
		if strings.Contains(err.Error(), "450 4.7.1") {
			result.IsValid = true
			result.ErrorMessage = "Reverse DNS lookup required but email might be valid"
			return result, nil
		}

		if strings.Contains(err.Error(), "550 5.1.1") {
			result.ErrorMessage = "User doesn't exist"
			return result, nil
		}

		if strings.Contains(err.Error(), "250") {
			result.IsValid = true
			result.IsCatchAll = true
			return result, nil
		}

		return result, err
	}

	result.IsValid = true
	return result, nil
}

// ValidateEmail validates the recipient's email address by checking its format,
// verifying the existence of MX records for the domain, and attempting to connect
// to the mail servers using SMTP.
//
// Parameters:
//   - recipientEmail: The email address of the recipient to be validated.
//   - senderEmail: The email address of the sender.
//
// Returns:
//   - *ValidationResult: A struct containing the validation result, including whether
//     the email is valid, if MX records were found, and any error messages.
//   - error: An error object if an error occurred during the validation process.
//
// The function performs the following steps:
//  1. Checks if the recipient email contains an "@" symbol and splits it into local
//     and domain parts.
//  2. Retrieves the MX records for the domain.
//  3. Gets the local hostname for the HELO command.
//  4. Attempts to connect to each mail server using SMTP, first without TLS and then
//     with TLS if the initial attempt fails.
//  5. Returns the validation result and any errors encountered during the process.
func (c *Client) ValidateEmail(recipientEmail string) (*ValidationResult, error) {
	// Basic format validation
	if !strings.Contains(recipientEmail, "@") {
		return &ValidationResult{
			IsValid:      false,
			ErrorMessage: "Invalid email format",
		}, nil
	}

	parts := strings.Split(recipientEmail, "@")
	if len(parts) != 2 {
		return &ValidationResult{
			IsValid:      false,
			ErrorMessage: "Invalid email format",
		}, nil
	}

	domain := parts[1]
	// fmt.Printf("Validating email domain: %s\n", domain)

	// Check MX records
	mailServers, err := c.GetMailServers(domain)
	if err != nil {
		return &ValidationResult{
			IsValid:      false,
			HasMX:        false,
			ErrorMessage: "No MX records found",
		}, nil
	}

	// Get hostname for HELO
	localName, err := c.GetHostname()
	if err != nil {
		return &ValidationResult{
			IsValid:      false,
			HasMX:        true,
			ErrorMessage: err.Error(),
		}, nil
	}
	// fmt.Printf("Using hostname for HELO: %s\n", localName)

	// Try each mail server
	var lastErr error
	for _, mailServer := range mailServers {
		smtpServer, err := c.GetSMTPServer(mailServer)
		if err != nil {
			lastErr = err
			continue
		}

		// fmt.Printf("Trying mail server: %s\n", mailServer)
		// fmt.Printf("SMTP server details: %+v\n", smtpServer)

		// try connecting with TLS
		result, err := c.TryConnectingSMTP(smtpServer, recipientEmail, localName, false)
		if err == nil {
			result.SMTPDetails = smtpServer
			return result, nil
		}
		// fmt.Printf("Validation attempt without TLS failed for server %s: %v\n", mailServer, err)
		// fmt.Println("trying to connect with TLS...")

		// Try connecting with TLS
		result, err = c.TryConnectingSMTP(smtpServer, recipientEmail, localName, true)
		if err == nil {
			result.SMTPDetails = smtpServer
			return result, nil
		}

		// fmt.Printf("Validation attempt with TLS failed for server %s: %v\n", mailServer, err)

		lastErr = err
	}

	return &ValidationResult{
		IsValid:      false,
		HasMX:        true,
		ErrorMessage: lastErr.Error(),
	}, nil
}

// Helper function to format validation results
// FormatValidationResult formats the validation result of an email address into a human-readable string.
//
// Parameters:
//   - email: The email address that was validated.
//   - result: A pointer to a ValidationResult struct containing the validation details.
//
// Returns:
//
//	A formatted string summarizing the validation results, including the email address, validation status,
//	presence of MX records, catch-all status, and any error message.
func (c *Client) FormatValidationResult(recipientEmail string, result *ValidationResult) string {
	status := "INVALID"
	if result.IsValid {
		status = "VALID"
	}

	return fmt.Sprintf(`
Email Validation Results for %s:
Status: %s
Has MX Records: %v
Catch-All: %v
Details: %s
`, recipientEmail, status, result.HasMX, result.IsCatchAll, result.ErrorMessage)
}

// ExtractDomainFromEmailAddress extracts the domain part from the given email address.
// It takes a recipient email as input and returns the domain as a string.
// If the email format is invalid, it returns an error.
//
// Parameters:
//   receipientEmail (string): The email address from which to extract the domain.
//
// Returns:
//   string: The domain part of the email address.
//   error: An error if the email format is invalid.
func (c *Client) ExtractDomainFromEmailAddress(receipientEmail string) (string, error) {

	parts := strings.Split(receipientEmail, "@")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid email format")
	}

	domain := parts[1]
	return domain, nil
}
