// -build ansi
// Copyright 2020 Thomas.Hoehenleitner [at] seerose.net
// Use of this source code is governed by a license that can be found in the LICENSE file.

package disp

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/mgutz/ansi"
)

var (
	// Out is an exported function pointer, which can be redirected for example to a client call
	Out = out

	// mux is for syncing line output
	mux sync.Mutex

	colorizeERROR     = ansi.ColorFunc("11:red")
	colorizeWARNING   = ansi.ColorFunc("11+i:red")
	colorizeMESSAGE   = ansi.ColorFunc("green+h:black")
	colorizeREAD      = ansi.ColorFunc("101:black")
	colorizeWRITE     = ansi.ColorFunc("101+i:black")
	colorizeTIME      = ansi.ColorFunc("108:blue")
	colorizeATTENTION = ansi.ColorFunc("11:green")
	colorizeDEBUG     = ansi.ColorFunc("130+i")
	colorizeDIAG      = ansi.ColorFunc("161+B")
	colorizeINTERRUPT = ansi.ColorFunc("13+i")
	colorizeSIGNAL    = ansi.ColorFunc("118+i")
	colorizeTEST      = ansi.ColorFunc("yellow+h:black")
	colorizeINFO      = ansi.ColorFunc("121+i")
)

// out displays ss and sets color.
// ss is a slice containing substring parts of one line.
// Each substring inside ss is result of Trice or contains prefix,
// timestamp or postfix and can have its own color prefix.
// The last substring inside the slice is definitly the last in
// the line and has already trimmed its newline.
// Linebreaks inside the substrings are not handled separately (yet).
func out(ss []string) error {
	var line string

	mux.Lock()
	for _, s := range ss {
		c := colorize(s)
		line += fmt.Sprint(c)
	}

	log.SetFlags(0)
	log.Print(line)

	mux.Unlock()
	return nil
}

// colorize prefixes s with a ansi color code according to this conditions
// if "\n" == s add ANSI reset code
// if COL: is begin of string add ANSI color code according to COL:
// if col: is begin of string replace col: with ANSI color code according to col:
func colorize(s string) string {
	if "off" == ColorPalette || "none" == ColorPalette {
		return s
	}
	if "\n" == s {
		return ansi.Reset + s
	}
	sc := strings.SplitN(s, ":", 2)
	if len(sc) < 2 {
		/*
		   Avoid case 1 == len(sc):
		   panic: runtime error: index out of range [1] with length 1

		   goroutine 123 [running]:
		   github.com/rokath/trice/pkg/disp.colorize(0xb521b3, 0x1, 0x1f, 0xb52160)
		           C:/GitRepos/trice/pkg/disp/dispAnsiOut.go:126 +0x1b90
		   github.com/rokath/trice/pkg/disp.out(0xc0004d4420, 0x15, 0x15, 0xc000477bb0, 0x413ad3)
		           C:/GitRepos/trice/pkg/disp/dispAnsiOut.go:50 +0xa1
		   github.com/rokath/trice/pkg/disp.(*Server).Out(0xbac7b8, 0xc0004d4420, 0x15, 0x15, 0xc0004947a0, 0x0, 0x0)
		           C:/GitRepos/trice/pkg/disp/dispServer.go:43 +0x55
		   reflect.Value.call(0xc00048a480, 0xc0000b2448, 0x13, 0x85f5cd, 0x4, 0xc000477f18, 0x3, 0x3, 0xc000477e88, 0xc00048fda0, ...)
		           c:/go/src/reflect/value.go:460 +0x5fd
		   reflect.Value.Call(0xc00048a480, 0xc0000b2448, 0x13, 0xc000477f18, 0x3, 0x3, 0xc000488450, 0x0, 0x0)
		           c:/go/src/reflect/value.go:321 +0xbb
		   net/rpc.(*service).call(0xc00049ca00, 0xc000126000, 0xc000088020, 0xc000088030, 0xc000120180, 0xc00007e580, 0x7c6c40, 0xc0004d2900, 0x197, 0x7ba9c0, ...)
		           c:/go/src/net/rpc/server.go:377 +0x176
		   created by net/rpc.(*Server).ServeCodec
		           c:/go/src/net/rpc/server.go:474 +0x432
		*/
		return s
	}
	switch sc[0] {
	case "e", "err", "error":
		s = sc[1] // remove channel info
		fallthrough
	case "E", "ERR", "ERROR":
		return colorizeERROR(s)
	case "w", "wrn", "warning":
		s = sc[1] // remove channel info
		fallthrough
	case "W", "WRN", "WARNING":
		return colorizeWARNING(s)
	case "m", "msg", "message":
		s = sc[1] // remove channel info
		fallthrough
	case "M", "MSG", "MESSAGE":
		return colorizeMESSAGE(s)
	case "rd", "rd_":
		s = sc[1] // remove channel info
		fallthrough
	case "RD", "RD_":
		return colorizeREAD(s)
	case "wr", "wr_":
		s = sc[1] // remove channel info
		fallthrough
	case "WR", "WR_":
		return colorizeWRITE(s)
	case "tim", "time":
		s = sc[1] // remove channel info
		fallthrough
	case "TIM", "TIME":
		return colorizeTIME(s)
	case "att", "attention":
		s = sc[1] // remove channel info
		fallthrough
	case "ATT", "ATTENTION":
		return colorizeATTENTION(s)
	case "d", "db", "dbg", "debug":
		s = sc[1] // remove channel info
		fallthrough
	case "D", "DB", "DBG", "DEBUG":
		return colorizeDEBUG(s)
	case "dia", "diag":
		s = sc[1] // remove channel info
		fallthrough
	case "DIA", "DIAG":
		return colorizeDIAG(s)
	case "isr", "interrupt":
		s = sc[1] // remove channel info
		fallthrough
	case "ISR", "INTERRUPT":
		return colorizeDIAG(s)
	case "s", "sig", "signal":
		s = sc[1] // remove channel info
		fallthrough
	case "S", "SIG", "SIGNAL":
		return colorizeSIGNAL(s)
	case "t", "tst", "test":
		s = sc[1] // remove channel info
		fallthrough
	case "T", "TST", "TEST":
		return colorizeTEST(s)
	case "i", "inf", "info", "informal":
		s = sc[1] // remove channel info
		fallthrough
	case "I", "INF", "INFO", "INFORMAL":
		return colorizeINFO(s)
	}
	return s
}