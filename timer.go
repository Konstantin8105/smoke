package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/Konstantin8105/ds"
	"github.com/Konstantin8105/glsymbol"
	"github.com/Konstantin8105/ms/window"
	"github.com/Konstantin8105/vl"
)

func init() {
	vl.SpecificSymbol(false)
}

func main() {
	ch := make(chan func() (fus bool), 1000)

	var design vl.Widget
	{
		var list vl.List
		list.Compress()

		var status vl.Text

		var limit int = 5400 // 1.5 hours
		var timer vl.Text
		start := time.Now()
		go func() {
			for {
				ch <- func() (fus bool) {
					diff := time.Now().Sub(start)
					sec := int(diff.Seconds())
					if sec < limit {
						status.SetText("")
					} else {
						status.SetText("TIME IS OVER")
					}
					timer.SetText(fmt.Sprintf("%06d", sec))
					return true
				}
				time.Sleep(time.Millisecond * 500)
			}
		}()
		list.Add(&timer)

		list.Add(&status)

		var reset vl.Button
		reset.SetText("Reset")
		reset.OnClick = func() {
			start = time.Now()
		}
		list.Add(&reset)

		step := int64(300) // seconds
		var add vl.Button
		add.SetText(fmt.Sprintf("Add %d sec", step))
		add.OnClick = func() {
			dur := time.Duration(-step * 1e9)
			start = start.Add(dur)
		}
		list.Add(&add)

		design = &list
	}

	var doc vl.Widget
	doc = vl.TextStatic(`Documentation:
1) Reset timer to zero.
2) Add on step seconds.
3) After 5400 seconds indicated`)

	var (
		left  = window.NewTui(design)
		right = window.NewTui(doc)

		ws = [2]ds.Window{left, right}
	)

	screen, err := ds.New("Timer", ws, &ch)
	if err != nil {
		panic(err)
	}

	// add fonts
	f, err := DefaultFont()
	if err != nil {
		return
	}
	left.SetFont(f)
	right.SetFont(f)

	quit := make(chan struct{})

	screen.Run(&quit)
}

// DefaultFont return default font
func DefaultFont() (_ *glsymbol.Font, err error) {
	var (
		low   = rune(byte(32))
		high  = rune(byte(127))
		scale = int32(16) // font size
	)
	return glsymbol.LoadTruetype(
		strings.NewReader(glsymbol.DefaultEmbeddedFont),
		scale,
		rune(byte(low)),
		rune(byte(high)),
	)
}
