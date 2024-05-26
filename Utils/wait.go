package Utils

import (
	"fmt"
	"net"
	"time"
)

func WaitForPort(host string, port string, timeout time.Duration) error {
	start := time.Now()
	for {
		_, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", host, port), time.Second)
		if err == nil {
			fmt.Printf("Port %s is available\n", port)
			return nil
		}
		if time.Since(start) > timeout {
			return fmt.Errorf("timeout: port %s is not available", port)
		}
		time.Sleep(time.Second) // Wait for 1 second before trying again
	}
}
