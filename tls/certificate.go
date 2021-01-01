// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/31

package tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"sort"
	"strings"

	"github.com/crochee/proxy/logger"
	"github.com/crochee/proxy/tls/generate"
)

// Certificate holds a SSL cert/key pair
// Certs and Key could be either a file path, or the file content itself.
type Certificate struct {
	CertFile FileOrContent `json:"certFile,omitempty" toml:"certFile,omitempty" yaml:"certFile,omitempty"`
	KeyFile  FileOrContent `json:"keyFile,omitempty" toml:"keyFile,omitempty" yaml:"keyFile,omitempty"`
}

type Certificates []Certificate

// GetCertificates retrieves the certificates as slice of tls.Certificate.
func (c Certificates) GetCertificates() []tls.Certificate {
	var certs []tls.Certificate

	for _, certificate := range c {
		cert, err := certificate.GetCertificate()
		if err != nil {
			logger.Debugf("Error while getting certificate: %v", err)
			continue
		}

		certs = append(certs, cert)
	}

	return certs
}

// GetCertificate retrieves Certificate as tls.Certificate.
func (c *Certificate) GetCertificate() (tls.Certificate, error) {
	certContent, err := c.CertFile.Read()
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("unable to read CertFile : %w", err)
	}

	var keyContent []byte
	if keyContent, err = c.KeyFile.Read(); err != nil {
		return tls.Certificate{}, fmt.Errorf("unable to read KeyFile : %w", err)
	}
	var cert tls.Certificate
	if cert, err = tls.X509KeyPair(certContent, keyContent); err != nil {
		return tls.Certificate{}, fmt.Errorf("unable to generate TLS certificate : %w", err)
	}
	return cert, nil
}

// GetTruncatedCertificateName truncates the certificate name.
func (c *Certificate) GetTruncatedCertificateName() string {
	certName := c.CertFile.String()

	// Truncate certificate information only if it's a well formed certificate content with more than 50 characters
	if !c.CertFile.IsPath() && strings.HasPrefix(certName, certificateHeader) && len(certName) > len(certificateHeader)+50 {
		certName = strings.TrimPrefix(c.CertFile.String(), certificateHeader)[:50]
	}

	return certName
}

// String is the method to format the flag's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (c *Certificates) String() string {
	if len(*c) == 0 {
		return ""
	}
	var result []string
	for _, certificate := range *c {
		result = append(result, certificate.CertFile.String()+","+certificate.KeyFile.String())
	}
	return strings.Join(result, ";")
}

// Set is the method to set the flag value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the flag.
// It's a comma-separated list, so we split it.
func (c *Certificates) Set(value string) error {
	certificates := strings.Split(value, ";")
	for _, certificate := range certificates {
		files := strings.Split(certificate, ",")
		if len(files) != 2 {
			return fmt.Errorf("bad certificates format: %s", value)
		}
		*c = append(*c, Certificate{
			CertFile: FileOrContent(files[0]),
			KeyFile:  FileOrContent(files[1]),
		})
	}
	return nil
}

// Type is type of the struct.
func (c *Certificates) Type() string {
	return "certificates"
}

// CreateTLSConfig creates a TLS config from Certificate structures.
func (c *Certificates) CreateTLSConfig(entryPointName string) (*tls.Config, error) {
	config := &tls.Config{}
	domainsCertificates := make(map[string]map[string]*tls.Certificate)

	if c.isEmpty() {
		config.Certificates = []tls.Certificate{}

		cert, err := generate.DefaultCertificate("", "")
		if err != nil {
			return nil, err
		}

		config.Certificates = append(config.Certificates, *cert)
	} else {
		for _, certificate := range *c {
			err := certificate.AppendCertificate(domainsCertificates, entryPointName)
			if err != nil {
				logger.Errorf("Unable to add a certificate to the entryPoint %q : %v", entryPointName, err)
				continue
			}

			for _, certDom := range domainsCertificates {
				for _, cert := range certDom {
					config.Certificates = append(config.Certificates, *cert)
				}
			}
		}
	}
	return config, nil
}

// isEmpty checks if the certificates list is empty.
func (c *Certificates) isEmpty() bool {
	if len(*c) == 0 {
		return true
	}
	var key int
	for _, cert := range *c {
		if len(cert.CertFile.String()) != 0 && len(cert.KeyFile.String()) != 0 {
			break
		}
		key++
	}
	return key == len(*c)
}

// AppendCertificate appends a Certificate to a certificates map keyed by entrypoint.
func (c *Certificate) AppendCertificate(certs map[string]map[string]*tls.Certificate, ep string) error {
	certContent, err := c.CertFile.Read()
	if err != nil {
		return fmt.Errorf("unable to read CertFile : %w", err)
	}
	var keyContent []byte
	if keyContent, err = c.KeyFile.Read(); err != nil {
		return fmt.Errorf("unable to read KeyFile : %w", err)
	}
	var tlsCert tls.Certificate
	if tlsCert, err = tls.X509KeyPair(certContent, keyContent); err != nil {
		return fmt.Errorf("unable to generate TLS certificate : %w", err)
	}

	parsedCert, _ := x509.ParseCertificate(tlsCert.Certificate[0])

	var SANs []string
	if parsedCert.Subject.CommonName != "" {
		SANs = append(SANs, strings.ToLower(parsedCert.Subject.CommonName))
	}
	if parsedCert.DNSNames != nil {
		sort.Strings(parsedCert.DNSNames)
		for _, dnsName := range parsedCert.DNSNames {
			if dnsName != parsedCert.Subject.CommonName {
				SANs = append(SANs, strings.ToLower(dnsName))
			}
		}
	}
	if parsedCert.IPAddresses != nil {
		for _, ip := range parsedCert.IPAddresses {
			if ip.String() != parsedCert.Subject.CommonName {
				SANs = append(SANs, strings.ToLower(ip.String()))
			}
		}
	}
	certKey := strings.Join(SANs, ",")

	certExists := false
	if certs[ep] == nil {
		certs[ep] = make(map[string]*tls.Certificate)
	} else {
		for domains := range certs[ep] {
			if domains == certKey {
				certExists = true
				break
			}
		}
	}
	if certExists {
		logger.Debugf("Skipping addition of certificate for domain(s) %q, to EntryPoint %s, "+
			"as it already exists for this Entrypoint.", certKey, ep)
	} else {
		logger.Debugf("Adding certificate for domain(s) %s", certKey)
		certs[ep][certKey] = &tlsCert
	}

	return err
}
