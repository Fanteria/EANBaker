//go:build !windows
// +build !windows

package app

func setHidden(string) error {
	return nil
}
