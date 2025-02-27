package _aes

import (
	"crypto/aes"
	"fmt"
	"testing"
)

var (
	key = "ea3969af1479eff8e2ab4a0e499eb3aa44a8029990708291156649699ce23bd1"
	iv  = "c19a9538d4c4ae888476788fe7e18916"
)

func TestKey(t *testing.T) {
	fmt.Println(Secret(aes.BlockSize))
}

func TestECB(t *testing.T) {
	data := "hello world"
	if cipher, err := EncryptECB(data, key); err != nil {
		t.Error(err)
	} else {
		t.Logf("加密内容: %s\n", cipher)
		// 解密
		if plain, err := DecryptECB(cipher, key); err != nil {
			t.Error(err)
		} else {
			t.Logf("解密内容: %s\n", plain)
		}
	}
}

func TestCBC(t *testing.T) {
	data := "hello world"
	if cipher, err := EncryptCBC(data, key, iv); err != nil {
		t.Error(err)
	} else {
		t.Logf("加密内容: %s\n", cipher)
		// 解密
		if plain, err := DecryptCBC(cipher, key, iv); err != nil {
			t.Error(err)
		} else {
			t.Logf("解密内容: %s\n", plain)
		}
	}
}

func TestCFB(t *testing.T) {
	data := "hello world"
	if cipher, err := EncryptCFB(data, key, iv); err != nil {
		t.Error(err)
	} else {
		t.Logf("加密内容: %s\n", cipher)
		// 解密
		if plain, err := DecryptCFB(cipher, key, iv); err != nil {
			t.Error(err)
		} else {
			t.Logf("解密内容: %s\n", plain)
		}
	}
}

func TestCTR(t *testing.T) {
	data := "hello world"
	if cipher, err := EncryptCTR(data, key, iv); err != nil {
		t.Error(err)
	} else {
		t.Logf("加密内容: %s\n", cipher)
		// 解密
		if plain, err := DecryptCTR(cipher, key, iv); err != nil {
			t.Error(err)
		} else {
			t.Logf("解密内容: %s\n", plain)
		}
	}
}

func TestOFB(t *testing.T) {
	data := "hello world"
	if cipher, err := EncryptOFB(data, key, iv); err != nil {
		t.Error(err)
	} else {
		t.Logf("加密内容: %s\n", cipher)
		// 解密
		if plain, err := DecryptOFB(cipher, key, iv); err != nil {
			t.Error(err)
		} else {
			t.Logf("解密内容: %s\n", plain)
		}
	}
}
