package authn

import "context"

// Principal models the authenticated caller that endpoint middleware can place
// into request context once Clerk integration is added.
type Principal interface {
	Subject() string
}

// Verifier is the seam between HTTP transport and a concrete auth provider.
// Clerk-backed verification should be added behind this contract later.
type Verifier interface {
	Verify(context.Context, string) (Principal, error)
}
