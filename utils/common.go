package utils

import (
	"net"
)

func IsTimeOut(err error) bool {
	if err, ok := err.(net.Error); ok && err.Timeout(){
		return true
	}
	return false
}
