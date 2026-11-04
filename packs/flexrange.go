package portscan

import (
	"fmt"
)

func Range(st, en uint16) []uint16 {
	//invalid ranges
	if st > en || st < 1 || en > 65535 {
		fmt.Printf("ERROR: Range(%d, %d) is invalid.\n", st, en)
		return nil
	}
	ports := make([]uint16, 0, en - st + 1)
	for i := st; i <= en; i++ {
		ports = append(ports, i)
	}
	return ports
}

func FlexRange(args ...any) []uint16 {
	if len(args) == 0 {
		return nil
	}
	portsMap := make(map[uint16]struct{})
	isValid := func(port uint16) bool {
		return port >= 1 && port <= 65535
	}
	
	for i, arg := range args {
		switch val := arg.(type) {
		case uint16:
			if isValid(val) {
				portsMap[val] = struct{}{}
			} else {
				fmt.Printf("ERROR: Invalid port %d at the argument %d\n", val, i)
				return nil
			}
		case []uint16:
			if len(val) == 0 {
				fmt.Printf("ERROR: Invalid argument %d. (empty or nil)\n", i)
				return nil	
			}
			for j, port := range val {
				if isValid(port) {
					portsMap[port] = struct{}{}
				} else {
					fmt.Printf("ERROR: Invalid port %d at the range %d, argument %d\n", val, i, j)
					return nil
				}
			}
		default:
			fmt.Printf("ERROR: Invalid argument %d. type: %T.\n", i, arg)
			return nil
		}
	}
	ports := make([]uint16, 0, len(portsMap))
	for port := range portsMap {
		ports = append(ports, port)
	}
	return ports
}