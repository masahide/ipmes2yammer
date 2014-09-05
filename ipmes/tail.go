package ipmes

import "os"

type TailFile struct {
	lastSize int64
	file     *os.File
}

func NewTailFile(name string) (*TailFile, error) {
	t := TailFile{}
	var err error
	t.file, err = os.Open(name)
	if err != nil {
		return nil, err
	}
	fi, err := t.file.Stat()
	if err != nil {
		return nil, err
	}
	t.lastSize = fi.Size()

	t.file.Seek(t.lastSize, 0)
	return &t, nil

}
func (t *TailFile) TailMessage() (string, error) {
	fi, err := t.file.Stat()
	if err != nil {
		return "", err
	}
	if t.lastSize < fi.Size() {
		b := make([]byte, 4096)
		n, err := t.file.Read(b)
		if err != nil {
			return "", err
		}
		t.lastSize += int64(n)
		return string(b[0:n]), nil

	}
	return "", nil
}

func (t *TailFile) Close() {
	t.file.Close()
}
