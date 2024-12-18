package main

import (
	"log"
	"fmt"

	"github.com/adarsh-jaiss/mailify"
)

func main() {
	// Create a new mailify client
	senderEmail := "its.adarshjaiss@gmail.com"
	receipientEmail := "hello@namanrai.tech"

	client, err := mailify.NewClient(senderEmail)
	if err != nil {
		log.Fatalf("Failed to create mailify client: %v", err)
	}

	// Get mail servers for a domain
	resp, err := client.GetMailServers("namanrai.tech")
	if err != nil {
		log.Fatalf("Failed to get mail servers: %v", err)
	}
	log.Println("Mail servers:", resp)

	// Get mail servers for a recepient email
	res, err := client.GetMailServersFromReceipientEmail(receipientEmail)
	if err != nil {
		log.Fatalf("Failed to get mail servers: %v", err)
	}
	log.Println("Mail servers:", res)

	// Validate an email address
	result, err := client.ValidateEmail(receipientEmail)
	if err!= nil {
		log.Fatalf("Failed to validate email: %v", err)
	}

	fmt.Println("Validation result:", client.FormatValidationResult(receipientEmail,result))

	// Validate all the email address in an Excel file, creates a new column with the validation result
	err = client.ProcessAndValidateEmailsViaExcel("emails.xlsx",client.SenderEmail)
	if err!= nil {
         fmt.Printf("Error processing file: %v\n", err)
         return
	}

}
