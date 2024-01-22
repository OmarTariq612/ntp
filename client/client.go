package client

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"io"
	"net"
	"time"

	"github.com/OmarTariq612/ntp/proto"
)

type Client struct {
	conn net.Conn
}

func New() (Client, error) {
	raddr, err := net.ResolveUDPAddr("udp", proto.DefaultNTPServerAddr)
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
	Header *proto.Header
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
		r.Header = new(proto.Header)
	}

	r.Header.LIVNMode = proto.LIVNModeToByte(proto.ClockUnsynchronized, proto.V4, proto.Client)
	r.Header.Precision = -18
	if _, err := io.ReadFull(rand.Reader, sendArr[:8]); err != nil {
		return err
	}
	r.Header.TransmitTimestamp = proto.TimestampFormat(binary.LittleEndian.Uint64(sendArr[:8]))
	sendBuf := bytes.NewBuffer(sendArr[:0])
	if err := binary.Write(sendBuf, binary.BigEndian, r.Header); err != nil {
		return err
	}
	if _, err := c.conn.Write(sendBuf.Bytes()); err != nil {
		return err
	}
	org := proto.TimestampFormatFromUnixNano(time.Now().UnixNano())
	nread, err := c.conn.Read(recvArr[:])
	if err != nil {
		return err
	}
	dst := proto.TimestampFormatFromUnixNano(time.Now().UnixNano())
	recvBuf := bytes.NewBuffer(recvArr[:nread])
	if err := binary.Read(recvBuf, binary.BigEndian, r.Header); err != nil {
		return err
	}

	r.Offset = proto.Offset(org, r.Header.ReceiveTimestamp, r.Header.TransmitTimestamp, dst)

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
		r.Header = new(proto.Header)
	}
	// TODO
	return nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
