package server

import (
	"archive/zip"
	"errors"
	"github.com/labstack/gommon/log"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// GetEnUsLang 获取英文的语言包
func GetEnUsLang() error {
	// 读取 Jar
	jarPath := JAR.GetOldPath()
	oldJar, err := zip.OpenReader(jarPath)
	if err != nil {
		log.Error("文件打开失败，格式不正确：", err.Error())
		return errors.New("文件打开失败，格式不正确")
	}
	defer oldJar.Close()

	// 判断是否存在 可汉化包 或 是否已经汉化
	EnUs := 0
	// 初始话 VER 防止后续存在问题
	VER = "lang"
	// 存放需要翻译的文本
	for _, file := range oldJar.File {
		if strings.Contains(strings.ToLower(file.Name), "zh_cn") {
			return errors.New("内置已有汉化")
		}
		if strings.Contains(strings.ToLower(file.Name), "en_us") {
			if strings.HasSuffix(file.Name, ".json") {
				//高版本翻译文件使用 json 存储
				VER = "json"
			}
			// 如果存在 en_us 则直接存在需要汉化列表
			EnUsFile, _ := file.Open()
			HanHuaServer(EnUsFile)
			EnUsFile.Close()
			EnUs = 1
		}
	}
	if EnUs == 0 {
		return errors.New("没有找到英文版的语言包")
	}
	return nil
}

// SaveJarInWin 把 Jar 保存到本地磁盘中
func SaveJarInWin() error {
	oldJar, err := zip.OpenReader(JAR.GetOldPath())
	defer oldJar.Close()

	// 创建一个新的jar文件用于写入
	newJar, err := os.Create(JAR.GetNewPath())
	if err != nil {
		log.Error("无法创建新的 Jar 文件：", err.Error())
		return errors.New("无法创建新的 Jar 文件：" + err.Error())
	}
	defer newJar.Close()

	jarWriter := zip.NewWriter(newJar)
	defer jarWriter.Close()

	tranPath := ""
	// 复制原有的文件到新的zip文件中，并添加新文件
	for _, file := range oldJar.File {
		fileReader, err1 := file.Open()
		if err1 != nil {
			log.Error("无法打开文件：", file.Name, err1.Error())
			continue
		}

		// TODO 配置压缩方式  Store:不压缩 | Deflate:压缩
		/*		var header *zip.FileHeader
				if strings.Contains(file.Name, "BOOT-INF") {
					header = &zip.FileHeader{
						Name:   file.Name,
						Method: zip.Store,
					}
				} else {
					header = &zip.FileHeader{
						Name:   file.Name,
						Method: zip.Deflate,
					}
				}*/
		header := &zip.FileHeader{
			Name:   file.Name,
			Method: zip.Deflate,
		}

		header.Modified = file.Modified
		header.SetMode(file.Mode())

		writer, err1 := jarWriter.CreateHeader(header)
		if err1 != nil {
			fileReader.Close()
			log.Error("无法创建zip头：", err1.Error())
			continue
		}

		_, err = io.Copy(writer, fileReader)
		fileReader.Close()
		if err != nil {
			log.Error("无法复制文件到新的zip：", file.Name, err.Error())
			continue
		}

		if strings.Contains(strings.ToLower(file.Name), "en_us") {
			tranPath = filepath.Dir(file.Name) + "/"
		}
	}

	// 添加新文件 zh_CN 到 lang 目录
	var transFile *os.File
	if "json" == VER {
		transFile, err = os.Open("./File/zh_CN.json")
	} else {
		transFile, err = os.Open("./File/zh_CN.lang")
	}
	if err != nil {
		log.Error("无法打开文件 zh_CN 文件 : ", err)
		return errors.New("无法打开文件 zh_CN 文件 : " + err.Error())
	}
	defer transFile.Close()

	var addWrite io.Writer
	if "json" == VER {
		addWrite, err = jarWriter.Create(tranPath + "zh_CN.json")
	} else {
		addWrite, err = jarWriter.Create(tranPath + "zh_CN.lang")
	}
	if err != nil {
		log.Error("无法在 zip 中创建文件 ：", err.Error())
		return errors.New("无法在 zip 中创建文件 ：" + err.Error())
	}

	_, err = io.Copy(addWrite, transFile)
	if err != nil {
		log.Error("无法将文件 zh_CN 复制到zip中：", err.Error())
		return errors.New("无法将文件 zh_CN 复制到 zip 中 ：" + err.Error())
	}
	return nil
}
