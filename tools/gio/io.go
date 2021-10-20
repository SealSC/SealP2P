package gio

import "os"

func FileExist(f string) bool {
	_, err := os.Stat(f)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}
