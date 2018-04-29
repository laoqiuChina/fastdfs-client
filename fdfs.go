package fdfs

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/wodog/pool"
)

type Client struct {
	timeout  time.Duration
	poolSize uint
	trackers []*tracker
}

func Default() *Client {
	return &Client{
		timeout:  10 * time.Second,
		poolSize: 10,
	}
}

func (c *Client) AddTracker(trackerAddr string) error {
	host, port, err := net.SplitHostPort(trackerAddr)
	if err != nil {
		return err
	}

	p, err := pool.NewDefault(func() (io.Closer, error) {
		return net.DialTimeout("tcp", fmt.Sprintf("%s:%s", host, port), c.timeout)
	})
	if err != nil {
		return err
	}

	t := &tracker{
		host: host,
		port: port,
		Pool: p,
	}

	c.trackers = append(c.trackers, t)
	return nil
}

func (c *Client) Upload(file io.Reader,ext string) (string, error) {
	tracker, err := c.getTracker()
	if err != nil {
		return "", err
	}
	storage, err := tracker.getUploadStorage()
	if err != nil {
		return "", err
	}
	fileID, err := storage.upload(file,ext)
	if err != nil {
		return "", err
	}
	return fileID, nil
}

func (c *Client) Download(fileID string, w io.Writer) error {
	tracker, err := c.getTracker()
	if err != nil {
		return err
	}
	storage, err := tracker.getDownloadStorage(fileID)
	if err != nil {
		return err
	}
	err = storage.download(fileID, w)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Delete(fileID string) error {
	tracker, err := c.getTracker()
	if err != nil {
		return err
	}
	storage, err := tracker.getUploadStorage()
	if err != nil {
		return err
	}
	err = storage.delete(fileID)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Info(fileID string) (map[string]string, error) {
	tracker, err := c.getTracker()
	if err != nil {
		return nil, err
	}
	storage, err := tracker.getUploadStorage()
	if err != nil {
		return nil, err
	}
	m, err := storage.info(fileID)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *Client) getTracker() (*tracker, error) {
	if len(c.trackers) == 0 {
		return nil, errors.New("tracker列表为空")
	}
	tracker := c.trackers[0]
	return tracker, nil
}
