package main

import (
	"crypto/aes"
	"crypto/cipher"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var (
	copyBuffer = make([]byte, 32<<10)

	aesKey   [32]byte
	cipherIV [aes.BlockSize]byte

	encryptOrDecrypt bool // true:encrypt, false:decrypt

	srcRoot string
	dstRoot string
)

func init() {
	flag.StringVar(&srcRoot, "src", "", "source root directory")
	flag.StringVar(&dstRoot, "dst", "", "destination root directory")
}

func getDstPath(srcPath, srcRoot, dstRoot string) (dstPath string, err error) {
	rel, err := filepath.Rel(srcRoot, srcPath)
	if err != nil {
		return
	}
	dstPath = filepath.Join(dstRoot, rel)
	return
}

func walk(path string, info os.FileInfo, inErr error) (err error) {
	if inErr != nil {
		return inErr
	}

	switch {
	case info.Mode().IsRegular():
		dstPath, err := getDstPath(path, srcRoot, dstRoot)
		if err != nil {
			return err
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, info.Mode().Perm())
		if err != nil {
			return err
		}
		defer dstFile.Close()

		block, err := aes.NewCipher(aesKey[:])
		if err != nil {
			return err
		}

		var stream cipher.Stream
		if encryptOrDecrypt {
			stream = cipher.NewCFBEncrypter(block, cipherIV[:])
		} else {
			stream = cipher.NewCFBDecrypter(block, cipherIV[:])
		}

		streamReader := cipher.StreamReader{S: stream, R: srcFile}

		_, err = io.CopyBuffer(dstFile, streamReader, copyBuffer)
		return err

	case info.IsDir():
		dstDir, err := getDstPath(path, srcRoot, dstRoot)
		if err != nil {
			return err
		}
		return os.MkdirAll(dstDir, 0666)

	default:
		return nil
	}
}

func main() {
	var key, iv string
	flag.StringVar(&key, "key", "", "aes key")
	flag.StringVar(&iv, "iv", "", "cipher iv")

	var doEncrypt, doDecrypt bool
	flag.BoolVar(&doEncrypt, "e", false, "do encrypt")
	flag.BoolVar(&doDecrypt, "d", false, "do decrypt")

	flag.Parse()

	copy(aesKey[:], key)
	copy(cipherIV[:], iv)

	switch {
	case doEncrypt:
		encryptOrDecrypt = true
	case doDecrypt:
		encryptOrDecrypt = false
	default:
		encryptOrDecrypt = true
	}

	srcRoot = filepath.Clean(srcRoot)
	dstRoot = filepath.Clean(dstRoot)

	if srcRoot == dstRoot {
		fmt.Println("src, dst 不能是同一个目录")
		return
	}

	srcFileInfo, err := os.Stat(srcRoot)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !srcFileInfo.IsDir() {
		fmt.Println("src 不是目录")
		return
	}

	dstFileInfo, err := os.Stat(dstRoot)
	if err != nil { // 可能是路径不正确, 也有可能是路径不存在
		if err := os.MkdirAll(dstRoot, 0666); err != nil {
			fmt.Println(err)
			return
		}
	} else {
		if !dstFileInfo.IsDir() {
			fmt.Println("dst 不是目录")
			return
		}
	}

	fmt.Println(filepath.Walk(srcRoot, walk))
}
