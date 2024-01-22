package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/OmarTariq612/ntp/adjtime"
	"github.com/OmarTariq612/ntp/proto"
)

func main() {
	var sentHeader proto.Header
	sentHeader.LIVNMode = proto.LIVNModeToByte(proto.ClockUnsynchronized, proto.V4, proto.Client)
	sentHeader.Precision = -18
	sentHeader.TransmitTimestamp = 555555

	var sendBuf bytes.Buffer
	if err := binary.Write(&sendBuf, binary.BigEndian, &sentHeader); err != nil {
		panic(err)
	}

	raddr, err := net.ResolveUDPAddr("udp", proto.DefaultNTPServerAddr)
	if err != nil {
		panic(err)
	}
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	var recvSlice [8196]byte

	for {
		org := proto.TimestampFormatFromUnixNano(time.Now().UnixNano())

		if _, err = conn.Write(sendBuf.Bytes()); err != nil {
			fmt.Println("one")
			panic(err)
		}

		nread, err := conn.Read(recvSlice[:])
		if err != nil {
			fmt.Println("two")
			panic(err)
		}

		dst := proto.TimestampFormatFromUnixNano(time.Now().UnixNano())

		recvBuf := bytes.NewBuffer(recvSlice[:nread])
		var recvHeader proto.Header
		if err = binary.Read(recvBuf, binary.BigEndian, &recvHeader); err != nil {
			fmt.Println("three")
			panic(err)
		}

		offsetDuration := proto.Offset(org, recvHeader.ReceiveTimestamp, recvHeader.TransmitTimestamp, dst)
		if err := adjtime.AdjTime(offsetDuration); err != nil {
			fmt.Println("four")
			panic(err)
		}

		fmt.Println(offsetDuration)

		if offsetDuration.Milliseconds() < 150 {
			break
		}

		time.Sleep(2 * time.Second)
	}
}
