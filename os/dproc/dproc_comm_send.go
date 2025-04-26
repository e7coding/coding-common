// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dproc

import (
	"github.com/coding-common/errs/werr"
	"io"

	"github.com/coding-common/internal/json"
	"github.com/coding-common/net/dtcp"
)

// Send sends data to specified process of given pid.
func Send(pid int, data []byte, group ...string) error {
	msg := MsgRequest{
		SenderPid:   Pid(),
		ReceiverPid: pid,
		Group:       defaultGroupNameForProcComm,
		Data:        data,
	}
	if len(group) > 0 {
		msg.Group = group[0]
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	var conn *dtcp.PoolConn
	conn, err = getConnByPid(pid)
	if err != nil {
		return err
	}
	defer conn.Close()
	// Do the sending.
	var result []byte
	result, err = conn.SendRecvPkg(msgBytes, dtcp.PkgOption{
		Retry: dtcp.Retry{
			Count: 3,
		},
	})
	if len(result) > 0 {
		response := new(MsgResponse)
		if err = json.UnmarshalUseNumber(result, response); err == nil {
			if response.Code != 1 {
				err = werr.WithMsg(response.Message)
			}
		}
	}
	// EOF is not really an error.
	if err == io.EOF {
		err = nil
	}
	return err
}
