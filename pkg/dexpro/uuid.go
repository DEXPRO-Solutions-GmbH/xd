package dexpro

import "github.com/google/uuid"

// GenerateProjectID generates a UUID from a domain name.
//
// This panics if the hardcoded namespace UUID is not valid. That should never happen. If it does,
// contact a library maintainer.
func GenerateProjectID(domain string) uuid.UUID {
	// The namespace is a predefined UUID. This will not change and is thus hardcoded.
	namespace := uuid.MustParse("4befef11-1d7c-5958-8343-7aab6779a086")

	return uuid.NewSHA1(namespace, []byte(domain))
}
