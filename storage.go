package fdfs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

type storage struct {
	host  string
	port  string
	group string
	index byte
}

func (s *storage) upload(file io.Reader,ext string) (string, error) {
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		return "", err
	}
	defer conn.Close()
	extbyte := []byte(ext)
	exttmp :=make([]byte, 6)
	copy(exttmp, extbyte)
	buf := &bytes.Buffer{}
	buf.WriteByte(s.index)
	buf.Write(lengthByte(uint64(len(b))))
	buf.Write(exttmp)
	buf.Write(b)
	h := header{
		uint64(buf.Len()),
		STORAGE_PROTO_CMD_UPLOAD_FILE,
		0,
	}
	p := newProtocol(h, conn)
	err = p.request(buf.Bytes())
	if err != nil {
		return "", err
	}
	b, err = p.body()
	if err != nil {
		return "", err
	}

	group := clearZero(string(b[:16]))
	path := clearZero(string(b[16:]))
	return group + "/" + path, nil
}

func (s *storage) download(fileID string, w io.Writer) error {
	ss := strings.SplitN(fileID, "/", 2)
	group := ss[0]
	path := ss[1]

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		return err
	}
	defer conn.Close()

	buf := &bytes.Buffer{}
	buf.Write(lengthByte(0))
	buf.Write(lengthByte(0))
	b := make([]byte, 16)
	copy(b, group)
	buf.Write(b)
	buf.WriteString(path)
	h := header{
		uint64(buf.Len()),
		STORAGE_PROTO_CMD_DOWNLOAD_FILE,
		0,
	}
	p := newProtocol(h, conn)
	err = p.request(buf.Bytes())
	if err != nil {
		return err
	}

	_, err = io.CopyN(w, p, int64(p.length))
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) delete(fileID string) error {
	ss := strings.SplitN(fileID, "/", 2)
	group := ss[0]
	path := ss[1]

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		return err
	}
	defer conn.Close()

	buf := &bytes.Buffer{}
	b := make([]byte, 16)
	copy(b, group)
	buf.Write(b)
	buf.WriteString(path)
	h := header{
		uint64(buf.Len()),
		STORAGE_PROTO_CMD_DELETE_FILE,
		0,
	}
	p := newProtocol(h, conn)
	err = p.request(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) info(fileID string) (map[string]string, error) {
	ss := strings.SplitN(fileID, "/", 2)
	group := ss[0]
	path := ss[1]

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	buf := &bytes.Buffer{}
	b := make([]byte, 16)
	copy(b, group)
	buf.Write(b)
	buf.WriteString(path)
	h := header{
		uint64(buf.Len()),
		STORAGE_PROTO_CMD_QUERY_FILE_INFO,
		0,
	}
	p := newProtocol(h, conn)
	err = p.request(buf.Bytes())
	if err != nil {
		return nil, err
	}
	b, err = p.body()
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"size":      strconv.FormatUint(binary.BigEndian.Uint64(b[:8]), 10),
		"timestamp": strconv.FormatUint(binary.BigEndian.Uint64(b[8:16]), 10),
		"crc32":     strconv.FormatUint(binary.BigEndian.Uint64(b[16:24]), 10),
		"ip":        clearZero(string(b[24:])),
	}, nil
}
