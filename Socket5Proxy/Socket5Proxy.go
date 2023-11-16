package Socket5Proxy

import (
	"io"
	"net"
	"strconv"
)

var (
	PackageWithAuth       = []byte{0x05, 0x02}
	PackageAuthSuccess    = []byte{0x05, 0x00}
	PackageAuthFailed     = []byte{0x05, 0x01}
	PackageConnectSuccess = []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
)

const (
	User     = "PaladinDu"
	PassWord = "286e71940aae4331bc6e4991e31fa8fc"
)

func Socket5Proxy(connect net.Conn) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()

	if connect == nil {
		return
	}
	defer func() {
		_ = connect.Close()
	}()
	b := make([]byte, 1024)

	n, err := connect.Read(b)
	if err != nil {
		return
	}

	if b[0] == 0x05 {
		_, _ = connect.Write(PackageWithAuth)

		n, err = connect.Read(b)
		if err != nil {
			return
		}

		userLength := int(b[1])
		user := string(b[2:(2 + userLength)])
		pass := string(b[(2 + userLength):])

		if User == user && PassWord == pass {
			_, _ = connect.Write(PackageAuthSuccess)
		} else {
			_, _ = connect.Write(PackageAuthFailed)
			return
		}

		n, err = connect.Read(b)
		var host string
		switch b[3] {
		case 0x01:
			host = net.IPv4(b[4], b[5], b[6], b[7]).String()
		case 0x03:
			host = string(b[5 : n-2])
		case 0x04:
			host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
		default:
			return
		}
		port := strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))

		server, err := net.Dial("tcp", net.JoinHostPort(host, port))
		if server != nil {
			defer func() {
				_ = server.Close()
			}()
		}
		if err != nil {
			return
		}
		_, _ = connect.Write(PackageConnectSuccess)

		go func() {
			_, _ = io.Copy(server, connect)
		}()
		_, _ = io.Copy(connect, server)
	}
}
