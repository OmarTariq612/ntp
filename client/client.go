package client

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"io"
	"net"
	"time"

	"github.com/OmarTariq612/ntp/ntp"
	"golang.org/x/sys/unix"
)

type Client struct {
	conn net.Conn
}

func New() (Client, error) {
	raddr, err := net.ResolveUDPAddr("udp", ntp.DefaultNTPServer)
	if err != nil {
		return Client{}, nil
	}
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return Client{}, nil
	}
	return Client{conn: conn}, nil
}

func NewWithConn(conn net.Conn) Client {
	return Client{conn: conn}
}

type Response struct {
	Header *ntp.Header
	Offset time.Duration
}

var (
	sendArr [48]byte
	recvArr [1024]byte
)

func (c *Client) Query() (*Response, error) {
	var resp Response
	if err := c.query(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) query(r *Response) error {
	if r.Header == nil {
		r.Header = new(ntp.Header)
	}

	r.Header.LIVNMode = ntp.LIVNModeToByte(ntp.ClockUnsynchronized, ntp.V4, ntp.Client)
	r.Header.Precision = -18
	if _, err := io.ReadFull(rand.Reader, sendArr[:8]); err != nil {
		return err
	}
	r.Header.TransmitTimestamp = ntp.TimestampFormat(binary.LittleEndian.Uint64(sendArr[:8]))
	sendBuf := bytes.NewBuffer(sendArr[:0])
	if err := binary.Write(sendBuf, binary.BigEndian, r.Header); err != nil {
		return err
	}
	if _, err := c.conn.Write(sendBuf.Bytes()); err != nil {
		return err
	}
	org := ntp.TimestampFormatFromUnixNano(time.Now().UnixNano())
	nread, err := c.conn.Read(recvArr[:])
	if err != nil {
		return err
	}
	dst := ntp.TimestampFormatFromUnixNano(time.Now().UnixNano())
	recvBuf := bytes.NewBuffer(recvArr[:nread])
	if err := binary.Read(recvBuf, binary.BigEndian, r.Header); err != nil {
		return err
	}

	r.Offset = ntp.Offset(org, r.Header.ReceiveTimestamp, r.Header.TransmitTimestamp, dst)

	return nil
}

func (c *Client) AdjustTime() (*Response, error) {
	var resp Response
	if err := c.adjustTime(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) adjustTime(r *Response) error {
	if r.Header == nil {
		r.Header = new(ntp.Header)
	}
	// TODO
	return nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func adjTime(delta time.Duration) error {
	newTime := time.Now().Add(delta).UnixMicro()
	return unix.Settimeofday(&unix.Timeval{
		Sec:  newTime / 1_000_000,
		Usec: newTime % 1_000_000,
	})
}
