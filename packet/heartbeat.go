package packet

// NewHeartBeatPacket 构造心跳包
func NewHeartBeatPacket() []byte {
	pkt := NewPacket(Plain, HeartBeat, nil)
	return pkt.Build()
}
