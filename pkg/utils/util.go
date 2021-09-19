package redisflu

import "os"

func GetOrCreateFileToWrite(filename string) (*os.File, error) {
	var file *os.File
	var err error

	file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(filename)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return file, err
}
