package options

import "github.com/nmorenor/chezmoi/net"

type SessionOptions struct {
	Username *string
	HostMode bool
	Session  *string
	Client   *net.RemoteClient
}

var (
	SessionInfo = &SessionOptions{Username: nil, HostMode: false, Session: nil, Client: nil}
)
