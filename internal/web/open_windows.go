//go:build windows && cgo

package web

import (
	"github.com/webview/webview_go"
)

func OpenDashboard() {
	url := "http://localhost:8080"
	
	w := webview.New(true)
	if w == nil {
		return
	}
	defer w.Destroy()
	w.SetTitle("Rune V2 | Management Console")
	w.SetSize(1024, 768, webview.HintNone)
	w.Navigate(url)
	w.Run()
}