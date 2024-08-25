package mstore

import (
	"bufio"
	"os"
	"path"
	"strconv"
	"strings"
)

type filePosition struct {
	currentPos uint64
	fileDir    string
	fileName   string
}

func NewFilePosition(path, name string) BlockPosition {
	bp := new(filePosition)
	bp.fileDir = path
	bp.fileName = name
	return bp
}

func (fp *filePosition) Update(n uint64) error {
	for {
		p := strconv.FormatUint(n, 10)
		f, err := os.OpenFile(path.Join(fp.fileDir, fp.fileName), os.O_WRONLY, 0757)
		if err != nil {
			return err
		}
		_, err = f.Seek(0, 0)
		if err != nil {
			return err
		}
		w := bufio.NewWriter(f)
		_, err = w.Write([]byte(p))
		if err != nil {
			return err
		}
		err = w.Flush()
		if err != nil {
			return err
		}
		err = f.Close()
		if err != nil {
			return err
		}
		return nil
	}
}

func (fp *filePosition) Exists() bool {
	_, err := os.Open(path.Join(fp.fileDir, fp.fileName))
	if err != nil {
		return false
	}
	return true
}

func (fp *filePosition) Create() error {
	if _, err := os.Stat(fp.fileDir); err != nil {
		err = os.MkdirAll(fp.fileDir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	fil, err := os.Create(path.Join(fp.fileDir, fp.fileName))
	if err != nil {
		return err
	}
	return fil.Close()
}

func (fp *filePosition) Read() (uint64, error) {
	fil, err := os.Open(path.Join(fp.fileDir, fp.fileName))
	if err != nil {
		return 0, err
	}
	scanner := bufio.NewScanner(fil)
	var block string
	for scanner.Scan() {
		block = strings.Trim(scanner.Text(), "\n")
	}
	blockI, err := strconv.ParseUint(block, 0, 64)
	if err != nil {
		return 0, err
	}
	return blockI, nil
}
