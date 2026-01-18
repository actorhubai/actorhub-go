// Package actorhub provides a Go client for the ActorHub.ai API.
//
// ActorHub.ai helps protect identities from unauthorized AI-generated content.
// This SDK provides methods to verify images against protected identities,
// check consent status, browse the marketplace, and purchase licenses.
//
// # Quick Start
//
//	client := actorhub.NewClient("your-api-key")
//
//	result, err := client.Verify(context.Background(), &actorhub.VerifyRequest{
//	    ImageURL: "https://example.com/image.jpg",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	if result.Protected {
//	    fmt.Println("Protected identity detected!")
//	}
//
// # Configuration
//
// The client can be configured with various options:
//
//	client := actorhub.NewClient("your-api-key",
//	    actorhub.WithBaseURL("https://custom.actorhub.ai"),
//	    actorhub.WithTimeout(60 * time.Second),
//	    actorhub.WithMaxRetries(5),
//	)
//
// # Error Handling
//
// The SDK returns typed errors for different scenarios:
//
//   - AuthenticationError: Invalid or missing API key (401)
//   - RateLimitError: Rate limit exceeded (429)
//   - ValidationError: Request validation failed (422)
//   - NotFoundError: Resource not found (404)
//   - ServerError: Server error (5xx)
//
// Example:
//
//	result, err := client.Verify(ctx, req)
//	if err != nil {
//	    switch e := err.(type) {
//	    case *actorhub.AuthenticationError:
//	        fmt.Println("Invalid API key")
//	    case *actorhub.RateLimitError:
//	        fmt.Printf("Rate limit exceeded, retry after %d seconds\n", e.RetryAfter)
//	    default:
//	        fmt.Println("Error:", err)
//	    }
//	}
package actorhub
