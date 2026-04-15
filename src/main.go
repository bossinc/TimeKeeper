package main

import "github.com/gotk3/gotk3/gtk"

func main() {
	gtk.Init(nil)
	w := newMainWindow()
	w.ShowAll()
	gtk.Main()
}
