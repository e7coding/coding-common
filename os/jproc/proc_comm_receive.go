// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jproc

import (
	"fmt"
	"github.com/e7coding/coding-common/errs/jerr"
	"net"

	"github.com/e7coding/coding-common/container/jqueue"

	"github.com/e7coding/coding-common/container/jatomic"

	"github.com/e7coding/coding-common/internal/json"
	"github.com/e7coding/coding-common/jutil/jconv"
	"github.com/e7coding/coding-common/net/jtcp"
	"github.com/e7coding/coding-common/os/jfile"
)

var (
	// tcpListened marks whether the receiving listening service started.
	tcpListened = jatomic.NewBool()
)

// Receive blocks and receives message from other process using local TCP listening.
// Note that, it only enables the TCP listening service when this function called.
func Receive(group ...string) *MsgRequest {
	// Use atomic operations to guarantee only one receiver goroutine listening.
	if tcpListened.CAS(false, true) {
		go receiveTcpListening()
	}
	var groupName string
	if len(group) > 0 {
		groupName = group[0]
	} else {
		groupName = defaultGroupNameForProcComm
	}
	queues := commReceiveQueues.GetOrPutFunc(groupName, func() interface{} {
		return jqueue.New(maxLengthForProcMsgQueue)
	}).(*jqueue.Queue)

	// Blocking receiving.
	if v := queues.Pop(); v != nil {
		return v.(*MsgRequest)
	}
	return nil
}

// receiveTcpListening scans local for available port and starts listening.
func receiveTcpListening() {
	var (
		listen  *net.TCPListener
		conn    net.Conn
		port    = jtcp.MustGetFreePort()
		address = fmt.Sprintf("127.0.0.1:%d", port)
	)
	tcpAddress, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		panic(jerr.WithMsgErr(err, `net.ResolveTCPAddr failed`))
	}
	listen, err = net.ListenTCP("tcp", tcpAddress)
	if err != nil {
		panic(jerr.WithMsgErrF(err, `net.ListenTCP failed for address "%s"`, address))
	}
	// Save the port to the pid file.
	if err = jfile.PutContents(getCommFilePath(Pid()), jconv.String(port)); err != nil {
		panic(err)
	}
	// Start listening.
	for {
		if conn, err = listen.Accept(); err != nil {

		} else if conn != nil {
			go receiveTcpHandler(jtcp.NewConnByNetConn(conn))
		}
	}
}

// receiveTcpHandler is the connection handler for receiving data.
func receiveTcpHandler(conn *jtcp.Conn) {
	var (
		result   []byte
		response MsgResponse
	)
	for {
		response.Code = 0
		response.Message = ""
		response.Data = nil
		buffer, err := conn.RecvPkg()
		if len(buffer) > 0 {
			// Package decoding.
			msg := new(MsgRequest)
			if err = json.UnmarshalUseNumber(buffer, msg); err != nil {
				continue
			}
			if msg.ReceiverPid != Pid() {
				// Not mine package.
				response.Message = fmt.Sprintf(
					"receiver pid not match, target: %d, current: %d",
					msg.ReceiverPid, Pid(),
				)
			} else if v := commReceiveQueues.Get(msg.Group); v == nil {
				// Group check.
				response.Message = fmt.Sprintf("group [%s] does not exist", msg.Group)
			} else {
				// Push to buffer queue.
				response.Code = 1
				v.(*jqueue.Queue).Push(msg)
			}
		} else {
			// Empty package.
			response.Message = "empty package"
		}
		if err == nil {
			result, err = json.Marshal(response)
			if err != nil {
				fmt.Println(err)
			}
			if err = conn.SendPkg(result); err != nil {
				fmt.Println(err)
			}
		} else {
			// Just close the connection if any error occurs.
			if err = conn.Close(); err != nil {
				fmt.Println(err)
			}
			break
		}
	}
}
