package mailify


// SMTPDetails holds the details required to connect to an SMTP server.
type SMTPDetails struct {
	// Server is the address of the SMTP server.
	Server string
	// Port is the port number on which the SMTP server is listening.
	Port string
	// Protocol is the protocol used by the SMTP server (e.g., "SMTP", "SMTPS").
	Protocol string
	// UsedTLS indicates whether TLS is used for the connection.
	UsedTLS bool
	// IPAddress is the IP address of the SMTP server.
	IPAddress string
}

// ValidationResult represents the result of an email validation check.
type ValidationResult struct {
	// IsValid indicates whether the email address is valid.
	IsValid bool
	// IsCatchAll indicates whether the domain has a catch-all address.
	IsCatchAll bool
	// HasMX indicates whether the domain has MX records.
	HasMX bool
	// ErrorMessage contains any error message encountered during validation.
	ErrorMessage string
	// SMTPDetails contains the SMTP server details used for validation.
	SMTPDetails *SMTPDetails
}

