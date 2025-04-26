// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dproc

import (
	"context"
	"fmt"
	"github.com/coding-common/container/wmap"
	"github.com/coding-common/errs/werr"
	"sync"

	"github.com/coding-common/internal/intlog"
	"github.com/coding-common/net/dtcp"
	"github.com/coding-common/os/dfile"
	"github.com/coding-common/util/dconv"
)

// MsgRequest is the request structure for process communication.
type MsgRequest struct {
	SenderPid   int    // Sender PID.
	ReceiverPid int    // Receiver PID.
	Group       string // Message group name.
	Data        []byte // Request data.
}

// MsgResponse is the response structure for process communication.
type MsgResponse struct {
	Code    int    // 1: OK; Other: Error.
	Message string // Response message.
	Data    []byte // Response data.
}

const (
	defaultFolderNameForProcComm = "gf_pid_port_mapping" // Default folder name for storing pid to port mapping files.
	defaultGroupNameForProcComm  = ""                    // Default group name.
	defaultTcpPortForProcComm    = 10000                 // Starting port number for receiver listening.
	maxLengthForProcMsgQueue     = 10000                 // Max size for each message queue of the group.
)

var (
	// commReceiveQueues is the group name to queue map for storing received data.
	// The value of the map is type of *gqueue.Queue.
	commReceiveQueues = wmap.NewSafeStrAnyMap()

	// commPidFolderPath specifies the folder path storing pid to port mapping files.
	commPidFolderPath string

	// commPidFolderPathOnce is used for lazy calculation for `commPidFolderPath` is necessary.
	commPidFolderPathOnce sync.Once
)

// getConnByPid creates and returns a TCP connection for specified pid.
func getConnByPid(pid int) (*dtcp.PoolConn, error) {
	port := getPortByPid(pid)
	if port > 0 {
		if conn, err := dtcp.NewPoolConn(fmt.Sprintf("127.0.0.1:%d", port)); err == nil {
			return conn, nil
		} else {
			return nil, err
		}
	}
	return nil, werr.WithMsgF(`could not find port for pid "%d"`, pid)
}

// getPortByPid returns the listening port for specified pid.
// It returns 0 if no port found for the specified pid.
func getPortByPid(pid int) int {
	path := getCommFilePath(pid)
	if path == "" {
		return 0
	}
	return dconv.Int(dfile.GetContentsWithCache(path))
}

// getCommFilePath returns the pid to port mapping file path for given pid.
func getCommFilePath(pid int) string {
	path, err := getCommPidFolderPath()
	if err != nil {
		intlog.Errorf(context.TODO(), `%+v`, err)
		return ""
	}
	return dfile.Join(path, dconv.String(pid))
}

// getCommPidFolderPath retrieves and returns the available directory for storing pid mapping files.
func getCommPidFolderPath() (folderPath string, err error) {
	commPidFolderPathOnce.Do(func() {
		availablePaths := []string{
			"/var/tmp",
			"/var/run",
		}
		if path, _ := dfile.Home(".config"); path != "" {
			availablePaths = append(availablePaths, path)
		}
		availablePaths = append(availablePaths, dfile.Temp())
		for _, availablePath := range availablePaths {
			checkPath := dfile.Join(availablePath, defaultFolderNameForProcComm)
			if !dfile.Exists(checkPath) && dfile.Mkdir(checkPath) != nil {
				continue
			}
			if dfile.IsWritable(checkPath) {
				commPidFolderPath = checkPath
				break
			}
		}
		if commPidFolderPath == "" {
			err = werr.WithMsgF(
				`cannot find available folder for storing pid to port mapping files in paths: %+v`,
				availablePaths,
			)
		}
	})
	folderPath = commPidFolderPath
	return
}
