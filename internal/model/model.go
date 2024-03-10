// Data models
package model

// Profile represents a user profile with an associated email and IP address.
// PendingIP is the IP address the user is attempting to associate with their account.
// AssociatedIP is the IP address currently associated with the user's account.
type Profile struct {
	Email        string `json:"email"`
	PendingIP    string `json:"pending_ip"`
	AssociatedIP string `json:"associated_ip"`
}
