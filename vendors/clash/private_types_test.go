package clash

import (
	"testing"

	"github.com/airportr/miaospeed/interfaces"
	"github.com/metacubex/mihomo/adapter"
)

// TestPrivateVendorTypesAcceptedByWhitelist proves the two gates agree end to
// end: the kernel parses the node and renders its own adapter identity through
// AdapterType.String(), and interfaces.Parse (the MiaoSpeed whitelist that
// vendors/clash/profile.go applies) accepts exactly that rendered string and
// maps it back to the expected ProxyType.
func TestPrivateVendorTypesAcceptedByWhitelist(t *testing.T) {
	cases := []struct {
		name    string
		mapping map[string]any
		want    interfaces.ProxyType
	}{
		{
			name: "shanlian",
			mapping: map[string]any{
				"name":     "Shanlian fixture",
				"type":     "shanlian",
				"server":   "node.example",
				"port":     443,
				"password": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			},
			want: interfaces.Shanlian,
		},
		{
			name: "fastup",
			mapping: map[string]any{
				"name":     "Fastup fixture",
				"type":     "fastup",
				"server":   "node.example",
				"port":     443,
				"password": "synthetic-password",
				"sni":      "www.example.com",
			},
			want: interfaces.Fastup,
		},
		{
			name: "ss viewturbo",
			mapping: map[string]any{
				"name":     "ViewTurbo fixture",
				"type":     "ss",
				"server":   "node.example",
				"port":     8388,
				"password": "synthetic-token#VT",
				"cipher":   "chacha20-ietf-poly1305",
			},
			want: interfaces.ViewTurbo,
		},
		{
			name: "ss blackstone",
			mapping: map[string]any{
				"name":     "Blackstone fixture",
				"type":     "ss",
				"server":   "node.example",
				"port":     8388,
				"password": "synthetic-password#BLACKSTONE",
			},
			want: interfaces.Blackstone,
		},
		{
			name: "snell oix",
			mapping: map[string]any{
				"name":   "OIX fixture",
				"type":   "snell",
				"server": "node.example",
				"port":   8388,
				"psk":    "synthetic-psk#oix",
			},
			want: interfaces.OIX,
		},
		{
			name: "vless x365",
			mapping: map[string]any{
				"name":   "X365 fixture",
				"type":   "vless",
				"server": "node.example",
				"port":   443,
				"uuid":   "00112233-4455-6677-8899-aabbccddeeff#x365",
			},
			want: interfaces.X365,
		},
		{
			name: "vless juzi",
			mapping: map[string]any{
				"name":   "Juzi fixture",
				"type":   "vless",
				"server": "node.example",
				"port":   443,
				"uuid":   "6ac24745-6118-1ce8-57da-1aab6a9b7b56#juzi",
			},
			want: interfaces.Juzi,
		},
		{
			name: "vless pure",
			mapping: map[string]any{
				"name":    "Pure fixture",
				"type":    "vless",
				"server":  "node.example",
				"port":    443,
				"uuid":    "00112233-4455-6677-8899-aabbccddeeff#pure",
				"tls":     true,
				"network": "ws",
			},
			want: interfaces.Pure,
		},
		{
			name: "ninja",
			mapping: map[string]any{
				"name":          "Ninja fixture",
				"type":          "ninja",
				"server":        "node.example",
				"port":          443,
				"password":      "synthetic-password",
				"node-password": "synthetic-password",
				"method":        "aes-128-gcm",
			},
			want: interfaces.Ninja,
		},
		{
			name: "oppa",
			mapping: map[string]any{
				"name":     "Oppa fixture",
				"type":     "oppa",
				"server":   "node.example",
				"port":     443,
				"password": "future-token",
				"sni":      "tls.example",
			},
			want: interfaces.Oppa,
		},
		{
			name: "xhttp plain",
			mapping: map[string]any{
				"name":     "XHttp fixture",
				"type":     "xhttp",
				"server":   "192.0.2.10",
				"port":     443,
				"password": "secret:0f1e2d3c4b5a69788796a5b4c3d2e1f0",
			},
			want: interfaces.XHttp,
		},
		{
			// The vendor's own xhttp nodes carry a gateway endpoint; the core
			// renders those as "Heysocks". A plain xhttp node without the
			// gateway marker stays "Xhttp" (the case above), so the two
			// identities are deliberately kept side by side.
			name: "heysocks gateway",
			mapping: map[string]any{
				"name":     "Heysocks fixture",
				"type":     "xhttp",
				"server":   "192.0.2.10",
				"port":     443,
				"password": "secret:0f1e2d3c4b5a69788796a5b4c3d2e1f0",
				"gateway":  "gateway.example:443",
			},
			want: interfaces.Heysocks,
		},
		{
			// anytls plus the vendor xhttp transport is the dedicated-line
			// shape, which the core renders as "Shanlian".
			name: "shanlian anytls xhttp transport",
			mapping: map[string]any{
				"name":             "Shanlian express fixture",
				"type":             "anytls",
				"server":           "198.51.100.7",
				"port":             6504,
				"password":         "00000000000000000000000000000000",
				"sni":              "mirror.example.com",
				"skip-cert-verify": true,
				"transport":        "xhttp",
				"xhttp-opts": map[string]any{
					"path":     "/rpcv2",
					"mode":     "stream-up",
					"obfs-key": "synthetic-key",
				},
			},
			want: interfaces.Shanlian,
		},
		{
			// A plain anytls node (no transport block) must keep reporting
			// "AnyTLS" even with a 32-hex password and no "#sl" marker.
			name: "anytls plain",
			mapping: map[string]any{
				"name":     "AnyTLS fixture",
				"type":     "anytls",
				"server":   "198.51.100.8",
				"port":     443,
				"password": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
				"sni":      "www.example.com",
			},
			want: interfaces.AnyTLS,
		},
		{
			// The kernel renders this adapter as "SNTP", not "SNTPEclipse".
			name: "sntp-eclipse",
			mapping: map[string]any{
				"name":              "SNTP fixture",
				"type":              "sntp-eclipse",
				"server":            "node.example",
				"port":              443,
				"uuid":              "00112233445566778899aabbccddeeff",
				"server-public-key": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
			},
			want: interfaces.SNTP,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			proxy, err := adapter.ParseProxy(tc.mapping)
			if err != nil {
				t.Fatalf("kernel rejected the node: %v", err)
			}
			defer proxy.Close()

			rendered := proxy.Type().String()
			got := interfaces.Parse(rendered)
			if got == interfaces.ProxyInvalid {
				t.Fatalf("kernel rendered %q but the miaospeed whitelist rejects it", rendered)
			}
			if got != tc.want {
				t.Fatalf("kernel rendered %q, whitelist mapped it to %q, want %q", rendered, got, tc.want)
			}
		})
	}
}
