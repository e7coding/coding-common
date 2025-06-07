// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package jfsnotify

import (
	"context"
	"github.com/e7coding/coding-common/container/jlist"
	"github.com/e7coding/coding-common/errs/jerr"

	"github.com/e7coding/coding-common/internal/intlog"
)

// Add monitors `path` with callback function `callbackFunc` to the watcher.
//
// The parameter `path` can be either a file or a directory path.
// The optional parameter `recursive` specifies whether monitoring the `path` recursively,
// which is true in default.
func (w *Watcher) Add(
	path string, callbackFunc func(event *Event), option ...WatchOption,
) (callback *Callback, err error) {
	return w.AddOnce("", path, callbackFunc, option...)
}

// AddOnce monitors `path` with callback function `callbackFunc` only once using unique name
// `name` to the watcher. If AddOnce is called multiple times with the same `name` parameter,
// `path` is only added to monitor once.
//
// It returns error if it's called twice with the same `name`.
//
// The parameter `path` can be either a file or a directory path.
// The optional parameter `recursive` specifies whether monitoring the `path` recursively,
// which is true in default.
func (w *Watcher) AddOnce(
	name, path string, callbackFunc func(event *Event), option ...WatchOption,
) (callback *Callback, err error) {
	var watchOption = w.getWatchOption(option...)
	w.nameSet.AddIfNotExistFunc(name, func() bool {
		// Firstly add the path to watcher.
		//
		// A path can only be watched once; watching it more than once is a no-op and will
		// not return an error.
		callback, err = w.addWithCallbackFunc(
			name, path, callbackFunc, option...,
		)
		if err != nil {
			return false
		}

		// If it's recursive adding, it then adds all sub-folders to the monitor.
		// NOTE:
		// 1. It only recursively adds **folders** to the monitor, NOT files,
		//    because if the folders are monitored and their sub-files are also monitored.
		// 2. It bounds no callbacks to the folders, because it will search the callbacks
		//    from its parent recursively if any event produced.
		if fileIsDir(path) && !watchOption.NoRecursive {
			for _, subPath := range fileAllDirs(path) {
				if fileIsDir(subPath) {
					if watchAddErr := w.watcher.Add(subPath); watchAddErr != nil {
						err = jerr.WithMsgErrF(
							err,
							`add watch failed for path "%s", err: %s`,
							subPath, watchAddErr.Error(),
						)
					} else {
						intlog.Printf(context.TODO(), "watcher adds monitor for: %s", subPath)
					}
				}
			}
		}
		if name == "" {
			return false
		}
		return true
	})
	return
}

func (w *Watcher) getWatchOption(option ...WatchOption) WatchOption {
	if len(option) > 0 {
		return option[0]
	}
	return WatchOption{}
}

// addWithCallbackFunc adds the path to underlying monitor, creates and returns a callback object.
// Very note that if it calls multiple times with the same `path`, the latest one will overwrite the previous one.
func (w *Watcher) addWithCallbackFunc(
	name, path string, callbackFunc func(event *Event), option ...WatchOption,
) (callback *Callback, err error) {
	var watchOption = w.getWatchOption(option...)
	// Check and convert the given path to absolute path.
	if realPath := fileRealPath(path); realPath == "" {
		return nil, jerr.WithMsgF(`"%s" does not exist`, path)
	} else {
		path = realPath
	}
	// Create callback object.
	callback = &Callback{
		Id:        callbackIdGenerator.Add(1),
		Func:      callbackFunc,
		Path:      path,
		name:      name,
		recursive: !watchOption.NoRecursive,
	}
	// Register the callback to watcher.
	w.callbacks.ByFunc(func(m map[string]interface{}) {
		lists := jlist.NewSafeList()
		if v, ok := m[path]; !ok {
			lists = jlist.NewSafeList()
			m[path] = lists
		} else {
			lists = v.(*jlist.SafeList)
		}
		callback.elem = lists.PushBack(callback)
	})
	// Add the path to underlying monitor.
	if err = w.watcher.Add(path); err != nil {
		err = jerr.WithMsgErrF(err, `add watch failed for path "%s"`, path)
	} else {
		intlog.Printf(context.TODO(), "watcher adds monitor for: %s", path)
	}
	// Add the callback to global callback map.
	callbackIdMap.Put(callback.Id, callback)
	return
}

// Close closes the watcher.
func (w *Watcher) Close() {
	close(w.closeChan)
	if err := w.watcher.Close(); err != nil {
		intlog.Errorf(context.TODO(), `%+v`, err)
	}
	w.events.Close()
}

// Remove removes watching and all callbacks associated with the `path` recursively.
// Note that, it's recursive in default if given `path` is a directory.
func (w *Watcher) Remove(path string) error {
	var (
		err          error
		subPaths     []string
		removedPaths = make([]string, 0)
	)
	removedPaths = append(removedPaths, path)
	if fileIsDir(path) {
		subPaths, err = fileScanDir(path, "*", true)
		if err != nil {
			return err
		}
		removedPaths = append(removedPaths, subPaths...)
	}

	for _, removedPath := range removedPaths {
		// remove the callbacks of the path.
		if value := w.callbacks.Del(removedPath); value != nil {
			list := value.(*jlist.SafeList)
			for {
				if item, _ := list.PopFront(); item != nil {
					callbackIdMap.Delete(item.(*Callback).Id)
				} else {
					break
				}
			}
		}
		// remove the monitor of the path from underlying monitor.
		if watcherRemoveErr := w.watcher.Remove(removedPath); watcherRemoveErr != nil {
			err = jerr.WithMsgErrF(
				err,
				`remove watch failed for path "%s", err: %s`,
				removedPath, watcherRemoveErr.Error(),
			)
		}
	}
	return err
}

// RemoveCallback removes callback with given callback id from watcher.
//
// Note that, it auto removes the path watching if there's no callback bound on it.
func (w *Watcher) RemoveCallback(callbackId int) {
	callback := (*Callback)(nil)
	if r := callbackIdMap.Get(callbackId); r != nil {
		callback = r.(*Callback)
	}
	if callback != nil {
		if r := w.callbacks.Get(callback.Path); r != nil {
			r.(*jlist.SafeList).Remove(callback.elem)
		}
		callbackIdMap.Delete(callbackId)
		if callback.name != "" {
			w.nameSet.Remove(callback.name)
		}
	}
}
