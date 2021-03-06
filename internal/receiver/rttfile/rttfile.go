// Copyright 2020 Thomas.Hoehenleitner [at] seerose.net
// Use of this source code is governed by a license that can be found in the LICENSE file.

// Package rttfile reads trice data from file.
package rttfile

import (
	"io"
	"os"

	"github.com/rokath/trice/internal/global"
)

// RTTFile delivers
type RTTFile struct {
	fh *os.File
}

// New creates an instance of simulator.
func New() *RTTFile {
	r := &RTTFile{}
	return r
}

// Read is the exported method.
func (p *RTTFile) Read(b []byte) (int, error) {
	n, err := p.fh.Read(b)
	if io.EOF == err {
		global.OsExit(1)
	}
	return n, err
}

// Close is part of the exported interface io.ReadCloser. It ends the connection.
func (p *RTTFile) Close() error {
	return p.fh.Close()
}

// Open initializes
func (p *RTTFile) Open(fn string) error {
	var err error
	p.fh, err = os.Open(fn)
	return err
}
