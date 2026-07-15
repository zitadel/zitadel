package ldap

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"sync"
	"testing"
	"time"

	ber "github.com/go-asn1-ber/asn1-ber"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootCAPool(t *testing.T) {
	t.Parallel()

	_, caPEM, err := generateTestCert("localhost")
	require.NoError(t, err)

	t.Run("empty returns nil", func(t *testing.T) {
		pool, err := rootCAPool(nil)
		require.NoError(t, err)
		assert.Nil(t, pool)
	})

	t.Run("valid pem", func(t *testing.T) {
		pool, err := rootCAPool(caPEM)
		require.NoError(t, err)
		require.NotNil(t, pool)
	})

	t.Run("invalid pem", func(t *testing.T) {
		pool, err := rootCAPool([]byte("not-a-cert"))
		assert.ErrorIs(t, err, ErrUnableToAppendRootCA)
		assert.Nil(t, pool)
	})
}

func TestGetConnection_RootCA(t *testing.T) {
	t.Parallel()

	serverCert, caPEM, err := generateTestCert("127.0.0.1")
	require.NoError(t, err)

	t.Run("ldaps with rootCA succeeds", func(t *testing.T) {
		t.Parallel()
		addr := startMockLDAPS(t, serverCert)
		conn, err := getConnection("ldaps://"+addr, false, time.Second, caPEM)
		require.NoError(t, err)
		require.NotNil(t, conn)
		conn.Close()
	})

	t.Run("ldaps without rootCA fails for private CA", func(t *testing.T) {
		t.Parallel()
		addr := startMockLDAPS(t, serverCert)
		conn, err := getConnection("ldaps://"+addr, false, time.Second, nil)
		require.Error(t, err)
		assert.Nil(t, conn)
	})

	t.Run("startTLS with rootCA succeeds", func(t *testing.T) {
		t.Parallel()
		addr := startMockLDAPStartTLS(t, serverCert)
		conn, err := getConnection("ldap://"+addr, true, time.Second, caPEM)
		require.NoError(t, err)
		require.NotNil(t, conn)
		conn.Close()
	})

	t.Run("startTLS without rootCA fails for private CA", func(t *testing.T) {
		t.Parallel()
		addr := startMockLDAPStartTLS(t, serverCert)
		conn, err := getConnection("ldap://"+addr, true, time.Second, nil)
		require.Error(t, err)
		assert.Nil(t, conn)
	})
}

func generateTestCert(commonName string) (tls.Certificate, []byte, error) {
	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, nil, err
	}
	serverKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, nil, err
	}

	now := time.Now()
	caTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "test-ca"},
		NotBefore:             now.Add(-time.Hour),
		NotAfter:              now.Add(24 * time.Hour),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}
	caDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		return tls.Certificate{}, nil, err
	}
	caCert, err := x509.ParseCertificate(caDER)
	if err != nil {
		return tls.Certificate{}, nil, err
	}

	serverTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject:      pkix.Name{CommonName: commonName},
		NotBefore:    now.Add(-time.Hour),
		NotAfter:     now.Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     []string{commonName},
	}
	if ip := net.ParseIP(commonName); ip != nil {
		serverTemplate.IPAddresses = []net.IP{ip}
		serverTemplate.DNSNames = nil
	}
	serverDER, err := x509.CreateCertificate(rand.Reader, serverTemplate, caCert, &serverKey.PublicKey, caKey)
	if err != nil {
		return tls.Certificate{}, nil, err
	}

	cert, err := tls.X509KeyPair(
		pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: serverDER}),
		pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: mustMarshalEC(serverKey)}),
	)
	if err != nil {
		return tls.Certificate{}, nil, err
	}

	caPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caDER})
	return cert, caPEM, nil
}

func mustMarshalEC(key *ecdsa.PrivateKey) []byte {
	b, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		panic(err)
	}
	return b
}

func startMockLDAPS(t *testing.T, cert tls.Certificate) string {
	t.Helper()
	ln, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{Certificates: []tls.Certificate{cert}})
	require.NoError(t, err)

	var once sync.Once
	closeLn := func() { once.Do(func() { _ = ln.Close() }) }
	t.Cleanup(closeLn)

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				_ = c.SetDeadline(time.Now().Add(2 * time.Second))
				// Keep the TLS connection open briefly so the client handshake can finish.
				buf := make([]byte, 1)
				_, _ = c.Read(buf)
			}(conn)
		}
	}()

	return ln.Addr().String()
}

func startMockLDAPStartTLS(t *testing.T, cert tls.Certificate) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	var once sync.Once
	closeLn := func() { once.Do(func() { _ = ln.Close() }) }
	t.Cleanup(closeLn)

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go handleStartTLSUpgrade(conn, cert)
		}
	}()

	return ln.Addr().String()
}

func handleStartTLSUpgrade(conn net.Conn, cert tls.Certificate) {
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(2 * time.Second))

	packet, err := ber.ReadPacket(conn)
	if err != nil {
		return
	}
	if len(packet.Children) == 0 {
		return
	}
	messageID := packet.Children[0].Value.(int64)

	response := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSequence, nil, "LDAP Response")
	response.AppendChild(ber.NewInteger(ber.ClassUniversal, ber.TypePrimitive, ber.TagInteger, messageID, "MessageID"))
	extended := ber.Encode(ber.ClassApplication, ber.TypeConstructed, 24, nil, "Extended Response")
	extended.AppendChild(ber.NewInteger(ber.ClassUniversal, ber.TypePrimitive, ber.TagInteger, 0, "Result Code"))
	extended.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "", "Matched DN"))
	extended.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "", "Diagnostic Message"))
	response.AppendChild(extended)
	if _, err := conn.Write(response.Bytes()); err != nil {
		return
	}

	tlsConn := tls.Server(conn, &tls.Config{Certificates: []tls.Certificate{cert}})
	if err := tlsConn.Handshake(); err != nil {
		return
	}
	defer tlsConn.Close()
	buf := make([]byte, 1)
	_, _ = tlsConn.Read(buf)
}
