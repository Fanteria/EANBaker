//go:build windows
// +build windows

package app

import (
	"syscall"
)

// Sets the Windows hidden attribute on the specified file.
// Uses Windows API to mark the file as hidden in the file-system.
func setHidden(filename string) error {
	pointer, err := syscall.UTF16PtrFromString(filename)
	if err != nil {
		return err
	}
	return syscall.SetFileAttributes(pointer, syscall.FILE_ATTRIBUTE_HIDDEN)
}
