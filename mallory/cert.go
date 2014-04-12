package mallory

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"math/big"
	"sync"
)

// simple config used to create signed certificate
type CertConfig struct {
	SerialNumber *big.Int
	CommonName   string
}

// root is a certificate loaded from external
// return a new certificate signed by the root CA and the root CA chain
func CreateSignedCert(root *tls.Certificate, config *CertConfig) (cert *tls.Certificate, err error) {
	// the first one is the signer
	signer, err := x509.ParseCertificate(root.Certificate[0])
	if err != nil {
		return
	}

	signee := &x509.Certificate{
		SerialNumber:          config.SerialNumber,
		Issuer:                signer.Issuer,
		Subject:               signer.Subject,
		NotBefore:             signer.NotBefore,
		NotAfter:              signer.NotAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// the key section, keep the common name same as the host name client want to connect
	signee.Subject.CommonName = config.CommonName

	// here use the root public key as the signee's to generate certificate
	priv := root.PrivateKey.(*rsa.PrivateKey)
	crtder, err := x509.CreateCertificate(rand.Reader, signee, signer, &priv.PublicKey, root.PrivateKey)
	if err != nil {
		return
	}

	cert = &tls.Certificate{
		Certificate: [][]byte{crtder, root.Certificate[0]},
		PrivateKey:  priv,
	}
	return
}

// certificates pool, able to thread safely add and fetch
type CertPool struct {
	data  map[string]*tls.Certificate
	mutex sync.RWMutex
}

func NewCertPool() *CertPool {
	pool := &CertPool{
		data: make(map[string]*tls.Certificate),
	}
	return pool
}

func (self *CertPool) AddSafe(key string, cert *tls.Certificate) {
	self.mutex.Lock()
	self.data[key] = cert
	self.mutex.Unlock()
}

func (self *CertPool) GetSafe(key string) *tls.Certificate {
	self.mutex.RLock()
	cert, ok := self.data[key]
	self.mutex.RUnlock()
	if !ok {
		return nil
	}
	return cert
}
