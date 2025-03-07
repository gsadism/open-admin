package osv

import (
	"crypto/rsa"
	"github.com/gsadism/open-admin/pkg/crypto/open_rsa"
	"sync"
)

type _RSA struct {
	public  *rsa.PublicKey
	private *rsa.PrivateKey
	once    sync.Once
}

var Rsa = new(_RSA)

func (r *_RSA) Init(pub *rsa.PublicKey, prv *rsa.PrivateKey) {
	if pub != nil && prv != nil {
		r.once.Do(func() {
			r.public = pub
			r.private = prv
		})
	}
}

func (r *_RSA) PublicKey() *rsa.PublicKey {
	return r.public
}

func (r *_RSA) PublicKeyWithString() (string, error) {
	return open_rsa.PublicKeyWithString(r.public)
}

func (r *_RSA) PrivateKey() *rsa.PrivateKey {
	return r.private
}

func (r *_RSA) PrivateKeyWithString() (string, error) {
	return open_rsa.PrivateKeyWithString(r.private)
}

func (r *_RSA) EncryptWithPKCS1v15(Text string) (string, error) {
	return open_rsa.EncryptWithPKCS1v15(r.public, Text)
}

func (r *_RSA) DecryptWithPKCS1v15(Text string) (string, error) {
	return open_rsa.DecryptWithPKCS1v15(r.private, Text)
}

func (r *_RSA) EncryptOAEP(Text, Label string) (string, error) {
	return open_rsa.EncryptOAEP(r.public, open_rsa.SHA384, Text, Label)
}

func (r *_RSA) DecryptOAEP(Text, Label string) (string, error) {
	return open_rsa.DecryptOAEP(r.private, open_rsa.SHA384, Text, Label)
}

func (r *_RSA) EncryptWithPKCS1v15WithString(Key string, Text string) (string, error) {
	// 加载key
	if pub, err := open_rsa.LoadPKIXPublicKey(Key); err != nil {
		return "", err
	} else {
		return open_rsa.EncryptWithPKCS1v15(pub, Text)
	}
}

func (r *_RSA) DecryptWithPKCS1v15WithString(Key string, Text string) (string, error) {
	// 加载 key
	if prv, err := open_rsa.LoadPKCS8PrivateKey(Key); err != nil {
		return "", err
	} else {
		return open_rsa.DecryptWithPKCS1v15(prv, Text)
	}
}
