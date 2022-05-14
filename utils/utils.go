package utils

import (
	"net"
	"os"
	"reflect"
	"sync"
)

const (
	Empty = ""
)

var (
	internalIPOnce sync.Once
	internalIP     = ""
)

func IsNotNil(Object interface{}) bool {
	return !IsNilObject(Object)
}

func IsNilObject(object interface{}) bool {
	if object == nil {
		return true
	}
	value := reflect.ValueOf(object)
	kind := value.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return true
	}

	return false
}

//GetInternal 获取内部本机ip
func GetInternal() string {
	internalIPOnce.Do(func() {
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			os.Stderr.WriteString("Oops:" + err.Error())
			os.Exit(1)
		}
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					internalIP = ipnet.IP.To4().String()
				}
			}
		}
	})
	return internalIP
}
