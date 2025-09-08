//go:build !windows
// +build !windows

package app

// Is a no-op function for non-Windows platforms.
// On Unix-like systems, files starting with '.' are automatically treated as hidden.
func setHidden(string) error {
	return nil
}
