// Copy of https://github.com/rivo/tview/blob/master/demos/textview/main.go
package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func tviewTextView() {
	app := tview.NewApplication()
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	numSelections := 0
	go func() {
		writer := tview.ANSIWriter(textView)

		content, err := os.ReadFile("diff.patch")
		if err != nil {
			fmt.Fprintf(writer, "Couldn't load file, err: %s", err)
			return
		}

		var out bytes.Buffer
		quick.Highlight(&out, string(content), "diff", "terminal16m", "monokai")
		prefix := fmt.Sprintf("This is tview's text area.\n\n")
		fmt.Fprintf(writer, "%s%s", prefix, out.String())
	}()
	textView.SetDoneFunc(func(key tcell.Key) {
		currentSelection := textView.GetHighlights()
		if key == tcell.KeyEnter {
			if len(currentSelection) > 0 {
				textView.Highlight()
			} else {
				textView.Highlight("0").ScrollToHighlight()
			}
		} else if len(currentSelection) > 0 {
			index, _ := strconv.Atoi(currentSelection[0])
			if key == tcell.KeyTab {
				index = (index + 1) % numSelections
			} else if key == tcell.KeyBacktab {
				index = (index - 1 + numSelections) % numSelections
			} else {
				return
			}
			textView.Highlight(strconv.Itoa(index)).ScrollToHighlight()
		}
	})
	textView.SetBorder(false)
	if err := app.SetRoot(textView, true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
}
