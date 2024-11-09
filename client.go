package mailify

// 
// Client represents an email client with a sender email address.
type Client struct {
	SenderEmail string
}

// NewClient creates a new Client instance with the provided sender email address.
// It returns a pointer to the Client and an error, if any.
//
// Parameters:
//   - SenderEmail: A string representing the sender's email address.
//
// Returns:
//   - *Client: A pointer to the newly created Client instance.
//   - error: An error if there is any issue during the creation of the Client.
func NewClient(SenderEmail string) (*Client, error) {
	return &Client{
		SenderEmail: SenderEmail,
	}, nil
}

