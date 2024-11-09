package mailify

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

// GetMailServers retrieves the mail servers (MX records) for a given domain.
// It uses a custom DNS resolver that queries Google's public DNS server (8.8.8.8).

// Parameters:
//   - domain: The domain name for which to look up MX records.

// Returns:
//   - A slice of strings containing the mail server hostnames.
//   - An error if there was an issue looking up the MX records.

// Example:
//   mailServers, err := GetMailServers("example.com")
//   if err != nil {
//       log.Fatalf("Failed to get mail servers: %v", err)
//   }
//   fmt.Println("Mail servers:", mailServers)

func(c *Client) GetMailServers(domain string) ([]string, error) {
	// Use custom DNS resolver to query Google's public DNS server

	resolver := net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, network, "8.8.8.8:53") // Use Google DNS
		},
	}

	// Lookup MX records for the domain
	mx, err := resolver.LookupMX(context.Background(), domain)
	// mx, err := net.LookupMX(domain)
	if err != nil {
		return nil, fmt.Errorf("error looking up MX records: %v", err)
	}

	// Extract mail server hostnames
	var mailServers []string
	for _, record := range mx {
		mailServers = append(mailServers, strings.TrimSuffix(record.Host, "."))
	}

	// Print mail servers
	// fmt.Printf("Found mail servers for %s: %v\n", domain, mailServers)
	return mailServers, nil
}

// GetSMTPServer attempts to find an available SMTP server for the given mail server.
// It performs a DNS lookup to get all IP addresses (both IPv4 and IPv6) associated with the mail server,
// and then tries to connect to common SMTP ports (587, 25, 465) on each IP address.
//
// If a connection is successfully established, it returns the SMTP server details including
// the server name, port, protocol, and IP address. If no available SMTP servers are found,
// it returns an error.
//
// Parameters:
//   - mailServer: The domain name of the mail server to look up.
//
// Returns:
//   - *SMTPDetails: A struct containing the details of the SMTP server if found.
//   - error: An error if no available SMTP servers are found or if there is a lookup failure.
func(c *Client) GetSMTPServer(mailServer string) (*SMTPDetails, error) {
	// Get all IPs (both IPv4 and IPv6)
	ips, err := net.LookupIP(mailServer)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup IP for %s: %v", mailServer, err)
	}

	// Try each IP address
	for _, ip := range ips {
		// Try common SMTP ports
		ports := []string{"587", "25", "465"}
		for _, port := range ports {
			// Format address based on IP version
			var address string
			if ip.To4() != nil {
				// IPv4
				address = fmt.Sprintf("%s:%s", ip.String(), port)
			} else {
				// IPv6 - wrap in square brackets
				address = fmt.Sprintf("[%s]:%s", ip.String(), port)
			}

			// Set timeout for connection
			smtpTimeout := time.Duration(time.Second * 5)

			// Try to connect
			conn, err := net.DialTimeout("tcp", address, smtpTimeout)
			if err != nil {
				continue
			}
			defer conn.Close()

			return &SMTPDetails{
				Server:    mailServer,
				Port:      port,
				Protocol:  "SMTP",
				IPAddress: ip.String(),
			}, nil
		}
	}
	return nil, fmt.Errorf("no available SMTP servers found for %s", mailServer)
}

// GetMailServersFromReceipientEmail extracts the domain from the given email address
// and retrieves the mail servers associated with that domain.
//
// Parameters:
//   email (string): The recipient's email address.
//
// Returns:
//   []string: A slice of mail server addresses.
//   error: An error object if there was an issue extracting the domain or retrieving the mail servers.
func(c *Client) GetMailServersFromReceipientEmail(email string) ([]string, error) {
	// Extract domain from email address
	domain,err := c.ExtractDomainFromEmailAddress(email)
	if err != nil {
		return nil, fmt.Errorf("error extracting domain from email address: %v", err)
	}
	
	return c.GetMailServers(domain)
}
