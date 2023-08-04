package services

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"go.mozilla.org/pkcs7"
	"io"
	"net/http"
	"os"
)

type RootCertificateProviderService interface {
	ProvideCertificates(ctx context.Context) ([]*x509.Certificate, error)
}

type FileReader interface {
	ReadFile(filename string) ([]byte, error)
}

type RealFileReader struct{}

func (reader RealFileReader) ReadFile(filename string) (content []byte, e error) {
	//#nosec G304 - only files specified by the person running the application will be loaded
	return os.ReadFile(filename)
}

type FileRootCertificateRetrieverService struct {
	FilePaths  []string
	FileReader FileReader
}

func (s FileRootCertificateRetrieverService) ProvideCertificates(ctx context.Context) (certs []*x509.Certificate, e error) {
	for _, pemFile := range s.FilePaths {
		bytes, e := s.FileReader.ReadFile(pemFile)
		if e != nil {
			return nil, e
		}

		intermediateCerts, e := parseCertificates(bytes)
		if e != nil {
			return nil, e
		}

		certs = append(certs, intermediateCerts...)
	}
	return
}

type OpcpRootCertificateRetrieverService struct {
	MoRootCertPool string
	MoOPCPToken    string
	HttpClient     *http.Client
}

func (s OpcpRootCertificateRetrieverService) ProvideCertificates(ctx context.Context) (certs []*x509.Certificate, e error) {
	body, err := s.retrieveCertificatesFromUrl(ctx)
	if err != nil {
		return nil, err
	}
	decodedBody, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		return nil, err
	}

	crtChain, err := pkcs7.Parse(decodedBody)
	if err != nil {
		return nil, err
	}

	for _, cert := range crtChain.Certificates {
		if cert.Issuer.String() == cert.Subject.String() {
			certs = append(certs, cert)
		}
	}
	return certs, nil
}

func (s OpcpRootCertificateRetrieverService) retrieveCertificatesFromUrl(ctx context.Context) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.MoRootCertPool, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/pkcs10, application/pkcs7")
	req.Header.Add("Content-Transfer-Encoding", "application/pkcs10")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.MoOPCPToken))

	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, HttpError(resp.StatusCode)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func parseCertificates(pemData []byte) ([]*x509.Certificate, error) {
	var certs []*x509.Certificate
	for {
		cert, rest, err := parseCertificate(pemData)
		if err != nil {
			return nil, err
		}
		if cert == nil {
			break
		}
		certs = append(certs, cert)
		pemData = rest
	}
	return certs, nil
}

func parseCertificate(pemData []byte) (cert *x509.Certificate, rest []byte, err error) {
	block, rest := pem.Decode(pemData)
	if block == nil {
		return
	}
	if block.Type != "CERTIFICATE" {
		return
	}
	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		cert = nil
		return
	}
	return
}