package server

import (
	"MCModHanHua/server/translate"
	"MCModHanHua/server/youdaoyunAPI"
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"time"
)

func HanHuaServer(file io.Reader) {
	// 存放待翻译字符
	scanner := bufio.NewScanner(file)
	//表示新版使用的 json 存储的语言文件
	if VER == "json" {
		PreTransJson(file)
		//SufTransJson()
	} else {
		for scanner.Scan() {
			line := scanner.Text()
			//分开后的值的对象 key为 "=" 前的值 val为等于号后的值
			if entry := ParseLineLand(line); entry != nil {
				NeedTransResults = append(NeedTransResults, *entry)
			}
		}
		if scanner.Err() != nil {
			fmt.Printf("读取文件异常: %v", scanner.Err())
			return
		}
	}

	//DeepLTranslate(res)
}

func DeepLTranslate(ctx context.Context) error {
	// 需要转换的符号
	replacements := []string{
		". ", ",",
		",", " ",
		": ", ":",
		"! ", "!",
		"? ", "?",
		"%", "%25",
		"<br>", "*QJK*",
	}

	// 打开文件流
	fw, err := NewFileWriter("./File/zh_CN.lang")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer func() {
		if err = fw.Close(); err != nil {
			fmt.Printf("关闭文件失败: %v\n", err)
		}
	}()

	urlPre := "https://www.deepl.com/zh/translator#en/zh-hans/"
	urls := make(map[int]string)

	// 拼接所有待翻译的 url 存入数组中
	for i, re := range NeedTransResults {
		if re.Val == "" {
			continue
		}
		// 转换标点 因为存在标点符号会让标签变更导致拿不到翻译
		re.Val = strings.NewReplacer(replacements...).Replace(re.Val)

		// 拼接URL
		urls[i] = urlPre + re.Val
	}

	// 获取翻译结果集
	results := translate.TranslateDeepL(urls, ctx)

	for _, result := range results {
		i := result.Index
		if result.Err != nil {
			err = fw.WriteLine(NeedTransResults[i].Key + "此行报错！！！！！！！！！！！！！！！！！！！！！: " + result.Err.Error())
			if err != nil {
				fmt.Printf("第 %v 行的 %v 写入异常: %v\n", i, NeedTransResults[i].Key, err)
				return err
			}
			continue
		}

		if strings.Contains(result.Translate, "*QJK*") {
			result.Translate = strings.ReplaceAll(result.Translate, "*QJK*", "<br>")
		}

		// 修正翻译中的百分号问题，确保不会错误地将 "%%" 转换为 "%"
		if strings.Contains(result.Translate, "% ") {
			result.Translate = strings.ReplaceAll(result.Translate, "% ", "%% ")
			result.Translate = strings.ReplaceAll(result.Translate, "%%%", "%%")
		}

		if VER == "json" {
			NeedTransResults[i].Val = result.Translate
			continue
		}
		err = fw.WriteLine(NeedTransResults[i].Key + result.Translate)
		if err != nil {
			fmt.Printf("第 %v 行的 %v 写入异常: %v\n", i, NeedTransResults[i].Key, err)
			return err
		}
	}
	if VER == "json" {
		SufTransJson()
	}
	return nil
}

// YoudaoTranslate 使用有道云API翻译 需要对应的 appKey 和 appSecret
func YoudaoTranslate(res []NeedTrans) {
	filePath := "./File/zh_CN.lang"
	fw, err := NewFileWriter(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(fw *FileWriter) {
		if err != nil {
			fmt.Printf("关闭文件失败: %v\n", err)
		}
	}(fw)

	//遍历并且输出参数
	for _, entry := range res {
		if entry.Val == "" {
			if err = fw.WriteLine(entry.Key); err != nil {
				fmt.Printf("写入文件异常: %v\n", err)
			}
			fmt.Println(entry.Key)
			continue
		}
		result := youdaoyunAPI.TransYouDaoYun(entry.Val)
		if err = fw.WriteLine(entry.Key + result); err != nil {
			fmt.Printf("写入文件异常: %v\n", err)
		}
		fmt.Println(entry.Key + result)

		// 根据API限制调整睡眠时间 这里用 2 秒翻译一次
		time.Sleep(2000 * time.Millisecond)
	}
}
