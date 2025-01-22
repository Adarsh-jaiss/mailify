# Mailify

Mailify is a Go package for validating email addresses by checking their format, verifying the existence of MX records for the domain, and attempting to connect to the mail servers using SMTP.

## CLI 

For more details, checkout [CLI DOCS](https://github.com/Adarsh-jaiss/mailify/blob/main/cli/README.md)


## Installation

To install the package, run:

```sh
go get github.com/adarsh-jaiss/mailify

```

## Usage

### Creating a Client

To create a new client, use the NewClient function:

```go
client, err := mailify.NewClient("sender@example.com")
if err != nil {
    log.Fatalf("Failed to create mailify client: %v", err)
}
```

### Validating an Email Address

To validate an email address, use the ValidateEmail method:

```go
result, err := client.ValidateEmail("recipient@example.com")
if err != nil {
    log.Fatalf("Failed to validate email: %v", err)
}

fmt.Println("Validation result:", client.FormatValidationResult("recipient@example.com", result))
```

### Getting Mail Servers
To get the mail servers for a domain, use the GetMailServers method:

```go
mailServers, err := client.GetMailServers("example.com")
if err != nil {
    log.Fatalf("Failed to get mail servers: %v", err)
}

fmt.Println("Mail servers:", mailServers)
```

To get the mail servers for a recipient email, use the GetMailServersFromReceipientEmail method:

```go

mailServers, err := client.GetMailServersFromReceipientEmail("recipient@example.com")
if err != nil {
    log.Fatalf("Failed to get mail servers: %v", err)
}

fmt.Println("Mail servers:", mailServers)
```


### Validate all the email addresses in an Excel file

This section demonstrates how to validate all email addresses in an Excel file using the `ProcessAndValidateEmailsViaExcel` method. The method takes the path to the Excel file and the sender's email as parameters. If there is an error during the processing of the file, it will print an error message and terminate the execution.

Example usage:

```go
	err = client.ProcessAndValidateEmailsViaExcel("emails.xlsx",client.SenderEmail)
	if err!= nil {
         fmt.Printf("Error processing file: %v\n", err)
         return
	}
```
### Example
Here is a complete example demonstrating how to use the package : [check examples](https://github.com/Adarsh-jaiss/mailify/blob/main/example/main.go)
