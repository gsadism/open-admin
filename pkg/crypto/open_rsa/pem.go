package open_rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"
)

type Bits string

const (
	SMALL  Bits = "small"
	MEDIUM Bits = "medium"
	LARGE  Bits = "large"
)

func Generate(bits Bits) (*rsa.PrivateKey, *rsa.PublicKey) {
	BitsMap := map[string]int{
		"small":  1024,
		"medium": 2048,
		"large":  4096,
	}
	var PrivateKey *rsa.PrivateKey
	if bit, ok := BitsMap[string(bits)]; ok {
		PrivateKey, _ = rsa.GenerateKey(rand.Reader, bit)
	} else {
		PrivateKey, _ = rsa.GenerateKey(rand.Reader, 1024)
	}
	return PrivateKey, &PrivateKey.PublicKey
}

// PublicKeyWithString :
func PublicKeyWithString(PublicKey *rsa.PublicKey) (string, error) {
	if asn, err := x509.MarshalPKIXPublicKey(PublicKey); err != nil {
		return "", err
	} else {
		pemBytes := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: asn,
		})
		return string(pemBytes), nil
	}
}

func PrivateKeyWithString(PrivateKey *rsa.PrivateKey) (string, error) {
	if asn, err := x509.MarshalPKCS8PrivateKey(PrivateKey); err != nil {
		return "", err
	} else {
		pemBytes := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: asn,
		})
		return string(pemBytes), nil
	}
}

func PrivatePKCS1KeyWithString(PrivateKey *rsa.PrivateKey) (string, error) {
	asn := x509.MarshalPKCS1PrivateKey(PrivateKey)
	if asn == nil {
		return "", errors.New("failed to marshal private key")
	} else {
		pemBytes := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: asn,
		})
		return string(pemBytes), nil
	}
}

func PrivatePKCS8KeyWithString(PrivateKey *rsa.PrivateKey) (string, error) {
	if asn, err := x509.MarshalPKCS8PrivateKey(PrivateKey); err != nil {
		return "", err
	} else {
		pemBytes := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: asn,
		})
		return string(pemBytes), nil
	}
}

// WritePKCS1PrivateKeyFile :
func WritePKCS1PrivateKeyFile(PrivateKey *rsa.PrivateKey, filename string, headers map[string]string) error {
	if f, err := os.Create(filename); err != nil {
		return err
	} else {
		defer f.Close()
		if err := pem.Encode(f, &pem.Block{
			Type:    "RSA PRIVATE KEY",
			Headers: headers,
			Bytes:   x509.MarshalPKCS1PrivateKey(PrivateKey),
		}); err != nil {
			return err
		}
		return nil
	}
}

// WritePKCS8PrivateKeyFile :
func WritePKCS8PrivateKeyFile(PrivateKey *rsa.PrivateKey, filename string, headers map[string]string) error {
	if f, err := os.Create(filename); err != nil {
		return err
	} else {
		defer f.Close()
		if Bytes, err := x509.MarshalPKCS8PrivateKey(PrivateKey); err != nil {
			return err
		} else {
			if err := pem.Encode(f, &pem.Block{
				Type:    "RSA PRIVATE KEY",
				Headers: headers,
				Bytes:   Bytes,
			}); err != nil {
				return err
			}
			return nil
		}
	}
}

// WritePKIXPublicKeyFile :
func WritePKIXPublicKeyFile(PublicKey *rsa.PublicKey, filename string, headers map[string]string) error {
	if f, err := os.Create(filename); err != nil {
		return err
	} else {
		defer f.Close()
		if x509PublicKey, err := x509.MarshalPKIXPublicKey(PublicKey); err != nil {
			return err
		} else {
			if err := pem.Encode(f, &pem.Block{
				Type:    "RSA PUBLIC KEY",
				Headers: headers,
				Bytes:   x509PublicKey,
			}); err != nil {
				return err
			}
			return nil
		}
	}
}

// LoadPKCS8PrivateKey : 通过私钥的文本加载私钥 PKCS8
func LoadPKCS8PrivateKey(PrivateKeyStr string) (*rsa.PrivateKey, error) {
	Block, _ := pem.Decode([]byte(PrivateKeyStr))
	if Block == nil {
		return nil, errors.New("decode private key error")
	}
	if PrivateKeyInterface, err := x509.ParsePKCS8PrivateKey(Block.Bytes); err != nil {
		return nil, err
	} else {
		if PrivateKey, ok := PrivateKeyInterface.(*rsa.PrivateKey); ok {
			return PrivateKey, nil
		} else {
			return nil, errors.New("not a RSA private key")
		}
	}
}

// LoadPKCS1PrivateKey : 通过私钥的文本加载私钥 PKCS1
func LoadPKCS1PrivateKey(PrivateKeyStr string) (*rsa.PrivateKey, error) {
	Block, _ := pem.Decode([]byte(PrivateKeyStr))
	if Block == nil {
		return nil, errors.New("decode private key error")
	}
	//if Block.Type != "RSA PRIVATE KEY" {
	//	return nil, errors.New("the kind of PEM should be RSA PRIVATE KEY")
	//}
	if PrivateKey, err := x509.ParsePKCS1PrivateKey(Block.Bytes); err != nil {
		return nil, err
	} else {
		return PrivateKey, nil
	}
}

// LoadPKIXPublicKey : 通过公钥的文本加载公钥
func LoadPKIXPublicKey(PublicKeyStr string) (*rsa.PublicKey, error) {
	Block, _ := pem.Decode([]byte(PublicKeyStr))
	if Block == nil {
		return nil, errors.New("decode public key error")
	}
	//if Block.Type != "RSA PUBLIC KEY" {
	//	return nil, errors.New("the kind of PEM should be RSA PUBLIC KEY")
	//}
	if PublicKeyInterface, err := x509.ParsePKIXPublicKey(Block.Bytes); err != nil {
		return nil, err
	} else {
		if PublicKey, ok := PublicKeyInterface.(*rsa.PublicKey); ok {
			return PublicKey, nil
		} else {
			return nil, errors.New("the kind of PEM should be RSA PUBLIC KEY")
		}
	}
}

// LoadPKCS1PrivateKeyWithFile : 通过私钥文件路径加载私钥 PKCS1
func LoadPKCS1PrivateKeyWithFile(filename string) (*rsa.PrivateKey, error) {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
	}
	if OpenFile, err := os.Open(filename); err != nil {
		return nil, err
	} else {
		defer OpenFile.Close()
		// 读取文件
		if Info, err := OpenFile.Stat(); err != nil {
			return nil, err
		} else {
			buf := make([]byte, Info.Size())
			// 读取文件内容
			if _, err := OpenFile.Read(buf); err != nil {
				return nil, err
			} else {
				return LoadPKCS1PrivateKey(string(buf))
			}
		}
	}
}

// LoadPKCS8PrivateKeyWithFile : 通过私钥文件路径加载私钥 PKCS8
func LoadPKCS8PrivateKeyWithFile(filename string) (*rsa.PrivateKey, error) {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
	}
	if OpenFile, err := os.Open(filename); err != nil {
		return nil, err
	} else {
		defer OpenFile.Close()
		// 读取文件
		if Info, err := OpenFile.Stat(); err != nil {
			return nil, err
		} else {
			buf := make([]byte, Info.Size())
			// 读取文件内容
			if _, err := OpenFile.Read(buf); err != nil {
				return nil, err
			} else {
				return LoadPKCS8PrivateKey(string(buf))
			}
		}
	}
}

// LoadPKIXPublicKeyWithFile : 通过公钥文件路径加载公钥
func LoadPKIXPublicKeyWithFile(filename string) (*rsa.PublicKey, error) {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
	}
	if OpenFile, err := os.Open(filename); err != nil {
		return nil, err
	} else {
		defer OpenFile.Close()
		if Info, err := OpenFile.Stat(); err != nil {
			return nil, err
		} else {
			buf := make([]byte, Info.Size())
			if _, err := OpenFile.Read(buf); err != nil {
				return nil, err
			} else {
				return LoadPKIXPublicKey(string(buf))
			}
		}
	}
}

// LoadCertificate : 通过证书的文本加载证书
func LoadCertificate(CertificateStr string) (*x509.Certificate, error) {
	Block, _ := pem.Decode([]byte(CertificateStr))
	if Block == nil {
		return nil, errors.New("decode certificate error")
	}
	if Certificate, err := x509.ParseCertificate(Block.Bytes); err != nil {
		return nil, err
	} else {
		return Certificate, nil
	}
}

// LoadCertificateWithFile : 通过证书文件路径加载证书
func LoadCertificateWithFile(filename string) (*x509.Certificate, error) {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
	}
	if OpenFile, err := os.Open(filename); err != nil {
		return nil, err
	} else {
		defer OpenFile.Close()
		if Info, err := OpenFile.Stat(); err != nil {
			return nil, err
		} else {
			buf := make([]byte, Info.Size())
			if _, err := OpenFile.Read(buf); err != nil {
				return nil, err
			} else {
				return LoadCertificate(string(buf))
			}
		}
	}
}

// CertificateSerialNumber : 从证书中获取序列号
func CertificateSerialNumber(Certificate x509.Certificate) string {
	return fmt.Sprintf("%X", Certificate.SerialNumber.Bytes())
}

// IsCertificateExpired 判定证书在特定时间是否过期
func IsCertificateExpired(certificate x509.Certificate, now time.Time) bool {

	return now.After(certificate.NotAfter)
}

// IsCertificateValid 判定证书在特定时间是否有效
func IsCertificateValid(certificate x509.Certificate, now time.Time) bool {
	return now.After(certificate.NotBefore) && now.Before(certificate.NotAfter)
}
