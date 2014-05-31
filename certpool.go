package mallory

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"sync"
)

// certificates pool, able to get and create certificates thread safely
type CertPool struct {
	// where to store certificates
	Dir string
	// CA to sign new certs
	CA *tls.Certificate

	data  map[string]*tls.Certificate
	mutex sync.RWMutex
}

// Create pool.
// dir is the path where is save and load certificates
// ca is the root CA to create and sign all new certificates
func NewCertPool(dir string, ca *tls.Certificate) (self *CertPool) {
	self = &CertPool{
		Dir:  dir,
		CA:   ca,
		data: make(map[string]*tls.Certificate),
	}
	return
}

func (self *CertPool) Get(host string) (cert *tls.Certificate, err error) {
	// 1. read from memory
	self.mutex.RLock()
	cert, ok := self.data[host]
	self.mutex.RUnlock()
	if ok {
		return
	}

	// Here we must do the lock, consider the case without lock:
	//   1) go A() and go B() load cert from disk at the same time
	//   2) both A() and B() are failed, and create their own certificate
	//   3) finally the same host may have different certificates...
	self.mutex.Lock()
	defer self.mutex.Unlock()

	// 2. find on disk
	crtnam := path.Join(self.Dir, host+".crt")
	der, err := ioutil.ReadFile(crtnam)
	if err == nil {
		rcert, err := tls.X509KeyPair(self.CA.Certificate[0], der)
		if err == nil {
			self.data[host] = &rcert
			return &rcert, err
		}
	}

	// 3. create new
	signer, err := x509.ParseCertificate(self.CA.Certificate[0])
	if err != nil {
		return
	}

	hash := sha1.Sum([]byte(host))
	signee := &x509.Certificate{
		SerialNumber:          new(big.Int).SetBytes(hash[:]),
		Issuer:                signer.Issuer,
		Subject:               signer.Subject,
		NotBefore:             signer.NotBefore,
		NotAfter:              signer.NotAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// the key section, keep the common name same as the host name client want to connect
	// FIXME: mismatch common name
	signee.Subject.CommonName = host

	// get the private key from CA
	key := self.CA.PrivateKey.(*rsa.PrivateKey)
	// here we use CA private key to create new certs
	der, err = x509.CreateCertificate(rand.Reader, signee, signer, &key.PublicKey, key)
	if err != nil {
		return
	}

	cert = &tls.Certificate{
		Certificate: [][]byte{der, self.CA.Certificate[0]},
		PrivateKey:  key,
	}
	self.data[host] = cert

	// save to disk
	fcrt, err := os.Create(crtnam)
	if err != nil {
		return
	}
	for _, c := range cert.Certificate {
		err = pem.Encode(fcrt, &pem.Block{Type: "CERTIFICATE", Bytes: c})
		if err != nil {
			defer os.Remove(crtnam)
			break
		}
	}
	defer fcrt.Close()
	return
}
