package ip

import (
	"github.com/zeromicro/go-zero/core/netx"
	"os"
	"strings"
)

const (
	allEths  = "0.0.0.0"
	envPodIP = "POD_IP"
)

/*func (util ipUtil) GetServerIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		// 排除回环接口和无效接口
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", errors.New("failed to determine server IP address")
}

func (util ipUtil) GetServiceAddress(port string) (*string, error) {
	ip, err := util.GetServerIP()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	address := ip + ":" + port
	return &address, nil
}*/

func FigureOutListenOn(listenOn string) string {
	fields := strings.Split(listenOn, ":")
	if len(fields) == 0 {
		return listenOn
	}

	host := fields[0]
	if len(host) > 0 && host != allEths {
		return listenOn
	}

	ip := GetInternalIp()
	if len(ip) == 0 {
		return listenOn
	}

	return strings.Join(append([]string{ip}, fields[1:]...), ":")
}

func GetInternalIp() string {
	ip := os.Getenv(envPodIP)
	if len(ip) == 0 {
		return netx.InternalIp()
	}
	return ip
}
