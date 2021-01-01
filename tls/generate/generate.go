// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/31

package generate

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/crochee/proxy/logger"
)

// DefaultDomain proxy domain for the default certificate.
const DefaultDomain = "PROXY DEFAULT CERT"

// DefaultCertificate generates random TLS certificates.
func DefaultCertificate(certPath, keyPath string) (*tls.Certificate, error) {
	randomBytes := make([]byte, 100)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, err
	}
	zBytes := sha256.Sum256(randomBytes)
	z := hex.EncodeToString(zBytes[:sha256.Size])
	domain := fmt.Sprintf("%s.%s.proxy.default", z[:32], z[32:])

	certPEM, keyPEM, err := KeyPair(domain, time.Time{})
	if err != nil {
		return nil, err
	}
	// write to file
	if certPath != "" && keyPath != "" {
		certFile, err := os.Create(certPath)
		if err != nil {
			logger.Errorf("create %s failed.Error:%w", certPath, err)
			return nil, err
		}
		if _, err = certFile.Write(certPEM); err != nil {
			logger.Errorf("write %s failed.Error:%w", certPath, err)
			return nil, err
		}
		if err = certFile.Close(); err != nil {
			logger.Warnf("close %s failed.Error:%w", certPath, err)
		}
		keyFile, err := os.Create(keyPath)
		if err != nil {
			logger.Errorf("create %s failed.Error:%w", keyPath, err)
			return nil, err
		}
		if _, err = keyFile.Write(keyPEM); err != nil {
			logger.Errorf("write %s failed.Error:%w", keyPath, err)
			return nil, err
		}
		if err = keyFile.Close(); err != nil {
			logger.Warnf("close %s failed.Error:%w", keyPath, err)
		}
	}
	var certificate tls.Certificate
	if certificate, err = tls.X509KeyPair(certPEM, keyPEM); err != nil {
		return nil, err
	}

	return &certificate, nil
}

// KeyPair generates cert and key files.
func KeyPair(domain string, expiration time.Time) ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	var certPEM []byte
	if certPEM, err = PemCert(privateKey, domain, expiration); err != nil {
		return nil, nil, err
	}
	return certPEM, keyPEM, nil
}

// PemCert generates PEM cert file.
func PemCert(privateKey *rsa.PrivateKey, domain string, expiration time.Time) ([]byte, error) {
	derBytes, err := derCert(privateKey, expiration, domain)
	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes}), nil
}

func derCert(privateKey *rsa.PrivateKey, expiration time.Time, domain string) ([]byte, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}

	if expiration.IsZero() {
		expiration = time.Now().Add(365 * (24 * time.Hour))
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: DefaultDomain,
		},
		NotBefore: time.Now(),
		NotAfter:  expiration,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageKeyAgreement | x509.KeyUsageDataEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{domain},
	}

	return x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
}
