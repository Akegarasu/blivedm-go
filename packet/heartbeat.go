package packet

// NewHeartBeatPacket 构造心跳包
func NewHeartBeatPacket() []byte {
	pkt := NewPacket(1, HeartBeat, nil)
	return pkt.Build()
}
