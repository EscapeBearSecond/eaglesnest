package utils

import (
	"archive/zip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 解压
func Unzip(zipFile string, destDir string) ([]string, error) {
	zipReader, err := zip.OpenReader(zipFile)
	var paths []string
	if err != nil {
		return []string{}, err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		if strings.Index(f.Name, "..") > -1 {
			return []string{}, fmt.Errorf("%s 文件名不合法", f.Name)
		}
		fpath := filepath.Join(destDir, f.Name)
		paths = append(paths, fpath)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return []string{}, err
			}

			inFile, err := f.Open()
			if err != nil {
				return []string{}, err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return []string{}, err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return []string{}, err
			}
		}
	}
	return paths, nil
}

func DecryptFile(data []byte, password, outputZipFile string) error {
	// 提取盐值、IV 和加密数据
	salt := data[:16]
	iv := data[16:32]
	encryptedContent := data[32:]

	// 使用 PBKDF2 生成密钥
	key := pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)

	// 创建 AES 解密器
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decryptedData := make([]byte, len(encryptedContent))
	mode.CryptBlocks(decryptedData, encryptedContent)

	// 去除填充
	decryptedData = unpad(decryptedData)

	// 将解密后的数据写入 zip 文件
	err = os.WriteFile(outputZipFile, decryptedData, 0644)
	if err != nil {
		return err
	}

	return nil
}

// 移除 PKCS7 填充
func unpad(data []byte) []byte {
	paddingLen := int(data[len(data)-1])
	return data[:len(data)-paddingLen]
}
