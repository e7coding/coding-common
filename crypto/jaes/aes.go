// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package daes 提供 AES 加解密算法的便捷 API。
package jaes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/e7coding/coding-common/errs/jerr"
	"io"
)

const (
	// IVDefaultValue 默认的初始化向量（IV）值。
	IVDefaultValue = "Root1234"
)

// Enc 是 EncCBC 的别名，用于 CBC 模式加密。
func Enc(plainText []byte, key []byte, iv ...[]byte) ([]byte, error) {
	return EncCBC(plainText, key, iv...)
}

// Dec 是 DecCBC 的别名，用于 CBC 模式解密。
func Dec(cipherText []byte, key []byte, iv ...[]byte) ([]byte, error) {
	return DecCBC(cipherText, key, iv...)
}

// EncCBC 使用 CBC 模式对明文进行加密，支持可选 IV 参数。
func EncCBC(plainText []byte, key []byte, iv ...[]byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		err = jerr.WithMsgErrF(err, `aes.NewCipher 失败，key="%s"`, key)
		return nil, err
	}
	blockSize := block.BlockSize()
	plainText = PKCS7Pad(plainText, blockSize)
	var ivValue []byte
	if len(iv) > 0 {
		ivValue = iv[0]
	} else {
		ivValue = []byte(IVDefaultValue)
	}
	mode := cipher.NewCBCEncrypter(block, ivValue)
	cipherText := make([]byte, len(plainText))
	mode.CryptBlocks(cipherText, plainText)
	return cipherText, nil
}

// DecCBC 使用 CBC 模式对密文进行解密，支持可选 IV 参数。
func DecCBC(cipherText []byte, key []byte, iv ...[]byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		err = jerr.WithMsgErrF(err, `aes.NewCipher 失败，key="%s"`, key)
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(cipherText) < blockSize {
		return nil, jerr.WithMsg("cipherText 太短")
	}
	var ivValue []byte
	if len(iv) > 0 {
		ivValue = iv[0]
	} else {
		ivValue = []byte(IVDefaultValue)
	}
	if len(cipherText)%blockSize != 0 {
		return nil, jerr.WithMsg("cipherText 长度不是块大小的整数倍")
	}
	mode := cipher.NewCBCDecrypter(block, ivValue)
	plainText := make([]byte, len(cipherText))
	mode.CryptBlocks(plainText, cipherText)
	plainText, err = PKCS7UnPad(plainText, blockSize)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

// PKCS5Pad 使用 PKCS#5 填充对数据进行对齐。
func PKCS5Pad(src []byte, blockSize ...int) []byte {
	size := 8
	if len(blockSize) > 0 {
		size = blockSize[0]
	}
	return PKCS7Pad(src, size)
}

// PKCS5UnPad 移除 PKCS#5 填充。
func PKCS5UnPad(src []byte, blockSize ...int) ([]byte, error) {
	size := 8
	if len(blockSize) > 0 {
		size = blockSize[0]
	}
	return PKCS7UnPad(src, size)
}

// PKCS7Pad 使用 PKCS#7 填充对数据进行对齐。
func PKCS7Pad(src []byte, blockSize int) []byte {
	pad := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(pad)}, pad)
	return append(src, padtext...)
}

// PKCS7UnPad 移除 PKCS#7 填充并验证合法性。
func PKCS7UnPad(src []byte, blockSize int) ([]byte, error) {
	length := len(src)
	if blockSize <= 0 {
		return nil, jerr.WithMsg(fmt.Sprintf("无效块大小：%d", blockSize))
	}
	if length == 0 || length%blockSize != 0 {
		return nil, jerr.WithMsg("无效数据长度")
	}
	unpad := int(src[length-1])
	if unpad == 0 || unpad > blockSize {
		return nil, jerr.WithMsg("无效填充大小")
	}
	for _, v := range src[length-unpad:] {
		if int(v) != unpad {
			return nil, jerr.WithMsg("填充不匹配")
		}
	}
	return src[:length-unpad], nil
}

// EncGCM 使用 AES-GCM 对明文进行加密，
// 支持可选的 nonce（IV）；如果未提供，则自动随机生成。
// 返回值格式：[nonce｜ciphertext]。
func EncGCM(plainText, key []byte, nonceOpt ...[]byte) ([]byte, error) {
	// 1. 创建 AES 块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, jerr.WithMsgErrF(err, `aes.NewCipher 失败，key="%x"`, key)
	}
	// 2. 创建 GCM AEAD
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, jerr.WithMsgErrF(err, `cipher.NewGCM 失败`)
	}
	// 3. 准备 nonce
	var nonce []byte
	if len(nonceOpt) > 0 {
		nonce = nonceOpt[0]
		if len(nonce) != aead.NonceSize() {
			return nil, jerr.WithMsg(`提供的 nonce 长度不符 GCM.NonceSize()`)
		}
	} else {
		nonce = make([]byte, aead.NonceSize())
		if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
			return nil, jerr.WithMsgErrF(err, `生成随机 nonce 失败`)
		}
	}
	// 4. Seal：在 nonce 之后拼接密文，AEAD 会把认证标签一并附加
	cipherText := aead.Seal(nil, nonce, plainText, nil)
	// 5. 返回 [nonce｜ciphertext]
	return append(nonce, cipherText...), nil
}

// DecCFB 使用 CFB 模式对密文进行解密。
// key 长度必须为 16/24/32 字节，iv 可选，若不提供则使用默认 IV。
// 参数 unPad 表示用于去除尾部零填充的字节数。
func DecCFB(cipherText []byte, key []byte, unPad int, iv ...[]byte) ([]byte, error) {
	// 创建 AES 块加密器
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, jerr.WithMsgErrF(
			err,
			`aes.NewCipher 失败，key 长度=%d`,
			len(key),
		)
	}
	blockSize := block.BlockSize()
	// 密文长度检查
	if len(cipherText) < blockSize {
		return nil, jerr.WithMsg("cipherText 长度不足")
	}
	// 选择 IV 并校验长度
	var ivValue []byte
	if len(iv) > 0 {
		ivValue = iv[0]
		if len(ivValue) != blockSize {
			return nil, jerr.WithMsg(fmt.Sprintf("IV 长度必须为 %d 字节，当前 %d 字节", blockSize, len(ivValue)))
		}
	} else {
		ivValue = []byte(IVDefaultValue)
	}
	// 使用 CFB 解密流
	stream := cipher.NewCFBDecrypter(block, ivValue)
	plainText := make([]byte, len(cipherText))
	stream.XORKeyStream(plainText, cipherText)
	// 去除尾部零填充
	if unPad > 0 && unPad <= len(plainText) {
		plainText = plainText[:len(plainText)-unPad]
	}
	return plainText, nil
}

// ZeroPad 对数据使用零字节进行填充，并返回填充长度。
func ZeroPad(cipherText []byte, blockSize int) ([]byte, int) {
	pad := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{0}, pad)
	return append(cipherText, padText...), pad
}

// ZeroUnPad 移除零字节填充。
func ZeroUnPad(plaintext []byte, unPadding int) []byte {
	return plaintext[:len(plaintext)-unPadding]
}
