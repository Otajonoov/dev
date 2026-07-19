package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
)

// SMPP over TLS (16-bob). SMPP'ning o'zida shifrlash YO'Q: parol (max 9
// bayt!) va xabar matni simda ochiq oqadi — himoya transport qatlamiga
// qoldirilgan. IANA'da faqat 2775 (smpp) ro'yxatdan o'tgan; 3550 ("ssmpp")
// — keng tarqalgan DE-FACTO konventsiya, IANA birlamchi manbasidan
// tasdiqlanmagan; amalda operatorlar ixtiyoriy port beradi.

// DefaultTLSConfig — production uchun oqilona boshlang'ich: TLS 1.2+.
// InsecureSkipVerify YO'Q va bo'lmaydi — uni test'dan prod'ga ko'chirish
// MITM himoyasini o'chiradigan klassik xato; test uchun o'z CA'ingizni
// RootCAs'ga qo'shing.
func DefaultTLSConfig(serverName string) *tls.Config {
	return &tls.Config{
		ServerName: serverName,
		MinVersion: tls.VersionTLS12,
	}
}

// dial TCP yoki (Config.TLS berilgan bo'lsa) TLS ulanish ochadi.
func (c *Client) dial(ctx context.Context) (net.Conn, error) {
	d := net.Dialer{}
	if c.cfg.TLS == nil {
		conn, err := d.DialContext(ctx, "tcp", c.cfg.Addr)
		if err != nil {
			return nil, fmt.Errorf("client: dial %s: %w", c.cfg.Addr, err)
		}
		return conn, nil
	}
	td := tls.Dialer{NetDialer: &d, Config: c.cfg.TLS}
	conn, err := td.DialContext(ctx, "tcp", c.cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("client: TLS dial %s: %w", c.cfg.Addr, err)
	}
	return conn, nil
}
