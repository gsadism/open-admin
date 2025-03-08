package open_aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

// ECB :
// 工作原理：每个分组独立加密，直接将明文分块，然后逐块加密，每块加密结果相同。
// 优点：简单，易于实现，支持并行处理。
// 缺点：缺乏随机性，明文相同则密文相同，容易被模式攻击，不适合加密大量相似数据。
// 应用场景：一般不建议用于数据加密，但可以用于简单数据结构的加密，如加密一个简单的ID或固定数据块。

// CBC :
// 工作原理：每个分组依赖于前一个密文块的结果，即前一个密文块和当前明文块异或后再进行加密。第一个块则用一个初始化向量（IV）进行异或操作。
// 优点：加入 IV 增加了随机性，能抵御模式攻击。
// 缺点：加密时不可并行，因每个分组依赖前一个分组，解密时可以并行。
// 应用场景：是对数据加密比较常用的模式，适用于数据完整性要求较高的场景。

// CFB :
// 工作原理：先对上一个密文块进行加密，然后将加密结果与当前明文块异或以生成密文块。
// 优点：可以加密小于块长度的数据，可以作为流加密使用。
// 缺点：密文错误会在下一个密文块传播。
// 应用场景：适用于流式数据加密，如网络通信、实时数据流等。

// CTR :
// 工作原理：每次加密时，用一个计数器（Counter）值加密生成密钥流，与明文块异或得到密文。
// 优点：支持并行加密和解密，性能高。
// 缺点：计数器必须唯一且不可重复，否则会降低安全性。
// 应用场景：高性能要求的场景，适合需要并行化的数据加密。

// OFB :
// 工作原理：先生成一个密钥流，每次对 IV 或上一个密钥流块加密得到一个新的密钥块，然后将其与明文异或生成密文。
// 优点：与 CFB 类似，可用作流加密，且不会在密文中传播错误。
// 缺点：模式攻击可能较为容易。
// 应用场景：适用于流加密且要求错误不传播的场景。

// GCM :
// 工作原理：在 CTR 模式基础上加入认证功能，通过加密和认证来确保数据的机密性和完整性。
// 优点：具有认证功能，能同时保证加密数据的完整性和安全性。
// 缺点：实现复杂度高。
// 应用场景：适用于需要认证和加密的场景，如 HTTPS、VPN。

func Secret(size int) (string, error) {
	nonce := make([]byte, size)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	} else {
		return hex.EncodeToString(nonce), nil
	}
}

// pkcs7padding 填充
func pkcs7padding(Data []byte, BlockSize int) []byte {
	Padding := BlockSize - len(Data)%BlockSize
	PadText := bytes.Repeat([]byte{byte(Padding)}, Padding)
	return append(Data, PadText...)
}

func pkcs7unPadding(Data []byte) []byte {
	length := len(Data)
	unPadding := int(Data[length-1])
	return Data[:(length - unPadding)]
}

func EncryptECB(Data, Key string) (string, error) {
	if k, err := hex.DecodeString(Key); err != nil {
		return "", err
	} else {
		if Block, err := aes.NewCipher(k); err != nil {
			return "", err
		} else {
			BlockSize := Block.BlockSize()
			PlainText := pkcs7padding([]byte(Data), BlockSize)
			CipherText := make([]byte, len(PlainText))
			for start := 0; start < len(PlainText); start += BlockSize {
				end := start + BlockSize
				Block.Encrypt(CipherText[start:end], PlainText[start:end])
			}
			return base64.StdEncoding.EncodeToString(CipherText), nil
		}
	}
}

func DecryptECB(Data, Key string) (string, error) {
	if k, err := hex.DecodeString(Key); err != nil {
		return "", err
	} else {
		if Block, err := aes.NewCipher(k); err != nil {
			return "", err
		} else {
			// 将密文进行base64解密
			if CipherText, err := base64.StdEncoding.DecodeString(Data); err != nil {
				return "", err
			} else {
				BlockSize := Block.BlockSize()
				PlainText := make([]byte, len(CipherText))
				for start := 0; start < len(CipherText); start += BlockSize {
					end := start + BlockSize
					Block.Decrypt(PlainText[start:end], CipherText[start:end])
				}
				// 解密后的内容
				return string(pkcs7unPadding(PlainText)), nil
			}
		}
	}
}

func EncryptCBC(Data, Key, IV string) (string, error) {
	k, err := hex.DecodeString(Key)
	if err != nil {
		return "", err
	}
	i, err := hex.DecodeString(IV)
	if err != nil {
		return "", err
	}
	if Block, err := aes.NewCipher(k); err != nil {
		return "", err
	} else {
		BlockSize := Block.BlockSize()
		PlainText := pkcs7padding([]byte(Data), BlockSize)
		CipherText := make([]byte, len(PlainText))
		cipher.NewCBCEncrypter(Block, i).CryptBlocks(CipherText, PlainText)
		return base64.StdEncoding.EncodeToString(CipherText), nil
	}
}

func DecryptCBC(Data, Key, IV string) (string, error) {
	k, err := hex.DecodeString(Key)
	if err != nil {
		return "", err
	}
	i, err := hex.DecodeString(IV)
	if err != nil {
		return "", err
	}
	if Block, err := aes.NewCipher(k); err != nil {
		return "", err
	} else {
		// 将密文进行base64解密
		if CipherText, err := base64.StdEncoding.DecodeString(Data); err != nil {
			return "", err
		} else {
			BlockSize := Block.BlockSize()
			if len(CipherText)%BlockSize != 0 {
				return "", errors.New("ciphertext is not a multiple of the block size")
			}
			PlainText := make([]byte, len(CipherText))
			cipher.NewCBCDecrypter(Block, i).CryptBlocks(PlainText, CipherText)
			return string(pkcs7unPadding(PlainText)), nil
		}
	}
}

func EncryptCFB(Data, Key, IV string) (string, error) {
	k, err := hex.DecodeString(Key)
	if err != nil {
		return "", err
	}
	i, err := hex.DecodeString(IV)
	if err != nil {
		return "", err
	}
	if Block, err := aes.NewCipher(k); err != nil {
		return "", err
	} else {
		CipherText := make([]byte, len([]byte(Data)))
		cipher.NewCFBEncrypter(Block, i).XORKeyStream(CipherText, []byte(Data))
		// 对密文进行base64加密
		return base64.StdEncoding.EncodeToString(CipherText), nil
	}
}

func DecryptCFB(Data, Key, IV string) (string, error) {
	k, err := hex.DecodeString(Key)
	if err != nil {
		return "", err
	}
	i, err := hex.DecodeString(IV)
	if err != nil {
		return "", err
	}
	if Block, err := aes.NewCipher(k); err != nil {
		return "", err
	} else {
		// 将密文进行base64解密
		if CipherText, err := base64.StdEncoding.DecodeString(Data); err != nil {
			return "", err
		} else {
			PlainText := make([]byte, len(CipherText))
			cipher.NewCFBDecrypter(Block, i).XORKeyStream(PlainText, CipherText)
			return string(PlainText), nil
		}
	}
}

func EncryptCTR(Data, Key, IV string) (string, error) {
	k, err := hex.DecodeString(Key)
	if err != nil {
		return "", err
	}
	i, err := hex.DecodeString(IV)
	if err != nil {
		return "", err
	}
	if Block, err := aes.NewCipher(k); err != nil {
		return "", err
	} else {
		CipherText := make([]byte, len([]byte(Data)))
		cipher.NewCTR(Block, i).XORKeyStream(CipherText, []byte(Data))
		// 对密文进行base64加密
		return base64.StdEncoding.EncodeToString(CipherText), nil
	}
}

func DecryptCTR(Data, Key, IV string) (string, error) {
	k, err := hex.DecodeString(Key)
	if err != nil {
		return "", err
	}
	i, err := hex.DecodeString(IV)
	if err != nil {
		return "", err
	}
	if Block, err := aes.NewCipher(k); err != nil {
		return "", err
	} else {
		// 将密文进行base64解密
		if CipherText, err := base64.StdEncoding.DecodeString(Data); err != nil {
			return "", err
		} else {
			PlainText := make([]byte, len(CipherText))
			cipher.NewCTR(Block, i).XORKeyStream(PlainText, CipherText)
			return string(PlainText), nil
		}
	}
}

func EncryptOFB(Data, Key, IV string) (string, error) {
	k, err := hex.DecodeString(Key)
	if err != nil {
		return "", err
	}
	i, err := hex.DecodeString(IV)
	if err != nil {
		return "", err
	}
	if Block, err := aes.NewCipher(k); err != nil {
		return "", err
	} else {
		CipherText := make([]byte, len([]byte(Data)))
		cipher.NewOFB(Block, i).XORKeyStream(CipherText, []byte(Data))
		return base64.StdEncoding.EncodeToString(CipherText), nil
	}
}

func DecryptOFB(Data, Key, IV string) (string, error) {
	k, err := hex.DecodeString(Key)
	if err != nil {
		return "", err
	}
	i, err := hex.DecodeString(IV)
	if err != nil {
		return "", err
	}
	if Block, err := aes.NewCipher(k); err != nil {
		return "", err
	} else {
		// 将密文进行base64解密
		if CipherText, err := base64.StdEncoding.DecodeString(Data); err != nil {
			return "", err
		} else {
			PlainText := make([]byte, len(CipherText))
			cipher.NewOFB(Block, i).XORKeyStream(PlainText, CipherText)
			return string(PlainText), nil
		}
	}
}

func EncryptGCM(Data, Key, IV, Add string) (string, error) {
	k, err := hex.DecodeString(Key)
	if err != nil {
		return "", err
	}
	i, err := hex.DecodeString(IV)
	if err != nil {
		return "", err
	}
	if Block, err := aes.NewCipher(k); err != nil {
		return "", err
	} else {
		if gcm, err := cipher.NewGCM(Block); err != nil {
			return "", err
		} else {
			return base64.StdEncoding.EncodeToString(gcm.Seal(nil, i, []byte(Data), nil)), nil
		}
	}
}

func DecryptGCM(Data, Key, IV, Add string) (string, error) {
	k, err := hex.DecodeString(Key)
	if err != nil {
		return "", err
	}
	i, err := hex.DecodeString(IV)
	if err != nil {
		return "", err
	}
	if Block, err := aes.NewCipher(k); err != nil {
		return "", err
	} else {
		// 将密文进行base64解密
		if CipherText, err := base64.StdEncoding.DecodeString(Data); err != nil {
			return "", err
		} else {
			if gcm, err := cipher.NewGCM(Block); err != nil {
				return "", err
			} else {
				if PlainText, err := gcm.Open(nil, i, CipherText, []byte(Add)); err != nil {
					return "", err
				} else {
					return string(PlainText), nil
				}
			}
		}
	}
}
