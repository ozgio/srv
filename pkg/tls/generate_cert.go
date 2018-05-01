// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Generate a self-signed X.509 certificate for a TLS server. Outputs to
// 'cert.pem' and 'key.pem' and will overwrite existing files.

// modified by @ozgio to be used as a library

package tls

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

//Options is the key and certficate options for generate functions
type Options struct {
	Host string
	//ValidFrom is creation date formatted as Jan 1 15:04:05 2011
	ValidFrom string
	//ValidFor is the duration that certificate is valid for
	ValidFor time.Duration
	//IsCA is whether this cert should be its own Certificate Authority
	IsCA bool
	//RSABits is the size of RSA key to generate. Ignored if EcdsaCurve is set
	RSABits int
	//ECDSA is the ECDSA curve to use to generate a key. Valid values are P224,
	//P256 (recommended), P384, P521
	ECDSACurve string
}

//GenerateCertToFiles generates key and cert and write them to the files
//specified with parameters keyFile and certFile
func GenerateCertToFiles(keyFile string, certFile string, opts Options) error {
	keyOut, err := os.OpenFile(keyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Printf("failed to open %sfor writing: %s", keyFile, err.Error())
		return err
	}

	certOut, err := os.Create(certFile)
	if err != nil {
		log.Printf("failed to open %s for writing: %s", certFile, err.Error())
		return err
	}

	var errs [3]error
	errs[0] = GenerateCert(certOut, keyOut, opts)

	errs[1] = keyOut.Close()
	errs[2] = certOut.Close()

	var errStrings []string
	for _, e := range errs {
		if e != nil {
			errStrings = append(errStrings, e.Error())
		}
	}
	if len(errStrings) > 0 {
		return errors.New(strings.Join(errStrings, "; "))
	}

	return nil
}

//GenerateCert generates key.pem, cert.pem files and write them through io.Writer
//interfaces
func GenerateCert(certOut io.Writer, keyOut io.Writer, opts Options) error {
	if opts.ValidFor == 0 {
		opts.ValidFor = 365 * 24 * time.Hour
	}
	if opts.RSABits == 0 {
		opts.RSABits = 2048
	}

	if len(opts.Host) == 0 {
		return errors.New("Options.Host must be set")
	}

	var priv interface{}
	var err error
	switch opts.ECDSACurve {
	case "":
		priv, err = rsa.GenerateKey(rand.Reader, opts.RSABits)
	case "P224":
		priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case "P256":
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "P384":
		priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "P521":
		priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		return fmt.Errorf("Unrecognized elliptic curve: %q", opts.ECDSACurve)
	}
	if err != nil {
		return fmt.Errorf("failed to generate private key: %s", err.Error())
	}

	var notBefore time.Time
	if len(opts.ValidFrom) == 0 {
		notBefore = time.Now()
	} else {
		notBefore, err = time.Parse("Jan 2 15:04:05 2006", opts.ValidFrom)
		if err != nil {
			return fmt.Errorf("failed to parse creation date: %s", err.Error())
		}
	}

	notAfter := notBefore.Add(opts.ValidFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %s", err.Error())
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(opts.Host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	if opts.IsCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		return fmt.Errorf("Failed to create certificate: %s", err)
	}

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return fmt.Errorf("failed to write data to cert file: %s", err)
	}

	if err := pem.Encode(keyOut, pemBlockForKey(priv)); err != nil {
		return fmt.Errorf("failed to write data to key file: %s", err)
	}

	return nil

}
