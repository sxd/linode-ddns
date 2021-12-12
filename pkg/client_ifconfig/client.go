package client_ifconfig

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"strings"
)

const ifconfigURL = "http://ifconfig.co"

type discoveredIP struct {
	ipAddrStr string
	ipV4Addr  net.IP
}

func (d *discoveredIP) GetString() string {
	return d.ipAddrStr
}

func getIPData(ip string) *discoveredIP {
	ipv4 := net.ParseIP(ip)

	return &discoveredIP{
		ipAddrStr: ip,
		ipV4Addr:  ipv4,
	}
}

func readResponse(closer io.ReadCloser) (string, error) {
	content, err := io.ReadAll(closer)
	if err != nil {
		return "", nil
	}
	ip := bytes.NewBuffer(content).String()

	return strings.TrimSpace(ip), nil
}

func callService(ctx context.Context) (string, error) {
	var body io.Reader
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ifconfigURL, body)
	if err != nil {
		return "", err
	}
	defer req.Body.Close()

	return readResponse(req.Body)
}

func GetIP(ctx context.Context) (*discoveredIP, error) {
	ip, err := callService(ctx)
	if err != nil {
		return nil, err
	}

	ipDiscovered := getIPData(ip)

	return ipDiscovered, nil
}
