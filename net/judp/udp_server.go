// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package judp

import (
	"fmt"
	"github.com/e7coding/coding-common/container/jmap"
	"github.com/e7coding/coding-common/errs/jerr"
	"net"
	"sync"

	"github.com/e7coding/coding-common/jutil/jconv"
	"github.com/e7coding/coding-common/text/jstr"
)

const (
	// FreePortAddress marks the server listens using random free port.
	FreePortAddress = ":0"
)

const (
	defaultServer = "default"
)

// Server is the UDP server.
type Server struct {
	// Used for Server.listen concurrent safety.
	// The golang test with data race checks this.
	mu sync.Mutex

	// UDP server connection object.
	conn *ServerConn

	// UDP server listening address.
	address string

	// Handler for UDP connection.
	handler ServerHandler
}

// ServerHandler handles all server connections.
type ServerHandler func(conn *ServerConn)

var (
	// serverMapping is used for instance name to its UDP server mappings.
	serverMapping = jmap.NewSafeStrAnyMap()
)

// GetServer creates and returns an udp server instance with given name.
func GetServer(name ...interface{}) *Server {
	serverName := defaultServer
	if len(name) > 0 && name[0] != "" {
		serverName = jconv.String(name[0])
	}
	if s := serverMapping.Get(serverName); s != nil {
		return s.(*Server)
	}
	s := NewServer("", nil)
	serverMapping.Put(serverName, s)
	return s
}

// NewServer creates and returns an udp server.
// The optional parameter `name` is used to specify its name, which can be used for
// GetServer function to retrieve its instance.
func NewServer(address string, handler ServerHandler, name ...string) *Server {
	s := &Server{
		address: address,
		handler: handler,
	}
	if len(name) > 0 && name[0] != "" {
		serverMapping.Put(name[0], s)
	}
	return s
}

// SetAddress sets the server address for UDP server.
func (s *Server) SetAddress(address string) {
	s.address = address
}

// SetHandler sets the connection handler for UDP server.
func (s *Server) SetHandler(handler ServerHandler) {
	s.handler = handler
}

// Close closes the connection.
// It will make server shutdowns immediately.
func (s *Server) Close() (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	err = s.conn.Close()
	if err != nil {
		err = jerr.WithMsgErr(err, "connection failed")
	}
	return
}

// Run starts listening UDP connection.
func (s *Server) Run() error {
	if s.handler == nil {
		return jerr.WithMsg(
			"start running failed: socket handler not defined",
		)
	}
	addr, err := net.ResolveUDPAddr("udp", s.address)
	if err != nil {
		err = jerr.WithMsgErrF(err, `net.ResolveUDPAddr failed for address "%s"`, s.address)
		return err
	}
	listenedConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		err = jerr.WithMsgErrF(err, `net.ListenUDP failed for address "%s"`, s.address)
		return err
	}
	s.mu.Lock()
	s.conn = NewServerConn(listenedConn)
	s.mu.Unlock()
	s.handler(s.conn)
	return nil
}

// GetListenedAddress retrieves and returns the address string which are listened by current server.
func (s *Server) GetListenedAddress() string {
	if !jstr.Contains(s.address, FreePortAddress) {
		return s.address
	}
	var (
		address      = s.address
		listenedPort = s.GetListenedPort()
	)
	address = jstr.Replace(address, FreePortAddress, fmt.Sprintf(`:%d`, listenedPort))
	return address
}

// GetListenedPort retrieves and returns one port which is listened to by current server.
func (s *Server) GetListenedPort() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ln := s.conn; ln != nil {
		return ln.LocalAddr().(*net.UDPAddr).Port
	}
	return -1
}
