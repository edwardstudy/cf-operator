package inmemorygenerator

import (
	"fmt"
	"time"

	"code.cloudfoundry.org/cf-operator/pkg/credsgen"
	"github.com/cloudflare/cfssl/cli/genkey"
	"github.com/cloudflare/cfssl/config"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/helpers"
	"github.com/cloudflare/cfssl/initca"
	"github.com/cloudflare/cfssl/signer"
	"github.com/cloudflare/cfssl/signer/local"
	"github.com/pkg/errors"
)

// GenerateCertificate generates a certificate using Cloudflare's TLS toolkit
func (g InMemoryGenerator) GenerateCertificate(name string, request credsgen.CertificateGenerationRequest) (credsgen.Certificate, error) {
	var certificate credsgen.Certificate
	var err error

	if request.IsCA {
		certificate, err = g.generateCACertificate()
		if err != nil {
			return credsgen.Certificate{}, errors.Wrap(err, "Generating certificate")
		}
	} else {
		certificate, err = g.generateCertificate(request)
		if err != nil {
			return credsgen.Certificate{}, errors.Wrap(err, "Generating CA")
		}
	}
	return certificate, nil
}

func (g InMemoryGenerator) generateCertificate(request credsgen.CertificateGenerationRequest) (credsgen.Certificate, error) {
	if !request.CA.IsCA {
		return credsgen.Certificate{}, fmt.Errorf("The passed CA is not a CA")
	}

	cert := credsgen.Certificate{
		IsCA: false,
	}

	// Generate certificate
	certReq := &csr.CertificateRequest{KeyRequest: &csr.BasicKeyRequest{A: g.Algorithm, S: g.Bits}}

	certReq.Hosts = append(certReq.Hosts, request.CommonName)
	for _, name := range request.AlternativeNames {
		certReq.Hosts = append(certReq.Hosts, name)
	}
	certReq.CN = certReq.Hosts[0]

	var signingReq []byte
	sslValidator := &csr.Generator{Validator: genkey.Validator}
	signingReq, privateKey, err := sslValidator.ProcessRequest(certReq)
	if err != nil {
		return credsgen.Certificate{}, errors.Wrap(err, "Generating certicate")
	}

	// Parse CA
	caCert, err := helpers.ParseCertificatePEM([]byte(request.CA.Certificate))
	if err != nil {
		return credsgen.Certificate{}, errors.Wrap(err, "Parsing CA PEM")
	}
	caKey, err := helpers.ParsePrivateKeyPEM([]byte(request.CA.PrivateKey))
	if err != nil {
		return credsgen.Certificate{}, errors.Wrap(err, "Parsing CA private key")
	}

	//Sign certificate
	signingProfile := &config.SigningProfile{
		Usage:        []string{"server auth", "client auth"},
		Expiry:       time.Duration(g.Expiry*24) * time.Hour,
		ExpiryString: fmt.Sprintf("%dh", g.Expiry*24),
	}
	policy := &config.Signing{
		Profiles: map[string]*config.SigningProfile{},
		Default:  signingProfile,
	}

	s, err := local.NewSigner(caKey, caCert, signer.DefaultSigAlgo(caKey), policy)
	if err != nil {
		return credsgen.Certificate{}, errors.Wrap(err, "Creating signer")
	}

	cert.Certificate, err = s.Sign(signer.SignRequest{Request: string(signingReq)})
	if err != nil {
		return credsgen.Certificate{}, errors.Wrap(err, "Signing certificate")
	}
	cert.PrivateKey = privateKey

	return cert, nil
}

func (g InMemoryGenerator) generateCACertificate() (credsgen.Certificate, error) {
	req := &csr.CertificateRequest{
		CA:         &csr.CAConfig{Expiry: fmt.Sprintf("%dh", g.Expiry*24)},
		CN:         "SCF CA",
		KeyRequest: &csr.BasicKeyRequest{A: g.Algorithm, S: g.Bits},
	}
	ca, _, privateKey, err := initca.New(req)
	if err != nil {
		return credsgen.Certificate{}, errors.Wrap(err, "Creating CA")
	}
	cert := credsgen.Certificate{
		IsCA:        true,
		Certificate: ca,
		PrivateKey:  privateKey,
	}

	return cert, nil
}
