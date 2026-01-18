// Example usage of the ActorHub Go SDK.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	actorhub "github.com/actorhub/actorhub-go"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("ACTORHUB_API_KEY")
	if apiKey == "" {
		log.Fatal("ACTORHUB_API_KEY environment variable is required")
	}

	// Create client
	client := actorhub.NewClient(apiKey)
	ctx := context.Background()

	// Example 1: Verify an image
	fmt.Println("=== Verifying Image ===")
	verifyResult, err := client.Verify(ctx, &actorhub.VerifyRequest{
		ImageURL:              "https://example.com/image.jpg",
		IncludeLicenseOptions: true,
	})
	if err != nil {
		log.Printf("Verify error: %v", err)
	} else {
		fmt.Printf("Protected: %v\n", verifyResult.Protected)
		fmt.Printf("Faces detected: %d\n", verifyResult.FacesDetected)
		for _, identity := range verifyResult.Identities {
			if identity.DisplayName != nil {
				fmt.Printf("  - Identity: %s (similarity: %.2f%%)\n",
					*identity.DisplayName,
					*identity.SimilarityScore*100)
			}
		}
	}

	// Example 2: Check consent
	fmt.Println("\n=== Checking Consent ===")
	consentResult, err := client.CheckConsent(ctx, &actorhub.ConsentCheckRequest{
		ImageURL:    "https://example.com/face.jpg",
		Platform:    "runway",
		IntendedUse: "video",
		Region:      "US",
	})
	if err != nil {
		log.Printf("Consent check error: %v", err)
	} else {
		fmt.Printf("Protected: %v\n", consentResult.Protected)
		for _, face := range consentResult.Faces {
			fmt.Printf("  - Video generation allowed: %v\n", face.Consent.VideoGeneration)
			fmt.Printf("  - Commercial use allowed: %v\n", face.Consent.CommercialUse)
			fmt.Printf("  - License available: %v\n", face.License.Available)
		}
	}

	// Example 3: List marketplace
	fmt.Println("\n=== Marketplace Listings ===")
	listings, err := client.ListMarketplace(ctx, &actorhub.MarketplaceListRequest{
		Category: "ACTOR",
		SortBy:   "popular",
		Limit:    5,
	})
	if err != nil {
		log.Printf("Marketplace error: %v", err)
	} else {
		fmt.Printf("Found %d listings:\n", len(listings))
		for _, listing := range listings {
			fmt.Printf("  - %s: $%.2f (%s)\n",
				listing.Title,
				listing.BasePriceUSD,
				listing.Category)
		}
	}

	// Example 4: Error handling
	fmt.Println("\n=== Error Handling Example ===")
	_, err = client.GetIdentity(ctx, "nonexistent-id")
	if err != nil {
		switch e := err.(type) {
		case *actorhub.NotFoundError:
			fmt.Printf("Not found: %s\n", e.Message)
		case *actorhub.AuthenticationError:
			fmt.Printf("Auth error: %s\n", e.Message)
		case *actorhub.RateLimitError:
			fmt.Printf("Rate limited, retry after %d seconds\n", e.RetryAfter)
		default:
			fmt.Printf("Other error: %v\n", err)
		}
	}

	fmt.Println("\nDone!")
}
