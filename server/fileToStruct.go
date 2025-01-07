package server

import (
	"github.com/bytedance/sonic"
	"github.com/labstack/gommon/log"
	"io"
	"strings"
)

// PreTransJson 将文件转换陈 Json 数据
func PreTransJson(read io.Reader) {
	date, err := io.ReadAll(read)
	if err != nil {
		log.Error("读取数据流异常 : ", err)
	}

	err = sonic.Unmarshal(date, &NeedTransJson)
	if err != nil {
		log.Error("Json序列化异常 : ", err)
		return
	}

	for key, val := range NeedTransJson {
		NeedTransResults = append(NeedTransResults, NeedTrans{Key: key, Val: val})
	}

	//return &NeedTransResults
}

// SufTransJson 将翻译后的数据输入 map 中 转换成 Json 的 byte 数据
func SufTransJson() {
	for _, result := range NeedTransResults {
		NeedTransJson[result.Key] = result.Val
	}
	TransString, err := sonic.MarshalString(NeedTransJson)
	if err != nil {
		log.Error("将 Json 转换成 String 时异常 : ", err.Error())
	}
	fw, err := NewFileWriter("./File/zh_CN.json")
	err = fw.WriteLine(TransString)
	fw.Close()
}

// ParseLineLand 转换并拼接字符 旧版 land 包
func ParseLineLand(line string) *NeedTrans {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) < 2 {
		// 处理没有 "=" 的情况，这里选择将信息存入 key 中 val 为空
		return &NeedTrans{Key: line, Val: ""}
	}
	return &NeedTrans{Key: parts[0] + "=", Val: parts[1]}
}
