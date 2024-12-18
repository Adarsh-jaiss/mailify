package cmd

import (
	"fmt"
	"os"

	"github.com/adarsh-jaiss/mailify"
	"github.com/spf13/cobra"
)

// senderEmail represents the email address of the sender.
var (
	senderEmail    string
	client         *mailify.Client
	emailToCheck   string
	excelFile      string
	domain         string
	receipientEmail string
)

// rootCmd represents the base command for the Mailify CLI tool
// rootCmd represents the base command when called without any subcommands.
// It provides functionality to validate email addresses and get mail server information.
// 
// Usage:
//   mailify [flags]
// 
// Flags:
//   -e, --email string       Email address to validate
//   -x, --excel string       Path to Excel file for bulk email validation
//   -d, --domain string      Domain to get mail servers for
//   -r, --receipient string  Email address to get mail servers for
// 
// Examples:
//   # Validate a single email address
//   mailify --email example@example.com
// 
//   # Bulk validate emails from an Excel file
//   mailify --excel emails.xlsx
// 
//   # Get mail servers for a domain
//   mailify --domain example.com
// 
//   # Get mail servers for an email address
//   mailify --receipient example@example.com
// 
// If no flags are provided, an error will be returned indicating that no operation was specified.
var rootCmd = &cobra.Command{
	Use:   "mailify",
	Short: "Mailify is a CLI tool for email validation and server information",
	Long: `Mailify CLI provides functionality to validate email addresses and get mail server information.
It can process single email addresses or bulk validate emails from Excel files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize client
		var err error
		client, err = mailify.NewClient(senderEmail)
		if err != nil {
			return fmt.Errorf("failed to create mailify client: %v", err)
		}

		// Handle single email validation
		if emailToCheck != "" {
			result, err := client.ValidateEmail(emailToCheck)
			if err != nil {
				return fmt.Errorf("failed to validate email: %v", err)
			}
			fmt.Println(client.FormatValidationResult(emailToCheck, result))
		}

		// Handle bulk validation from Excel
		if excelFile != "" {
			err := client.ProcessAndValidateEmailsViaExcel(excelFile, client.SenderEmail)
			if err != nil {
				return fmt.Errorf("failed to process Excel file: %v", err)
			}
			fmt.Println("Successfully processed and validated emails in", excelFile)
		}

		// Handle domain mail servers
		if domain != "" {
			servers, err := client.GetMailServers(domain)
			if err != nil {
				return fmt.Errorf("failed to get mail servers: %v", err)
			}
			fmt.Println("Mail servers for", domain+":")
			for _, server := range servers {
				fmt.Println("-", server)
			}
		}

		// Handle email mail servers
		if receipientEmail != "" {
			servers, err := client.GetMailServersFromReceipientEmail(receipientEmail)
			if err != nil {
				return fmt.Errorf("failed to get mail servers: %v", err)
			}
			fmt.Println("Mail servers for", receipientEmail+":")
			for _, server := range servers {
				fmt.Println("-", server)
			}
		}

		// Check if no flags were provided
		if emailToCheck == "" && excelFile == "" && domain == "" && receipientEmail == "" {
			return fmt.Errorf("no operation specified. Use --help to see available flags")
		}

		return nil
	},
}

// Execute runs the root command and handles any errors that occur during its execution.
// If an error is encountered, it prints the error message and exits the program with a status code of 1.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// init initializes the command-line flags for the root command.
// It sets up the following flags:
// - sender: Required flag for specifying the sender email address.
// - validate: Optional flag for validating a single email address.
// - excel: Optional flag for processing and validating emails from an Excel file.
// - domain: Optional flag for getting mail servers for a domain.
// - receipient: Optional flag for getting mail servers for a recipient email.
func init() {
	// Required sender email flag
	rootCmd.Flags().StringVarP(&senderEmail, "sender", "s", "", "Sender email address (required)")
	rootCmd.MarkFlagRequired("sender")

	// Operation flags
	rootCmd.Flags().StringVarP(&emailToCheck, "validate", "v", "", "Validate a single email address")
	rootCmd.Flags().StringVarP(&excelFile, "excel", "e", "", "Process and validate emails from an Excel file")
	rootCmd.Flags().StringVarP(&domain, "domain", "d", "", "Get mail servers for a domain")
	rootCmd.Flags().StringVarP(&receipientEmail, "receipient", "r", "", "Get mail servers for a receipient email")
}