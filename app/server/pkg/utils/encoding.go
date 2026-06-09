package utils

import (
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
)

// DecodeToUTF8 验证并解析 UTF-8 编码数据
// 职责单一：仅验证数据是否为有效的 UTF-8 编码
// 若数据是有效 UTF-8 则原样返回，否则返回 nil 表示非 UTF-8 编码
// 适用于已确认是 UTF-8 的读取展示场景
func DecodeToUTF8(data []byte) []byte {
	if utf8.Valid(data) {
		return data
	}
	return nil
}

// TryDecodeToUTF8 检测并转换多种编码格式到 UTF-8
// 支持：GBK/GB2312 → UTF-8
// 检测策略：先尝试 GBK 解码，失败则回退检查 UTF-8 有效性
// 适用于文件上传、入库等需要处理多种编码的场景
func TryDecodeToUTF8(data []byte) []byte {
	if len(data) == 0 {
		return data
	}

	// 1. 先尝试 GBK 解码（GBK 解码器可处理 GB2312/GBK）
	decoder := simplifiedchinese.GBK.NewDecoder()
	utf8Data, err := decoder.Bytes(data)
	if err == nil {
		return utf8Data
	}

	// 2. GBK 解码失败 → 检查是否是有效 UTF-8
	if utf8.Valid(data) {
		return data
	}

	// 3. 都不是，返回原始数据（放弃转码）
	return data
}
