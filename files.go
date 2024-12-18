package mailify

import (
	"strings"
	"fmt"

	"github.com/xuri/excelize/v2"
)

// ProcessAndValidateEmails reads an Excel file, validates emails, and writes results back
// ProcessAndValidateEmailsViaExcel processes and validates emails from an Excel file.
// It reads the email addresses from the specified Excel file, validates each email,
// and writes the validation results back to the Excel file in a new column.
//
// Parameters:
//   - filename: The path to the Excel file containing the email addresses.
//   - senderEmail: The email address of the sender (not used in the current implementation).
//
// Returns:
//   - error: An error if any issue occurs during the process, otherwise nil.
//
// The function performs the following steps:
//   1. Opens the specified Excel file.
//   2. Reads all rows from the first sheet ("Sheet1").
//   3. Creates a map of headers from the first row.
//   4. Adds a new column header for email validation results if it doesn't exist.
//   5. Iterates over each row, validates the email address, and writes the validation result to the new column.
//   6. Saves the modified Excel file with the validation results.
//
// The function prints progress and summary information to the console.
func(c *Client) ProcessAndValidateEmailsViaExcel(filename string, senderEmail string) error {
	fmt.Println("\n=== Starting Email Validation Process ===")

	// Open the Excel file
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Warning: failed to close excel file: %v\n", err)
		}
	}()

	fmt.Printf("Successfully opened Excel file: %s\n", filename)

	// Get all the rows in Sheet1
	rows, err := f.GetRows("sheet1")
	if err != nil {
		return fmt.Errorf("failed to get rows: %w", err)
	}

	if len(rows) < 2 {
		return fmt.Errorf("excel file has no data except field names")
	}

	fmt.Printf("Found %d rows in the Excel file (including header)\n", len(rows))

	// Create headers map and add new column
	headers := make(map[string]int)
	for i, cell := range rows[0] {
		header := strings.ToLower(strings.ReplaceAll(cell, " ", "_"))
		headers[header] = i
	}

	// Add new column for email validation if it doesn't exist
	isValidEmailCol := len(rows[0])
	headers["is_valid_email"] = isValidEmailCol

	// Add the new column header
	err = f.SetCellValue("Sheet1", fmt.Sprintf("%s1", columnToLetter(isValidEmailCol)), "is_valid_email")
	if err != nil {
		return fmt.Errorf("failed to add header: %w", err)
	}

	fmt.Println("\nStarting email validation process...")
	fmt.Println("=====================================")

	validCount := 0
	invalidCount := 0

	// Process each row
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) == 0 {
			continue
		}

		// Get email from the row
		var email string
		if idx, ok := headers["email"]; ok && idx < len(row) {
			email = strings.TrimSpace(row[idx])
		}

		if email != "" {
			fmt.Printf("Validating email %d/%d: %s... ", i, len(rows)-1, email)

			// Validate email
			result, err := c.ValidateEmail(email)
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
				continue
			}

			// Write validation result to the new column
			cellRef := fmt.Sprintf("%s%d", columnToLetter(isValidEmailCol), i+1)
			err = f.SetCellValue("Sheet1", cellRef, result.IsValid)
			if err != nil {
				fmt.Printf("ERROR: Failed to write result: %v\n", err)
				continue
			}

			if result.IsValid {
				fmt.Println("VALID ✓")
				validCount++
			} else {
				fmt.Println("INVALID ✗")
				invalidCount++
			}
		}
	}

	// Save the modified Excel file
	fmt.Println("\nSaving results to Excel file...")
	err = f.Save()
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	fmt.Println("\n=== Email Validation Summary ===")
	fmt.Printf("Total emails processed: %d\n", validCount+invalidCount)
	fmt.Printf("Valid emails: %d\n", validCount)
	fmt.Printf("Invalid emails: %d\n", invalidCount)
	fmt.Printf("Results have been written to: %s\n", filename)
	fmt.Println("===============================")

	return nil
}

// columnToLetter converts a given column number (0-indexed) to its corresponding
// Excel-style column letter. For example, 0 -> "A", 1 -> "B", 25 -> "Z", 26 -> "AA", etc.
// 
// Parameters:
//   col (int): The column number to convert.
//
// Returns:
//   string: The corresponding Excel-style column letter.
func columnToLetter(col int) string {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := ""

	for col >= 0 {
		result = string(alphabet[col%26]) + result
		col = col/26 - 1
	}

	return result
}
