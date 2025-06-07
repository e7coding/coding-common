// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gspath implements file index and search for folders.
//
// It searches file internally with high performance in order by the directory adding sequence.
// Note that:
// If caching feature enabled, there would be a searching delay after adding/deleting files.
package jspath

import (
	"github.com/e7coding/coding-common/container/jarr"
	"github.com/e7coding/coding-common/container/jmap"
	"github.com/e7coding/coding-common/errs/jerr"
	"os"
	"sort"
	"strings"

	"github.com/e7coding/coding-common/internal/intlog"
	"github.com/e7coding/coding-common/os/jfile"
	"github.com/e7coding/coding-common/text/jstr"
)

// SPath manages the path searching feature.
type SPath struct {
	paths *jarr.SafeStrArr    // The searching directories array.
	cache *jmap.SafeStrStrMap // Searching cache map, it is not enabled if it's nil.
}

// SPathCacheItem is a cache item for searching.
type SPathCacheItem struct {
	path  string // Absolute path for file/dir.
	isDir bool   // Is directory or not.
}

var (
	// Path to searching object mapping, used for instance management.
	pathsMap = jmap.NewSafeStrAnyMap()
)

// New creates and returns a new path searching manager.
func New(path string, cache bool) *SPath {
	sp := &SPath{
		paths: jarr.NewSafeStrArr(),
	}
	if cache {
		sp.cache = jmap.NewSafeStrStrMap()
	}
	if len(path) > 0 {
		if _, err := sp.Add(path); err != nil {
			// intlog.Print(err)
		}
	}
	return sp
}

// Get creates and returns an instance of searching manager for given path.
// The parameter `cache` specifies whether using cache feature for this manager.
// If cache feature is enabled, it asynchronously and recursively scans the path
// and updates all sub files/folders to the cache using package gfsnotify.
func Get(root string, cache bool) *SPath {
	if root == "" {
		root = "/"
	}
	return pathsMap.GetOrPutFunc(root, func() interface{} {
		return New(root, cache)
	}).(*SPath)
}

// Search searches file `name` under path `root`.
// The parameter `root` should be an absolute path. It will not automatically
// convert `root` to absolute path for performance reason.
// The optional parameter `indexFiles` specifies the searching index files when the result is a directory.
// For example, if the result `filePath` is a directory, and `indexFiles` is [index.html, main.html], it will also
// search [index.html, main.html] under `filePath`. It returns the absolute file path if any of them found,
// or else it returns `filePath`.
func Search(root string, name string, indexFiles ...string) (filePath string, isDir bool) {
	return Get(root, false).Search(name, indexFiles...)
}

// SearchWithCache searches file `name` under path `root` with cache feature enabled.
// The parameter `root` should be an absolute path. It will not automatically
// convert `root` to absolute path for performance reason.
// The optional parameter `indexFiles` specifies the searching index files when the result is a directory.
// For example, if the result `filePath` is a directory, and `indexFiles` is [index.html, main.html], it will also
// search [index.html, main.html] under `filePath`. It returns the absolute file path if any of them found,
// or else it returns `filePath`.
func SearchWithCache(root string, name string, indexFiles ...string) (filePath string, isDir bool) {
	return Get(root, true).Search(name, indexFiles...)
}

// Set deletes all other searching directories and sets the searching directory for this manager.
func (sp *SPath) Set(path string) (realPath string, err error) {
	realPath = jfile.RealPath(path)
	if realPath == "" {
		realPath, _ = sp.Search(path)
		if realPath == "" {
			realPath = jfile.RealPath(jfile.Pwd() + jfile.Separator + path)
		}
	}
	if realPath == "" {
		return realPath, jerr.WithMsgF(`path "%s" does not exist`, path)
	}
	// The set path must be a directory.
	if jfile.IsDir(realPath) {
		realPath = strings.TrimRight(realPath, jfile.Separator)
		if sp.paths.IndexOf(realPath) != -1 {
			for _, v := range sp.paths.Slice() {
				sp.removeMonitorByPath(v)
			}
		}
		intlog.Print("paths clear:", sp.paths)
		sp.paths.Clear()
		if sp.cache != nil {
			sp.cache.Clear()
		}
		sp.paths.Append(realPath)
		sp.updateCacheByPath(realPath)
		sp.addMonitorByPath(realPath)
		return realPath, nil
	} else {
		return "", jerr.WithMsg(path + " should be a folder")
	}
}

// Add adds more searching directory to the manager.
// The manager will search file in added order.
func (sp *SPath) Add(path string) (realPath string, err error) {
	realPath = jfile.RealPath(path)
	if realPath == "" {
		realPath, _ = sp.Search(path)
		if realPath == "" {
			realPath = jfile.RealPath(jfile.Pwd() + jfile.Separator + path)
		}
	}
	if realPath == "" {
		return realPath, jerr.WithMsgF(`path "%s" does not exist`, path)
	}
	// The added path must be a directory.
	if jfile.IsDir(realPath) {
		// fmt.Println("gspath:", realPath, sp.paths.Search(realPath))
		// It will not add twice for the same directory.
		if sp.paths.IndexOf(realPath) < 0 {
			realPath = strings.TrimRight(realPath, jfile.Separator)
			sp.paths.Append(realPath)
			sp.updateCacheByPath(realPath)
			sp.addMonitorByPath(realPath)
		}
		return realPath, nil
	} else {
		return "", jerr.WithMsg(path + " should be a folder")
	}
}

// Search searches file `name` in the manager.
// The optional parameter `indexFiles` specifies the searching index files when the result is a directory.
// For example, if the result `filePath` is a directory, and `indexFiles` is [index.html, main.html], it will also
// search [index.html, main.html] under `filePath`. It returns the absolute file path if any of them found,
// or else it returns `filePath`.
func (sp *SPath) Search(name string, indexFiles ...string) (filePath string, isDir bool) {
	// No cache enabled.
	if sp.cache == nil {
		sp.paths.ByFunc(func(array []string) {
			path := ""
			for _, v := range array {
				path = jfile.Join(v, name)
				if stat, err := os.Stat(path); stat != nil && !os.IsNotExist(err) {
					path = jfile.Abs(path)
					// Security check: the result file path must be under the searching directory.
					if len(path) >= len(v) && path[:len(v)] == v {
						filePath = path
						isDir = stat.IsDir()
						break
					}
				}
			}
		})
		if len(indexFiles) > 0 && isDir {
			if name == "/" {
				name = ""
			}
			path := ""
			for _, file := range indexFiles {
				path = filePath + jfile.Separator + file
				if jfile.Exists(path) {
					filePath = path
					isDir = false
					break
				}
			}
		}
		return
	}
	// Using cache feature.
	name = sp.formatCacheName(name)
	if v := sp.cache.Get(name); v != "" {
		filePath, isDir = sp.parseCacheValue(v)
		if len(indexFiles) > 0 && isDir {
			if name == "/" {
				name = ""
			}
			for _, file := range indexFiles {
				if v = sp.cache.Get(name + "/" + file); v != "" {
					return sp.parseCacheValue(v)
				}
			}
		}
	}
	return
}

// Remove deletes the `path` from cache files of the manager.
// The parameter `path` can be either an absolute path or just a relative file name.
func (sp *SPath) Remove(path string) {
	if sp.cache == nil {
		return
	}
	if jfile.Exists(path) {
		for _, v := range sp.paths.Slice() {
			name := jstr.Replace(path, v, "")
			name = sp.formatCacheName(name)
			sp.cache.Delete(name)
		}
	} else {
		name := sp.formatCacheName(path)
		sp.cache.Delete(name)
	}
}

// Paths returns all searching directories.
func (sp *SPath) Paths() []string {
	return sp.paths.Slice()
}

// AllPaths returns all paths cached in the manager.
func (sp *SPath) AllPaths() []string {
	if sp.cache == nil {
		return nil
	}
	paths := sp.cache.Keys()
	if len(paths) > 0 {
		sort.Strings(paths)
	}
	return paths
}

// Size returns the count of the searching directories.
func (sp *SPath) Size() int {
	return sp.paths.Len()
}
