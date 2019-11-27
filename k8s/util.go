package k8s

import (
	"net"
)

func (pod *Pod)Ping() error {
	if pod.Addr == inmemPodAddr {
		return nil
	}

	c, err := net.Dial("tcp", pod.Addr + ":" + pod.Port)
	if err == nil {
		c.Close()
	}
	return err
}
