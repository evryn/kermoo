package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"math/rand"

	"gopkg.in/yaml.v3"
)

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

func RandomFloat(min float32, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func RandomDuration(min, max time.Duration) (*time.Duration, error) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	diff := int64(max - min)

	if diff <= 0 {
		return nil, fmt.Errorf("duration is invalid since the range is zero or negative")
	}

	dur := min + time.Duration(r.Int63n(diff))

	return &dur, nil
}

func YamlToJSON(yamlStr string) (string, error) {
	m := make(map[interface{}]interface{})

	// Unmarshal the YAML into a map
	err := yaml.Unmarshal([]byte(yamlStr), &m)
	if err != nil {
		return "", err
	}

	// Convert map keys to strings since JSON only allows string keys
	strMap := convertMapKeysToString(m)

	// Marshal the map into a JSON string
	jsonBytes, err := json.Marshal(strMap)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func convertMapKeysToString(m map[interface{}]interface{}) map[string]interface{} {
	n := make(map[string]interface{})
	for k, v := range m {
		switch child := v.(type) {
		case map[interface{}]interface{}:
			n[fmt.Sprint(k)] = convertMapKeysToString(child)
		default:
			n[fmt.Sprint(k)] = v
		}
	}
	return n
}

func GetDuplicates(items []string) []string {
	itemCount := make(map[string]int)
	duplicates := []string{}

	for _, item := range items {
		itemCount[item]++
	}

	for item, count := range itemCount {
		if count > 1 {
			duplicates = append(duplicates, item)
		}
	}

	return duplicates
}
