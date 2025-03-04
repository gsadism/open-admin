package _rsa

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
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
func EncryptWithPKCS1v15(PublicKey *rsa.PublicKey, PlainText []byte) ([]byte, error) {
	PublicKeySize, SrcSize := PublicKey.Size(), len(PlainText) // 秘钥模长,E 代表指数
	OffSet, Once := 0, PublicKeySize-11
	buffer := bytes.Buffer{}
	if SrcSize <= PublicKeySize-11 {
		// 不需要分段加密
		if ByteOnce, err := rsa.EncryptPKCS1v15(rand.Reader, PublicKey, PlainText); err != nil {
			return nil, err
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
				return nil, err
			} else {
				buffer.Write(ByteOnce)
				OffSet = EndIndex
			}
		}
	}
	BytesEncrypt := buffer.Bytes()
	return BytesEncrypt, nil
}

// DecryptWithPKCS1v15 :
func DecryptWithPKCS1v15(PrivateKey *rsa.PrivateKey, CipherText []byte) ([]byte, error) {
	PrivateKeySize, SrcSize := PrivateKey.Size(), len(CipherText)
	OffSet := 0
	buffer := bytes.Buffer{}
	if SrcSize <= PrivateKeySize {
		// 不需要分段解密
		if BytesOnce, err := rsa.DecryptPKCS1v15(rand.Reader, PrivateKey, CipherText); err != nil {
			return nil, err
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
				return nil, err
			} else {
				buffer.Write(BytesOnce)
				OffSet = EndIndex
			}
		}
	}
	BytesDecrypt := buffer.Bytes()
	return BytesDecrypt, nil
}

// EncryptOAEP :
// 原文长度必须小于等于密钥模长-(2 * 摘要值长度)减- 2 字节
func EncryptOAEP(PublicKey *rsa.PublicKey, alg AbstractAlgorithm, PlainText []byte, label []byte) ([]byte, error) {
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
		if ByteOnce, err := rsa.EncryptOAEP(Algorithm, rand.Reader, PublicKey, PlainText, label); err != nil {
			return nil, err
		} else {
			buffer.Write(ByteOnce)
		}
	} else {
		for OffSet < SrcSize {
			EndIndex := OffSet + Once
			if EndIndex > SrcSize {
				EndIndex = SrcSize
			}
			if ByteOnce, err := rsa.EncryptOAEP(Algorithm, rand.Reader, PublicKey, PlainText[OffSet:EndIndex], label); err != nil {
				return nil, err
			} else {
				buffer.Write(ByteOnce)
				OffSet = EndIndex
			}
		}
	}
	BytesEncrypt := buffer.Bytes()
	return BytesEncrypt, nil
}

// DecryptOAEP :
func DecryptOAEP(PrivateKey *rsa.PrivateKey, alg AbstractAlgorithm, CipherText []byte, label []byte) ([]byte, error) {
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
	PrivateKeySize, SrcSize := PrivateKey.Size(), len(CipherText)
	OffSet := 0
	buffer := bytes.Buffer{}
	if SrcSize < PrivateKeySize {
		// 不需要分段
		if cipher, err := rsa.DecryptOAEP(Algorithm, rand.Reader, PrivateKey, CipherText, label); err != nil {
			return nil, err
		} else {
			buffer.Write(cipher)
		}
	} else {
		for OffSet < SrcSize {
			EndIndex := OffSet + PrivateKeySize
			if EndIndex > SrcSize {
				EndIndex = SrcSize
			}
			if BytesOnce, err := rsa.DecryptOAEP(Algorithm, rand.Reader, PrivateKey, CipherText[OffSet:EndIndex], label); err != nil {
				return nil, err
			} else {
				buffer.Write(BytesOnce)
				OffSet = EndIndex
			}
		}
	}
	BytesDecrypt := buffer.Bytes()
	return BytesDecrypt, nil
}
