package main

func main() {
	// var sentHeader ntp.Header
	// sentHeader.LIVNMode = ntp.LIVNModeToByte(ntp.ClockUnsynchronized, ntp.V4, ntp.Client)
	// sentHeader.Precision = -18
	// sentHeader.TransmitTimestamp = 555555

	// var sendBuf bytes.Buffer
	// if err := binary.Write(&sendBuf, binary.BigEndian, &sentHeader); err != nil {
	// 	panic(err)
	// }

	// raddr, err := net.ResolveUDPAddr("udp", ntp.DefaultNTPServer)
	// if err != nil {
	// 	panic(err)
	// }
	// conn, err := net.DialUDP("udp", nil, raddr)
	// if err != nil {
	// 	panic(err)
	// }
	// defer conn.Close()

	// var recvSlice [8196]byte

	// for {
	// 	org := ntp.TimestampFormatFromUnixNano(time.Now().UnixNano())

	// 	if _, err = conn.Write(sendBuf.Bytes()); err != nil {
	// 		panic(err)
	// 	}

	// 	nread, err := conn.Read(recvSlice[:])
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	dst := ntp.TimestampFormatFromUnixNano(time.Now().UnixNano())

	// 	recvBuf := bytes.NewBuffer(recvSlice[:nread])
	// 	var recvHeader ntp.Header
	// 	if err = binary.Read(recvBuf, binary.BigEndian, &recvHeader); err != nil {
	// 		panic(err)
	// 	}

	// 	offsetDuration := ntp.Offset(org, recvHeader.ReceiveTimestamp, recvHeader.TransmitTimestamp, dst)
	// 	if err := adjTime(offsetDuration); err != nil {
	// 		panic(err)
	// 	}

	// 	fmt.Println(offsetDuration)

	// 	if offsetDuration.Milliseconds() < 150 {
	// 		break
	// 	}

	// 	time.Sleep(2 * time.Second)
	// }
}
