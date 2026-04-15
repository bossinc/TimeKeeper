//go:build windows

package main

import (
	"syscall"
	"unsafe"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	getForeground    = user32.NewProc("GetForegroundWindow")
	getWindowText    = user32.NewProc("GetWindowTextW")
	getLastInputInfo = user32.NewProc("GetLastInputInfo")
	getTickCount     = kernel32.NewProc("GetTickCount")
)

type lastInputInfo struct {
	cbSize uint32
	dwTime uint32
}

func getActiveWindowTitle() string {
	hwnd, _, _ := getForeground.Call()
	if hwnd == 0 {
		return ""
	}
	buf := make([]uint16, 512)
	getWindowText.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return syscall.UTF16ToString(buf)
}

func getIdleTimeMs() int64 {
	var info lastInputInfo
	info.cbSize = uint32(unsafe.Sizeof(info))
	getLastInputInfo.Call(uintptr(unsafe.Pointer(&info)))
	tick, _, _ := getTickCount.Call()
	return int64(uint32(tick) - info.dwTime)
}
