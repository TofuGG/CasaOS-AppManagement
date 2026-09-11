package downloadHelper

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-getter"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

// Download downloads src into dst using go-getter restricted to https only.
//
// SECURITY:
//   - Only the https getter is registered, so file://, git::, s3://, gcs://, hg
//     etc. sources are rejected outright (prevents local file disclosure).
//   - The URL is validated before any network I/O: https scheme, no userinfo,
//     and the host must not resolve to a private/reserved address (prevents
//     SSRF against the host's internal network).
func Download(src string, dst string) error {
	if err := ValidateAppStoreURL(src); err != nil {
		return err
	}

	backgroundCtx := context.Background()
	client := &getter.Client{
		Ctx:  backgroundCtx,
		Src:  src,
		Dst:  dst,
		Mode: getter.ClientModeAny,
		// Restrict to https only — do NOT enable the default file/git/s3/gcs/hg getters.
		Getters: map[string]getter.Getter{
			"https": &getter.HttpGetter{Client: httpClient},
		},
		Options: []getter.ClientOption{},
	}

	return client.Get()
}

// ValidateAppStoreURL checks that an app store source URL is https and does not
// resolve to a private or reserved address.
func ValidateAppStoreURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("invalid app store URL: %s", err.Error())
	}

	if u.Scheme != "https" {
		return fmt.Errorf("app store URL scheme must be https, got %q", u.Scheme)
	}
	if u.Host == "" {
		return fmt.Errorf("app store URL has no host")
	}
	if u.User != nil {
		return fmt.Errorf("app store URL must not contain userinfo")
	}

	host := u.Hostname()
	ips, err := net.LookupHost(host)
	if err != nil {
		return fmt.Errorf("app store URL host %q could not be resolved: %s", host, err.Error())
	}

	for _, ipStr := range ips {
		ip := net.ParseIP(strings.TrimSpace(ipStr))
		if ip == nil {
			continue
		}
		if !isPublicIP(ip) {
			return fmt.Errorf("app store URL host %q resolves to a private or reserved address (%s)", host, ipStr)
		}
	}

	return nil
}

func isPublicIP(ip net.IP) bool {
	if v4 := ip.To4(); v4 != nil {
		ip = v4
	}
	return !(ip.IsLoopback() ||
		ip.IsPrivate() ||
		ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() ||
		ip.IsMulticast() ||
		ip.IsUnspecified() ||
		ip.Equal(net.IPv4bcast))
}