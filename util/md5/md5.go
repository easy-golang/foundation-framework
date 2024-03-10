package md5

import (
	"crypto/md5"
	"encoding/hex"
)

func Encode(str string) string {
	// 创建一个 MD5 散列对象
	hasher := md5.New()
	// 写入要计算散列的数据
	hasher.Write([]byte(str))
	// 计算 MD5 散列值
	hashBytes := hasher.Sum(nil)
	// 将散列值转换为十六进制字符串
	return hex.EncodeToString(hashBytes)
}
