// Copyright 2020 Thomas.Hoehenleitner [at] seerose.net
// Use of this source code is governed by a license that can be found in the LICENSE file.

// blackbox test
package disp_test

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/rokath/trice/internal/disp"
	"github.com/rokath/trice/pkg/cage"
	"github.com/rokath/trice/pkg/lib"
)

func lineGenerator(t *testing.T, s string, len, count int, wg *sync.WaitGroup) {
	var ss []string
	var i int
	// create line slice
	for i = 0; i < len; i++ { // line length
		var c string
		switch i & 11 { // 11 color channels
		case 0:
			c = "err:"
		case 1:
			c = "wrn:"
		case 2:
			c = "rd_:"
		case 3:
			c = "wr_:"
		case 4:
			c = "tim:"
		case 5:
			c = "att:"
		case 6:
			c = "dbg:"
		case 7:
			c = "dia:"
		case 8:
			c = "isr:"
		case 9:
			c = "sig:"
		case 10:
			c = "tst:"
		}
		su := fmt.Sprintf(c+s+"%02x ", i) // line position
		ss = append(ss, su)
	}
	ss = append(ss, "\n") // line end

	wg.Add(1)
	i = 0
	go func() {
		defer wg.Done()
		for ; i < count; i++ { // line count
			lib.Ok(t, disp.Out(ss))
		}
	}()
}

// TestServerMutex checks if several lines are displayed coreectly line by line
func TestServerMutex(t *testing.T) {
	disp.IPPort = "61498"
	cage.Name = "./testdata/serverMutexTest.log"
	uniqName := "./testdata/serverMutexUniq.txt"
	os.Remove(cage.Name)
	os.Remove(uniqName)

	var wg sync.WaitGroup

	var exe string
	if "windows" == runtime.GOOS {
		exe = "C:\\Users\\ms\\go\\bin\\trice.exe" // how to solve this universal?
	} else {
		t.Fail() // todo
	}
	disp.StartServer(exe)
	err := disp.Connect()
	disp.Out = disp.RemoteOut // re-direct output
	lib.Ok(t, err)

	var result int64
	err = disp.PtrRPC.Call("Server.ColorPalette", []string{"off"}, &result)
	lib.Ok(t, err)

	ss := []string{"first line", "\n"}
	err = disp.Out(ss)
	lib.Ok(t, err)

	ll := 21
	lc := 10
	lv := 9
	for i := 0; i < lv; i++ {
		li := fmt.Sprintf("line%d", i)
		lineGenerator(t, li, ll, lc, &wg)
	}

	wg.Wait()
	lib.Ok(t, disp.ScShutdownRemoteDisplayServer(0))
	n := lib.UniqLines(t, cage.Name, uniqName)
	lib.Equals(t, n, lv+4) // first line + 9 lines + server line
	os.Remove(uniqName)
	time.Sleep(200 * time.Millisecond) // Why to wait here? 50ms is not enough.
	os.Remove(cage.Name)
}
