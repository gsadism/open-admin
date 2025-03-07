package open_rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
)

// Sign : 签名
func Sign(privateKey *rsa.PrivateKey, message []byte, sHash crypto.Hash) ([]byte, error) {
	h := sHash.New()
	h.Write(message)
	if sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, sHash, h.Sum(nil)); err != nil {
		return nil, err
	} else {
		return sign, nil
	}
}

// SignVerify : 验签
func SignVerify(publicKey *rsa.PublicKey, message []byte, sign []byte, sHash crypto.Hash) error {
	h := sHash.New()
	h.Write(message)
	return rsa.VerifyPKCS1v15(publicKey, sHash, h.Sum(nil), sign)
}
