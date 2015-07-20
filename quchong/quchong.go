package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
)

type FileSum struct {
	size     int64
	crc32Sum uint32
	md5Sum   [16]byte
	sha1Sum  [20]byte
}

var (
	// 顺序执行, 下面变量无需加锁
	copyBuffer = make([]byte, 32<<10)

	crc32Hash = crc32.NewIEEE()
	md5Hash   = md5.New()
	sha1Hash  = sha1.New()

	FileSumSet = make(map[FileSum]string) // map[FileSum]path
)

func walk(path string, info os.FileInfo, inErr error) (err error) {
	if inErr != nil {
		return inErr
	}
	if !info.Mode().IsRegular() {
		return
	}

	shouldRemove := false

	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer func() {
		file.Close()
		if shouldRemove {
			err = os.Remove(path)
		}
	}()

	crc32Hash.Reset()
	md5Hash.Reset()
	sha1Hash.Reset()

	md5TeeReader := io.TeeReader(file, md5Hash)
	sha1TeeReader := io.TeeReader(md5TeeReader, sha1Hash)
	io.CopyBuffer(crc32Hash, sha1TeeReader, copyBuffer)

	var key FileSum
	key.size = info.Size()
	key.crc32Sum = crc32Hash.Sum32()
	copy(key.md5Sum[:], md5Hash.Sum(nil))
	copy(key.sha1Sum[:], sha1Hash.Sum(nil))

	fmt.Printf("path: %s\r\nsize: %d bytes\r\ncrc32: %x\r\nmd5: %s\r\nsha1: %s\r\n\r\n",
		path,
		key.size,
		key.crc32Sum,
		hex.EncodeToString(key.md5Sum[:]),
		hex.EncodeToString(key.sha1Sum[:]),
	)

	if pathx, ok := FileSumSet[key]; ok {
		shouldRemove = true
		fmt.Printf("%s 与 %s 重复, 将被移除\r\n", path, pathx)
	} else {
		FileSumSet[key] = path
	}

	return
}

func main() {
	var root string
	flag.StringVar(&root, "d", "", "root directory")

	flag.Parse()

	if root == "" {
		fmt.Println("请指定 d 参数")
		return
	}

	fmt.Println(filepath.Walk(root, walk))
}
