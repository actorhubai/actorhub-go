# ActorHub Go SDK

Official Go SDK for [ActorHub.ai](https://actorhub.ai) - Verify AI-generated content against protected identities.

## Installation

```bash
go get github.com/actorhubai/actorhub-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    actorhub "github.com/actorhubai/actorhub-go"
)

func main() {
    // Initialize the client
    client := actorhub.NewClient("your-api-key")

    // Verify if an image contains protected identities
    result, err := client.Verify(context.Background(), &actorhub.VerifyRequest{
        ImageURL: "https://example.com/image.jpg",
    })
    if err != nil {
        log.Fatal(err)
    }

    if result.Protected {
        fmt.Println("Protected identity detected!")
        for _, identity := range result.Identities {
            if identity.DisplayName != nil {
                fmt.Printf("  - %s (similarity: %.2f)\n", *identity.DisplayName, *identity.SimilarityScore)
            }
        }
    }
}
```

## Features

- **Identity Verification**: Check if images contain protected identities
- **Consent Checking**: Verify consent before AI generation
- **Marketplace Access**: Browse and license identities
- **Automatic Retries**: Built-in retry logic with exponential backoff
- **Context Support**: Full context.Context support for cancellation
- **Typed Errors**: Specific error types for easy handling

## Usage Examples

### Verify Image

```go
// From URL
result, err := client.Verify(ctx, &actorhub.VerifyRequest{
    ImageURL: "https://example.com/image.jpg",
})

// From base64
result, err := client.Verify(ctx, &actorhub.VerifyRequest{
    ImageBase64: "base64-encoded-data...",
})

fmt.Printf("Protected: %v\n", result.Protected)
fmt.Printf("Faces detected: %d\n", result.FacesDetected)
```

### Check Consent (for AI Platforms)

```go
result, err := client.CheckConsent(ctx, &actorhub.ConsentCheckRequest{
    ImageURL:    "https://example.com/face.jpg",
    Platform:    "runway",
    IntendedUse: "video",
    Region:      "US",
})

if result.Protected {
    for _, face := range result.Faces {
        fmt.Printf("Consent for video: %v\n", face.Consent.VideoGeneration)
        fmt.Printf("License available: %v\n", face.License.Available)
    }
}
```

### Browse Marketplace

```go
// Search listings
featured := true
listings, err := client.ListMarketplace(ctx, &actorhub.MarketplaceListRequest{
    Category: "ACTOR",
    Featured: &featured,
    SortBy:   "popular",
    Limit:    10,
})

for _, listing := range listings {
    fmt.Printf("%s - $%.2f\n", listing.Title, listing.BasePriceUSD)
}
```

### Purchase License

```go
purchase, err := client.PurchaseLicense(ctx, &actorhub.PurchaseLicenseRequest{
    IdentityID:         "uuid-here",
    LicenseType:        string(actorhub.LicenseTypeStandard),
    UsageType:          string(actorhub.UsageTypeCommercial),
    ProjectName:        "My AI Project",
    ProjectDescription: "Creating promotional content",
    DurationDays:       30,
})

// Redirect user to Stripe checkout
fmt.Printf("Checkout URL: %s\n", purchase.CheckoutURL)
```

### Get My Licenses

```go
licenses, err := client.GetMyLicenses(ctx, "active", 1, 20)

for _, license := range licenses {
    fmt.Printf("%s - %s - Expires: %v\n",
        license.IdentityName,
        license.LicenseType,
        license.ExpiresAt)
}
```

## Error Handling

```go
import "errors"

result, err := client.Verify(ctx, req)
if err != nil {
    var authErr *actorhub.AuthenticationError
    var rateLimitErr *actorhub.RateLimitError
    var validationErr *actorhub.ValidationError
    var notFoundErr *actorhub.NotFoundError

    switch {
    case errors.As(err, &authErr):
        fmt.Println("Invalid API key")
    case errors.As(err, &rateLimitErr):
        fmt.Printf("Rate limit exceeded. Retry after: %d seconds\n", rateLimitErr.RetryAfter)
    case errors.As(err, &validationErr):
        fmt.Printf("Validation error: %s\n", validationErr.Message)
    case errors.As(err, &notFoundErr):
        fmt.Println("Resource not found")
    default:
        fmt.Printf("Error: %v\n", err)
    }
}
```

## Configuration

```go
import "time"

client := actorhub.NewClient("your-api-key",
    actorhub.WithBaseURL("https://custom.actorhub.ai"),
    actorhub.WithTimeout(60 * time.Second),
    actorhub.WithMaxRetries(5),
)
```

## API Reference

### Client Methods

| Method | Description |
|--------|-------------|
| `Verify()` | Verify if image contains protected identities |
| `GetIdentity()` | Get identity details by ID |
| `CheckConsent()` | Check consent status for AI generation |
| `ListMarketplace()` | Search marketplace listings |
| `GetMyLicenses()` | Get user's purchased licenses |
| `PurchaseLicense()` | Purchase a license |
| `GetActorPack()` | Get Actor Pack status |

## Requirements

- Go 1.21+

## License

MIT License - see [LICENSE](LICENSE) for details.

## Links

- [Documentation](https://docs.actorhub.ai)
- [API Reference](https://api.actorhub.ai/docs)
- [GitHub](https://github.com/actorhubai/actorhub-go)
- [pkg.go.dev](https://pkg.go.dev/github.com/actorhubai/actorhub-go)
