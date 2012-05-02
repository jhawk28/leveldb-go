// Copyright 2012 The LevelDB-Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"io"
	"os"
)

// File is a readable, writable sequence of bytes.
//
// Typically, it will be an *os.File, but test code may choose to substitute
// memory-backed implementations.
type File interface {
	io.Closer
	io.ReaderAt
	io.Writer
	Stat() (stat os.FileInfo, err error)
}

// FileSystem is a namespace for files.
//
// The names are filepath names: they may be / separated or \ separated,
// depending on the underlying operating system.
type FileSystem interface {
	Create(name string) (File, error)
	Open(name string) (File, error)
	Remove(name string) error

	// Lock locks the given file. A nil Closer is returned if an error occurred.
	// Otherwise, close that Closer to release the lock.
	//
	// On Linux, a lock has the same semantics as fcntl(2)'s advisory locks.
	// In particular, closing any other file descriptor for the same file will
	// release the lock prematurely.
	//
	// Lock is not yet implemented on other operating systems, and calling it
	// will return an error.
	Lock(name string) (io.Closer, error)

	// List returns a listing of the given directory. The names returned are
	// relative to dir.
	List(dir string) ([]string, error)
}

// DefaultFileSystem is a FileSystem implementation backed by the underlying
// operating system's file system.
var DefaultFileSystem FileSystem = defFS{}

type defFS struct{}

func (defFS) Create(name string) (File, error) {
	return os.Create(name)
}

func (defFS) Open(name string) (File, error) {
	return os.Open(name)
}

func (defFS) Remove(name string) error {
	return os.Remove(name)
}

func (defFS) List(dir string) ([]string, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.Readdirnames(-1)
}
