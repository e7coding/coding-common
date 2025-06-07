// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package ddes 提供 DES 加密/解密 算法的实用接口
package jdes

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"github.com/e7coding/coding-common/errs/jerr"
)

const (
	NoPad    = iota // 不使用填充
	PKCS5PAD        // 使用 PKCS#5 填充
)

// EncECB 使用 ECB 模式对明文进行 DES 加密，pad 参数指定填充方式
func EncECB(plainText []byte, key []byte, pad int) ([]byte, error) {
	text, err := Pad(plainText, pad)
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, len(text))
	block, err := des.NewCipher(key)
	if err != nil {
		err = jerr.WithMsgErrF(err, `des.NewCipher 失败，key="%s"`, key)
		return nil, err
	}

	blockSize := block.BlockSize()
	for i, count := 0, len(text)/blockSize; i < count; i++ {
		begin, end := i*blockSize, i*blockSize+blockSize
		block.Encrypt(cipherText[begin:end], text[begin:end])
	}
	return cipherText, nil
}

// DecECB 使用 ECB 模式对密文进行 DES 解密，pad 参数指定填充方式
func DecECB(cipherText []byte, key []byte, pad int) ([]byte, error) {
	text := make([]byte, len(cipherText))
	block, err := des.NewCipher(key)
	if err != nil {
		err = jerr.WithMsgErrF(err, `des.NewCipher 失败，key="%s"`, key)
		return nil, err
	}

	blockSize := block.BlockSize()
	for i, count := 0, len(text)/blockSize; i < count; i++ {
		begin, end := i*blockSize, i*blockSize+blockSize
		block.Decrypt(text[begin:end], cipherText[begin:end])
	}

	plainText, err := UnPad(text, pad)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

// EncECBTriple 使用 TripleDES 的 ECB 模式加密，key 长度应为 16 或 24 字节
func EncECBTriple(plainText []byte, key []byte, pad int) ([]byte, error) {
	if len(key) != 16 && len(key) != 24 {
		return nil, jerr.WithMsg("key 长度错误")
	}

	text, err := Pad(plainText, pad)
	if err != nil {
		return nil, err
	}

	var newKey []byte
	if len(key) == 16 {
		newKey = append([]byte{}, key...)
		newKey = append(newKey, key[:8]...)
	} else {
		newKey = append([]byte{}, key...)
	}

	block, err := des.NewTripleDESCipher(newKey)
	if err != nil {
		err = jerr.WithMsgErrF(err, `des.NewTripleDESCipher 失败，key="%s"`, newKey)
		return nil, err
	}

	blockSize := block.BlockSize()
	cipherText := make([]byte, len(text))
	for i, count := 0, len(text)/blockSize; i < count; i++ {
		begin, end := i*blockSize, i*blockSize+blockSize
		block.Encrypt(cipherText[begin:end], text[begin:end])
	}
	return cipherText, nil
}

// DecECBTriple 使用 TripleDES 的 ECB 模式解密，key 长度应为 16 或 24 字节
func DecECBTriple(cipherText []byte, key []byte, pad int) ([]byte, error) {
	if len(key) != 16 && len(key) != 24 {
		return nil, jerr.WithMsg("key 长度错误")
	}

	var newKey []byte
	if len(key) == 16 {
		newKey = append([]byte{}, key...)
		newKey = append(newKey, key[:8]...)
	} else {
		newKey = append([]byte{}, key...)
	}

	block, err := des.NewTripleDESCipher(newKey)
	if err != nil {
		err = jerr.WithMsgErrF(err, `des.NewTripleDESCipher 失败，key="%s"`, newKey)
		return nil, err
	}

	blockSize := block.BlockSize()
	text := make([]byte, len(cipherText))
	for i, count := 0, len(text)/blockSize; i < count; i++ {
		begin, end := i*blockSize, i*blockSize+blockSize
		block.Decrypt(text[begin:end], cipherText[begin:end])
	}

	plainText, err := UnPad(text, pad)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

// EncCBC 使用 CBC 模式对明文进行 DES 加密，iv 为初始化向量，pad 指定填充方式
func EncCBC(plainText []byte, key []byte, iv []byte, pad int) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		err = jerr.WithMsgErrF(err, `des.NewCipher 失败，key="%s"`, key)
		return nil, err
	}

	if len(iv) != block.BlockSize() {
		return nil, jerr.WithMsg("无效的 iv 长度")
	}

	text, err := Pad(plainText, pad)
	if err != nil {
		return nil, err
	}
	cipherText := make([]byte, len(text))

	encryptor := cipher.NewCBCEncrypter(block, iv)
	encryptor.CryptBlocks(cipherText, text)

	return cipherText, nil
}

// DecCBC 使用 CBC 模式对密文进行 DES 解密，iv 为初始化向量，pad 指定填充方式
func DecCBC(cipherText []byte, key []byte, iv []byte, pad int) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		err = jerr.WithMsgErrF(err, `des.NewCipher 失败，key="%s"`, key)
		return nil, err
	}

	if len(iv) != block.BlockSize() {
		return nil, jerr.WithMsg("无效的 iv 长度")
	}

	text := make([]byte, len(cipherText))
	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypter.CryptBlocks(text, cipherText)

	plainText, err := UnPad(text, pad)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

// EncCBCTriple 使用 TripleDES 的 CBC 模式加密
func EncCBCTriple(plainText []byte, key []byte, iv []byte, pad int) ([]byte, error) {
	if len(key) != 16 && len(key) != 24 {
		return nil, jerr.WithMsg("key 长度错误")
	}

	var newKey []byte
	if len(key) == 16 {
		newKey = append([]byte{}, key...)
		newKey = append(newKey, key[:8]...)
	} else {
		newKey = append([]byte{}, key...)
	}

	block, err := des.NewTripleDESCipher(newKey)
	if err != nil {
		err = jerr.WithMsgErrF(err, `des.NewTripleDESCipher 失败，key="%s"`, newKey)
		return nil, err
	}

	if len(iv) != block.BlockSize() {
		return nil, jerr.WithMsg("无效的 iv 长度")
	}

	text, err := Pad(plainText, pad)
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, len(text))
	encrypter := cipher.NewCBCEncrypter(block, iv)
	encrypter.CryptBlocks(cipherText, text)

	return cipherText, nil
}

// DecCBCTriple 使用 TripleDES 的 CBC 模式解密
func DecCBCTriple(cipherText []byte, key []byte, iv []byte, pad int) ([]byte, error) {
	if len(key) != 16 && len(key) != 24 {
		return nil, jerr.WithMsg("key 长度错误")
	}

	var newKey []byte
	if len(key) == 16 {
		newKey = append([]byte{}, key...)
		newKey = append(newKey, key[:8]...)
	} else {
		newKey = append([]byte{}, key...)
	}

	block, err := des.NewTripleDESCipher(newKey)
	if err != nil {
		err = jerr.WithMsgErrF(err, `des.NewTripleDESCipher 失败，key="%s"`, newKey)
		return nil, err
	}

	if len(iv) != block.BlockSize() {
		return nil, jerr.WithMsg("无效的 iv 长度")
	}

	text := make([]byte, len(cipherText))
	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypter.CryptBlocks(text, cipherText)

	plainText, err := UnPad(text, pad)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

// PadPKCS5 对数据应用 PKCS#5 填充
func PadPKCS5(text []byte, blockSize int) []byte {
	pad := blockSize - len(text)%blockSize
	padText := bytes.Repeat([]byte{byte(pad)}, pad)
	return append(text, padText...)
}

// UnPadPKCS5 去除 PKCS#5 填充
func UnPadPKCS5(text []byte) []byte {
	length := len(text)
	pad := int(text[length-1])
	return text[:length-pad]
}

// Pad 根据指定填充方式对数据进行填充
func Pad(text []byte, pad int) ([]byte, error) {
	switch pad {
	case NoPad:
		if len(text)%8 != 0 {
			return nil, jerr.WithMsg("明文长度无效")
		}
	case PKCS5PAD:
		return PadPKCS5(text, 8), nil
	default:
		return nil, jerr.WithMsgF("不支持的填充类型 %d", pad)
	}
	return text, nil
}

// UnPad 根据指定填充方式去除填充
func UnPad(text []byte, pad int) ([]byte, error) {
	switch pad {
	case NoPad:
		if len(text)%8 != 0 {
			return nil, jerr.WithMsg("密文长度无效")
		}
	case PKCS5PAD:
		return UnPadPKCS5(text), nil
	default:
		return nil, jerr.WithMsgF("不支持的填充类型 %d", pad)
	}
	return text, nil
}
