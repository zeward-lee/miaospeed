package interfaces

import (
	"strings"
)

type ProxyType string

const (
	Shadowsocks  ProxyType = "Shadowsocks"
	ShadowsocksR ProxyType = "ShadowsocksR"
	Snell        ProxyType = "Snell"
	Socks5       ProxyType = "Socks5"
	Http         ProxyType = "Http"
	Vmess        ProxyType = "Vmess"
	Trojan       ProxyType = "Trojan"

	Vless       ProxyType = "Vless"
	Hysteria    ProxyType = "Hysteria"
	Hysteria2   ProxyType = "Hysteria2"
	TUIC        ProxyType = "TUIC"
	Wireguard   ProxyType = "Wireguard"
	SSH         ProxyType = "SSH"
	Mieru       ProxyType = "Mieru"
	AnyTLS      ProxyType = "AnyTLS"
	Sudoku      ProxyType = "Sudoku"
	Masque      ProxyType = "Masque"
	TrustTunnel ProxyType = "TrustTunnel"
	GostRelay   ProxyType = "GostRelay"
	ShadowQuic  ProxyType = "ShadowQuic"
	X365        ProxyType = "X365"
	Miu         ProxyType = "Miu"

	// Vendor variants served by our own core (Clash.Meta v1.19.34 + private
	// patches). Every string below is the exact value constant.AdapterType
	// String() returns in that core; the whitelist check is case-insensitive but
	// the spelling is kept identical on purpose.
	Ninja      ProxyType = "Ninja"
	Oppa       ProxyType = "Oppa"
	XHttp      ProxyType = "XHttp"
	Shanlian   ProxyType = "Shanlian"
	Fastup     ProxyType = "Fastup"
	ViewTurbo  ProxyType = "ViewTurbo"
	Blackstone ProxyType = "Blackstone"
	OIX        ProxyType = "OIX"
	Juzi       ProxyType = "Juzi"
	Pure       ProxyType = "Pure"
	// SNTPEclipse is the core's `sntp-eclipse` adapter, whose AdapterType is
	// rendered as "SNTP" (not "SNTPEclipse"), so the whitelist entry follows the
	// rendered name.
	SNTP ProxyType = "SNTP"

	// OpenVPN support is currently treated as experimental in miaospeed.
	OpenVPN ProxyType = "OpenVPN"

	ProxyInvalid ProxyType = "Invalid"
)

var AllProxyTypes = []ProxyType{
	Shadowsocks, ShadowsocksR, Snell, Socks5, Http, Vmess, Trojan,
	Vless, Hysteria, Hysteria2, TUIC, Wireguard, SSH, Mieru, AnyTLS, Sudoku, Masque, TrustTunnel, GostRelay, ShadowQuic, X365, Miu,
	Ninja, Oppa, XHttp, Shanlian, Fastup, ViewTurbo, Blackstone, OIX, Juzi, Pure, SNTP,
	OpenVPN,
}

func (pt *ProxyType) Equal(other ProxyType) bool {
	ptStr := strings.ToLower(string(*pt))
	otherStr := strings.ToLower(string(other))
	return ptStr == otherStr
}

func Valid(proxyType ProxyType) bool {
	for _, pt := range AllProxyTypes {
		if pt.Equal(proxyType) {
			return true
		}
	}
	return false
}

func Parse(proxyType string) ProxyType {
	pType := ProxyType(proxyType)
	if Valid(pType) {
		return pType
	}
	return ProxyInvalid
}

type ProxyInfo struct {
	Name    string
	Address string
	Type    ProxyType
}

func (pi *ProxyInfo) Map() map[string]string {
	return map[string]string{
		"Name":    pi.Name,
		"Address": pi.Address,
		"Type":    string(pi.Type),
	}
}
