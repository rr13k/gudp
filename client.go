package gudp

import (
	"bytes"
	"net"
	"sync"
	"time"
)

type Client struct {
	ID string

	sessionID []byte // Session ID is a secret byte array that indicates the client is already done the handshake process, the client must prepend these bytes into the start of each record body before encryption

	addr *net.UDPAddr // UDP address of the client

	eKey []byte // Client encryption key to decrypt & encrypt a record body with the symmetric encryption algorithm

	lastHeartbeat *time.Time // last time when a record has been received from the client.

	sync.Mutex
}

// ValidateSessionID compares the client session ID with the given one
func (c *Client) ValidateSessionID(sessionID []byte) bool {
	return bytes.Equal(c.sessionID, sessionID)
}
