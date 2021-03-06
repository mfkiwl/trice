// Copyright 2020 Thomas.Hoehenleitner [at] seerose.net
// Use of this source code is governed by a license that can be found in the LICENSE file.

// Package randomdummy_test is a blackbox test
package randomdummy_test

import (
	"errors"
	"os/exec"
	"testing"

	"github.com/rokath/trice/internal/randomdummy"
	"github.com/rokath/trice/pkg/lib"
)

// TestZeroMode checks Read mode 0
func TestZeroMode(t *testing.T) {
	s := randomdummy.New(randomdummy.ZeroRandomSeed, "ChaosMode", randomdummy.NoDelay, randomdummy.NoLimit)
	act := make([]byte, 19)
	s.Read(act)
	exp := []byte{0x9f, 0x90, 0xa3, 0xe2, 0x50, 0xc4, 0x9d, 0xa7, 0xc6, 0x30, 0x5d, 0xc5, 0x77, 0xe4, 0x93, 0x33, 0x29, 0x50, 0xdc}
	lib.Equals(t, exp, act)
}

// TestZeroModeLimit checks Read mode 0 with limit
func TestZeroModeLimit(t *testing.T) {
	s := randomdummy.New(randomdummy.ZeroRandomSeed, "ChaosMode", randomdummy.NoDelay, 9)
	act := make([]byte, 19)
	n, err := s.Read(act)
	exp := []byte{0x9f, 0x90, 0xa3, 0xe2, 0x50, 0xc4, 0x9d, 0xa7, 0xc6}
	act = act[:n]
	lib.Equals(t, exp, act)
	lib.Equals(t, errors.New("EOF"), err)
}

// TestUnknownMode checks Read mode unknown
func TestUnknownMode(t *testing.T) {
	s := randomdummy.New(randomdummy.ZeroRandomSeed, "unknownMode", randomdummy.NoDelay, randomdummy.NoLimit)
	act := make([]byte, 19)
	s.Read(act)
	exp := []byte{0x9f, 0x90, 0xa3, 0xe2, 0x50, 0xc4, 0x9d, 0xa7, 0xc6, 0x30, 0x5d, 0xc5, 0x77, 0xe4, 0x93, 0x33, 0x29, 0x50, 0xdc}
	lib.Equals(t, exp, act)
}

// TestBareModeShort checks Read mode bare
func TestBareModeTooShort(t *testing.T) {
	s := randomdummy.New(randomdummy.ZeroRandomSeed, "BareModeNoSync", randomdummy.NoDelay, randomdummy.NoLimit)
	b := make([]byte, 3)
	act, err := s.Read(b)
	lib.Equals(t, 0, act)
	lib.Equals(t, errors.New("too short for bare"), err)
}

// TestWrapModeTooShort checks Read mode wrap
func TestWrapModeTooShort(t *testing.T) {
	s := randomdummy.New(randomdummy.ZeroRandomSeed, "WrapModeWithValidCrc", randomdummy.NoDelay, randomdummy.NoLimit)
	b := make([]byte, 3)
	act, err := s.Read(b)
	lib.Equals(t, 0, act)
	lib.Equals(t, errors.New("too short for wrap"), err)
}

// TestBareMode checks Read mode bare
func TestBareMode(t *testing.T) {
	s := randomdummy.New(randomdummy.ZeroRandomSeed, "BareModeNoSync", randomdummy.NoDelay, randomdummy.NoLimit)
	act := make([]byte, 19)
	n, _ := s.Read(act)
	act = act[:n]
	exp := []byte{0x9f, 0x90, 0xa3, 0xe2}
	lib.Equals(t, exp, act)

	s.Read(act)
	exp = []byte{0x50, 0xc4, 0x9d, 0xa7}
	lib.Equals(t, exp, act)
}

// TestWrapMode checks Read mode wrap
func TestWrapMode(t *testing.T) {
	// create instance
	s := randomdummy.New(randomdummy.ZeroRandomSeed, "WrapModeWithValidCrc", randomdummy.NoDelay, randomdummy.NoLimit)
	act := make([]byte, 19) // receive space
	n, _ := s.Read(act)
	act = act[:n] // limit to read bytes
	exp := []byte{0xeb, 0x60, 0x60, 0x45, 0x50, 0xc4, 0x9d, 0xa7}
	lib.Equals(t, exp, act)

	s.Read(act)
	exp = []byte{0xeb, 0x60, 0x60, 0xd8, 0x77, 0xe4, 0x93, 0x33}
	lib.Equals(t, exp, act)
}

// TestRNDchaos expects installed trice compiled from actual sources.
//
// This kind of test does not work just with HandleArgs function, because of os.Exit(0) on io.EOF in -source RND.
// Endless waiting there does also not work, so this apprpoach is just a quick solution.
func TestRNDchaos(t *testing.T) {
	cmd := exec.Command("trice", "log", "-idlist", "c:/repos/trice/test/til.json", "-source", "RND", "-rndLimit", "10", "-rndMode", "ChaosMode", "-lg", "off", "-ts", "off", "-color", "off")
	out, err := cmd.CombinedOutput()
	if nil != err {
		t.Fail()
	}
	act := string(out)
	exp := `RND: trice:discarding byte 0x9f (dez 159, char ' ')
RND: trice:discarding byte 0x90 (dez 144, char ' ')
RND: trice:discarding byte 0xa3 (dez 163, char ' ')
`
	lib.Equals(t, exp, act)
}
