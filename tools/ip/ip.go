package ip

import "net"

func Available() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	return doPacking(addrs), nil
}

func ipv4(ip net.IP) bool {
	ip4 := ip.To4()
	return ip4 != nil
}

func ipv6(ip net.IP) bool {
	if ip4 := ip.To4(); ip4 != nil {
		return false
	} else if len(ip) == net.IPv6len {
		return true
	}
	return false
}

func doPacking(addrs []net.Addr) (arr []string) {
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			if ipNet.IP.IsLoopback() {
				continue
			}
			if ipv4(ipNet.IP) {
				ipNet.IP = ipNet.IP.To4()
				if ipNet.IP[0] == 169 && ipNet.IP[1] == 254 {
					continue
				}
			} else if ipv6(ipNet.IP) {
				ipNet.IP = ipNet.IP.To16()
				if ipNet.IP[0] == 0xFe && ipNet.IP[1] == 0x80 {
					continue
				}
			} else {
				continue
			}
			arr = append(arr, ipNet.IP.String())
		}
	}
	return
}
