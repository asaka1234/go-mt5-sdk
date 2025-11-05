package pumping

import (
	"bytes"
	"fmt"
	"github.com/golang/snappy"
	"github.com/valyala/bytebufferpool"
	"github.com/vmihailenco/msgpack/v5"
)

// EncodeMsgpackSnappy 独立的编码函数（使用对象池）
func Encode[T any](v T) ([]byte, error) {
	// 从池中获取buffer
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf) // 归还给池

	// 配置编码器
	enc := msgpack.NewEncoder(buf)
	enc.SetCustomStructTag("msgpack")
	enc.UseCompactInts(true)
	enc.UseCompactFloats(true)

	// 序列化
	if err := enc.Encode(v); err != nil {
		return nil, err
	}

	// 压缩
	return snappy.Encode(nil, buf.Bytes()), nil
}

// DecodeMsgpackSnappy 独立的解码函数（使用对象池）
func Decode[T any](data []byte, v *T) error {
	// 解压
	decompressed, err := snappy.Decode(nil, data)
	if err != nil {
		return err
	}

	// 反序列化
	dec := msgpack.NewDecoder(bytes.NewReader(decompressed))
	dec.SetCustomStructTag("msgpack")
	// 解码器没有 UseCompactInts 和 UseCompactFloats 方法
	// 解码器会自动处理所有兼容的格式

	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("msgpack decode failed: %w", err)
	}

	return nil
}
