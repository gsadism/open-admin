package open_rsa

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"hash"
)

type AbstractAlgorithm string

const (
	SHA1   AbstractAlgorithm = "sha1"
	SHA256 AbstractAlgorithm = "sha256"
	SHA384 AbstractAlgorithm = "sha384"
	SHA512 AbstractAlgorithm = "sha512"
)

// EncryptWithPKCS1v15 :
// 使用 PKCS1v15 加密时,原文长度必须小于等于秘钥模长减去 11 字节数
func EncryptWithPKCS1v15(PublicKey *rsa.PublicKey, Text string) (string, error) {
	PlainText := []byte(Text)
	PublicKeySize, SrcSize := PublicKey.Size(), len(PlainText) // 秘钥模长,E 代表指数
	OffSet, Once := 0, PublicKeySize-11
	buffer := bytes.Buffer{}
	if SrcSize <= PublicKeySize-11 {
		// 不需要分段加密
		if ByteOnce, err := rsa.EncryptPKCS1v15(rand.Reader, PublicKey, PlainText); err != nil {
			return "", err
		} else {
			buffer.Write(ByteOnce)
		}
	} else {
		// 进行分段加密
		for OffSet < SrcSize {
			EndIndex := OffSet + Once
			if EndIndex > SrcSize {
				EndIndex = SrcSize
			}
			// 加密一部分
			if ByteOnce, err := rsa.EncryptPKCS1v15(rand.Reader, PublicKey, PlainText[OffSet:EndIndex]); err != nil {
				return "", err
			} else {
				buffer.Write(ByteOnce)
				OffSet = EndIndex
			}
		}
	}
	BytesEncrypt := buffer.Bytes()
	// 将加密内容进行base64编码
	return base64.StdEncoding.EncodeToString(BytesEncrypt), nil
}

// DecryptWithPKCS1v15 :
func DecryptWithPKCS1v15(PrivateKey *rsa.PrivateKey, Text string) (string, error) {
	// 将文本内容进行base64解密
	if CipherText, err := base64.StdEncoding.DecodeString(Text); err != nil {
		return "", err
	} else {
		PrivateKeySize, SrcSize := PrivateKey.Size(), len(CipherText)
		OffSet := 0
		buffer := bytes.Buffer{}
		if SrcSize <= PrivateKeySize {
			// 不需要分段解密
			if BytesOnce, err := rsa.DecryptPKCS1v15(rand.Reader, PrivateKey, CipherText); err != nil {
				return "", err
			} else {
				buffer.Write(BytesOnce)
			}
		} else {
			for OffSet < SrcSize {
				EndIndex := OffSet + PrivateKeySize
				if EndIndex > SrcSize {
					EndIndex = SrcSize
				}
				if BytesOnce, err := rsa.DecryptPKCS1v15(rand.Reader, PrivateKey, CipherText[OffSet:EndIndex]); err != nil {
					return "", err
				} else {
					buffer.Write(BytesOnce)
					OffSet = EndIndex
				}
			}
		}
		BytesDecrypt := buffer.Bytes()
		return string(BytesDecrypt), nil
	}
}

// EncryptOAEP :
// 原文长度必须小于等于密钥模长-(2 * 摘要值长度)减- 2 字节
func EncryptOAEP(PublicKey *rsa.PublicKey, alg AbstractAlgorithm, Text string, label string) (string, error) {
	PlainText := []byte(Text)
	Label := []byte(label)
	AbstractAlgorithmMap := map[AbstractAlgorithm]hash.Hash{
		"sha1":   sha1.New(),
		"sha256": sha256.New(),
		"sha384": sha512.New384(),
		"sha512": sha512.New(),
	}
	var Algorithm hash.Hash
	if hs, ok := AbstractAlgorithmMap[alg]; ok {
		Algorithm = hs
	} else {
		Algorithm = sha1.New()
	}

	PublicKeySize, SrcSize := PublicKey.Size(), len(PlainText)
	OffSet, Once := 0, PublicKeySize-(2*Algorithm.Size())-2
	buffer := bytes.Buffer{}

	if SrcSize <= PublicKeySize-(2*Algorithm.Size())-2 {
		// 不需要分段加密
		if ByteOnce, err := rsa.EncryptOAEP(Algorithm, rand.Reader, PublicKey, PlainText, Label); err != nil {
			return "", err
		} else {
			buffer.Write(ByteOnce)
		}
	} else {
		for OffSet < SrcSize {
			EndIndex := OffSet + Once
			if EndIndex > SrcSize {
				EndIndex = SrcSize
			}
			if ByteOnce, err := rsa.EncryptOAEP(Algorithm, rand.Reader, PublicKey, PlainText[OffSet:EndIndex], Label); err != nil {
				return "", err
			} else {
				buffer.Write(ByteOnce)
				OffSet = EndIndex
			}
		}
	}
	BytesEncrypt := buffer.Bytes()
	return base64.StdEncoding.EncodeToString(BytesEncrypt), nil
}

// DecryptOAEP :
func DecryptOAEP(PrivateKey *rsa.PrivateKey, alg AbstractAlgorithm, Text string, label string) (string, error) {
	Label := []byte(label)
	AbstractAlgorithmMap := map[AbstractAlgorithm]hash.Hash{
		"sha1":   sha1.New(),
		"sha256": sha256.New(),
		"sha384": sha512.New384(),
		"sha512": sha512.New(),
	}
	var Algorithm hash.Hash
	if hs, ok := AbstractAlgorithmMap[alg]; ok {
		Algorithm = hs
	} else {
		Algorithm = sha1.New()
	}
	if CipherText, err := base64.StdEncoding.DecodeString(Text); err != nil {
		return "", err
	} else {
		PrivateKeySize, SrcSize := PrivateKey.Size(), len(CipherText)
		OffSet := 0
		buffer := bytes.Buffer{}
		if SrcSize < PrivateKeySize {
			// 不需要分段
			if cipher, err := rsa.DecryptOAEP(Algorithm, rand.Reader, PrivateKey, CipherText, Label); err != nil {
				return "", err
			} else {
				buffer.Write(cipher)
			}
		} else {
			for OffSet < SrcSize {
				EndIndex := OffSet + PrivateKeySize
				if EndIndex > SrcSize {
					EndIndex = SrcSize
				}
				if BytesOnce, err := rsa.DecryptOAEP(Algorithm, rand.Reader, PrivateKey, CipherText[OffSet:EndIndex], Label); err != nil {
					return "", err
				} else {
					buffer.Write(BytesOnce)
					OffSet = EndIndex
				}
			}
		}
		BytesDecrypt := buffer.Bytes()
		return string(BytesDecrypt), nil
	}
}
