package expr

import "github.com/rs/xid"

// NewID generates a random ID guaranteed to be unique (enough)
func NewID() string {
	return xid.New().String()
}
