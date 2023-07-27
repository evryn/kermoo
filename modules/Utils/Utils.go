package Utils

import (
	"fmt"
	"net"
	"os"
	"time"

	"math/rand"

	"github.com/spf13/cobra"
)

func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func GetDurationFlag(cmd *cobra.Command, flag string) (*time.Duration, error) {
	value, _ := cmd.Flags().GetString(flag)
	duration, err := time.ParseDuration(value)

	if err != nil {
		return nil, fmt.Errorf("'%s' is not a valid duration. Valid duration examples: 200ms, 5s, 10m, 2h, 1d", value)
	}

	return &duration, nil
}

func GetIpList() []string {
	list := []string{}

	ifaces, err := net.Interfaces()
	if err != nil {
		return list
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // Is down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // Is a lookpack interface
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}

			if ip.To4() != nil {
				list = append(list, ip.To4().String())
			} else if ip.To16() != nil {
				list = append(list, ip.To16().String())
			}
		}
	}
	return list
}

func GenerateRandomFloat32Between(min float32, max float32) float32 {
	return min + rand.Float32()*(max-min)
}
