package utils

import (
	"encoding/binary"
	"io"
)

// DecodeTCPMsgWithLePrefix 编码：长度前缀 + 数据
func EncodeTCPMsgWithLePrefix(data []byte) ([]byte, error) {
	length := uint32(len(data))
	buf := make([]byte, 4+len(data))

	// 前4字节存储长度（大端序）
	binary.BigEndian.PutUint32(buf[0:4], length)

	// 复制数据
	copy(buf[4:], data)

	return buf, nil
}

// DecodeTCPMsgWithLePrefix 解码：读取长度前缀 + 数据
func DecodeTCPMsgWithLePrefix(reader io.Reader) ([]byte, error) {
	// 先读取4字节的长度前缀
	lengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(reader, lengthBuf); err != nil {
		return nil, err
	}

	// 解析长度
	length := binary.BigEndian.Uint32(lengthBuf)

	// 读取实际数据
	data := make([]byte, length)
	if _, err := io.ReadFull(reader, data); err != nil {
		return nil, err
	}

	return data, nil
}
