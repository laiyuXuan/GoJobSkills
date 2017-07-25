package log

import (
	"log"
	"bytes"
	"os"
)

func GetLogger() *log.Logger{
	var buffer bytes.Buffer
	buffer.WriteString("[")
	buffer.WriteString("LOG")
	buffer.WriteString("] --->> ")
	return log.New(os.Stdout, buffer.String(), log.LstdFlags | log.Llongfile)
}
