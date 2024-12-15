package utils

import "os"

func DeleteFileExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}
