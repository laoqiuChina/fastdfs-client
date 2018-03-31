package fdfs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type tracker struct {
	host string
	port string
}

func (t tracker) getUploadStorage() (*storage, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", t.host, t.port), timeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	h := header{
		0,
		TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ONE,
		0,
	}
	p := newProtocol(h, conn)
	b, err := p.request(nil)
	if err != nil {
		return nil, err
	}

	group := clearZero(string(b[:16]))
	host := clearZero(string(b[16 : 16+15]))
	port := strconv.Itoa(int(binary.BigEndian.Uint64(b[16+15 : 16+15+8])))
	index := b[16+15+8 : 16+15+8+1][0]

	return &storage{
		group: group,
		host:  host,
		port:  port,
		index: index,
	}, nil
}

func (t tracker) getDownloadStorage(fileID string) (*storage, error) {
	ss := strings.SplitN(fileID, "/", 2)
	group := ss[0]
	path := ss[1]

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", t.host, t.port), timeout)
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
		TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ONE,
		0,
	}
	p := newProtocol(h, conn)
	b, err = p.request(buf.Bytes())
	if err != nil {
		return nil, err
	}

	group = clearZero(string(b[:16]))
	host := clearZero(string(b[16 : 16+15]))
	port := strconv.FormatUint(binary.BigEndian.Uint64(b[16+15:16+15+8]), 10)

	return &storage{
		group: group,
		host:  host,
		port:  port,
	}, nil
}
