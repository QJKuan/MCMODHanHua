package server

import "path/filepath"

// JAR 是单例 用来管理上传的文件以及翻译后的文件
var JAR JarManager

// VER 表示新旧版的全局变量
var VER string

// NeedTransJson 需要翻译的英文文本存储位置 | 新版的 Json 格式语言文件
var NeedTransJson map[string]string

// NeedTrans 是对需要翻译的对象的抽象结果集 | 不管是新旧版 后续都会把值存入这里
type NeedTrans struct {
	Key string
	Val string
}

// NeedTransResults 全局变量 需要翻译集合
var NeedTransResults []NeedTrans

type JarManager struct {
	OldPath string
	NewPath string
	Name    string
}

func (jar *JarManager) GetName() string {
	return jar.Name
}

func (jar *JarManager) GetOldPath() string {
	return jar.OldPath
}

func (jar *JarManager) GetNewPath() string {
	return jar.NewPath
}

func (jar *JarManager) SetOldPath(path string) {
	jar.OldPath = path
	jar.Name = filepath.Base(path)
	return
}

func (jar *JarManager) SetNewPath(path string) {
	jar.NewPath = path
	return
}

func (jar *JarManager) GetTotal() int {
	return len(NeedTransResults)
}
