package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"net"
	"time"
)

type Certificate struct {
	ServerKey []byte
	ServerPem []byte
}

func baseCertificate(serialNumber int64) *x509.Certificate {
	return &x509.Certificate{
		SerialNumber: big.NewInt(serialNumber),
		Subject:      pkix.Name{},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
}

func GenerateCACertificate() (*x509.Certificate, *rsa.PrivateKey) {
	ca := baseCertificate(1653)
	ca.Subject.Organization = []string{"Monkey-Cat"}
	ca.SubjectKeyId = []byte{1, 2, 3, 4, 5}
	ca.BasicConstraintsValid = true
	ca.IsCA = true

	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalln("generate CA private key:", err)
	}

	return ca, caKey
}

func GenerateCertificate(ca *x509.Certificate, caKey *rsa.PrivateKey, organizations []string) Certificate {
	cert := baseCertificate(1658)
	cert.Subject.Organization = organizations
	cert.SubjectKeyId = []byte{1, 2, 3, 4, 6}
	cert.DNSNames = append(cert.DNSNames, "127.0.0.1")
	cert.IPAddresses = append(cert.IPAddresses, net.ParseIP("0.0.0.0"))

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalln("generate private key:", err)
	}
	pub := &priv.PublicKey
	privPm := priv
	if caKey != nil {
		privPm = caKey
	}
	ca_b, err := x509.CreateCertificate(rand.Reader, cert, ca, pub, privPm)
	if err != nil {
		log.Fatalln("create certificate:", err)
	}

	return Certificate{
		ServerKey: pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: ca_b,
		}),
		ServerPem: pem.EncodeToMemory(&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		}),
	}
}
