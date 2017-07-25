package utils

import (
	"os"
)

func Save2File(filePath string, content string) error  {
	file, err := os.OpenFile(filePath, os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0666)
	defer file.Close()

	if err != nil{
		return err
	}
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}
