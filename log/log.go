package log

import (
	"log"
	"runtime"
	"bytes"
	"os"
)

func GetLogger() *log.Logger{
	_, file, _, _ := runtime.Caller(0);
	var buffer bytes.Buffer
	buffer.WriteString("[")
	buffer.WriteString(file)
	buffer.WriteString("] ->>>")
	return log.New(os.Stdout, buffer.String(), log.LstdFlags)
}
