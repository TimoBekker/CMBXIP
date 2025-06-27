package main

import (
	"nikeron/cmbxip/config"
	"syscall"
)

func init() {
	if !config.LeaveWindowsConsole() {
		modkernel32 := syscall.NewLazyDLL("kernel32.dll")
		moduser32 := syscall.NewLazyDLL("user32.dll")

		hwnd, _, _ := modkernel32.NewProc("GetConsoleWindow").Call()
		arg := uint32(0)
		moduser32.NewProc("ShowWindow").Call(uintptr(hwnd), uintptr(arg))
	}
}
