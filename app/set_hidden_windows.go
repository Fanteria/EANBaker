//go:build windows
// +build windows

package app

import (
	"syscall"
)

// setHidden sets the Windows hidden attribute
func setHidden(filename string) error {
	pointer, err := syscall.UTF16PtrFromString(filename)
	if err != nil {
		return err
	}
	return syscall.SetFileAttributes(pointer, syscall.FILE_ATTRIBUTE_HIDDEN)
}
