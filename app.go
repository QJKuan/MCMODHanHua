package main

import (
	"MCModHanHua/server"
	"MCModHanHua/server/translate"
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"strconv"
)

// SIGN 标记 查看是否正在操作 | 0 未有操作 | 1 正在操作
var SIGN int

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Shutdown(ctx context.Context) {
	translate.TransClose()
}

func (a *App) GetCtx() context.Context {
	return a.ctx
}

// SelectFile 选择文件方法 会记录选择文件的名称以及路径
func (a *App) SelectFile() {
	// 初始化 NeedTransResults
	server.NeedTransResults = make([]server.NeedTrans, 0)

	// 初始化 文件 以及 JAR管理类
	server.InitFile()
	server.JAR = server.JarManager{}

	// Wails 文件选择器
	jarPath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择 mod 文件",
		// JAR 过滤器 只能选择 jar 文件
		Filters: []runtime.FileFilter{
			{DisplayName: "选择一个 .jar 文件", Pattern: "*.jar"},
		},
	})
	// 选择文件时点击取消的情况
	if "" == jarPath {
		return
	}
	if err != nil {
		runtime.LogError(a.ctx, "选择器异常 : "+err.Error())
		server.MessageDialogHandle(a.ctx, "选择器异常 : "+err.Error(), "error")
		return
	}

	// 解析 jar 提取 en_us.lang 文件
	server.JAR.SetOldPath(jarPath)
	err = server.GetEnUsLang()
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		server.MessageDialogHandle(a.ctx, err.Error(), "error")
		return
	}

	server.MessageDialogHandle(a.ctx, "目标识别完成,可以开始汉化", "info")

}

// TranslateJar 翻译入口
func (a *App) TranslateJar() {
	// 查看是否已经有操作
	if 1 == SIGN {
		server.MessageDialogHandle(a.ctx, "正在处理上一步操作，先别搞我", "error")
		return
	}

	// 查看是否已经选择了文件
	if server.JAR.GetName() == "" {
		server.MessageDialogHandle(a.ctx, "还没选文件，先别急 \n去选文件先 !!! ", "error")
		return
	}

	// 翻译时间较长 添加一个是否翻译的交互
	total := strconv.Itoa(server.JAR.GetTotal())
	handle := server.MessageDialogHandle(a.ctx, "确定要现在翻译？\n共有 "+total+" 个条目需要翻译\n因为是实时翻译,所以时间可能会比较久\n平均一秒翻译一个", "question")
	if "No" == handle {
		server.MessageDialogHandle(a.ctx, "取消翻译成功", "info")
		return
	}

	SIGN = 1
	// 开始翻译
	err := server.DeepLTranslate(a.ctx)
	if err != nil {
		runtime.LogError(a.ctx, "汉化失败 : "+err.Error())
		server.MessageDialogHandle(a.ctx, "汉化失败 : "+err.Error(), "error")
		SIGN = 0
		return
	}
	server.MessageDialogHandle(a.ctx, "汉化完成", "info")
	SIGN = 0
}

// SaveJar 保存翻译后文件的入口
func (a *App) SaveJar() {
	if 1 == SIGN {
		server.MessageDialogHandle(a.ctx, "正在处理上一步操作，先别搞我", "error")
		return
	}

	newJarPath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultFilename: server.JAR.Name,
		Title:           "选择目录保存汉化后的 mod 文件",
	})
	// 选择文件时点击取消的情况
	if "" == newJarPath {
		return
	}
	SIGN = 1
	if err != nil {
		runtime.LogError(a.ctx, "保存文件异常 : "+err.Error())
		server.MessageDialogHandle(a.ctx, "保存文件异常 : "+err.Error(), "error")
		SIGN = 0
		return
	}

	server.JAR.SetNewPath(newJarPath)

	// 保存
	err = server.SaveJarInWin()
	if err != nil {
		runtime.LogError(a.ctx, "保存文件异常 : "+err.Error())
		server.MessageDialogHandle(a.ctx, "保存文件异常 : "+err.Error(), "error")
		os.Remove(server.JAR.GetNewPath())
		SIGN = 0
		return
	}
	server.MessageDialogHandle(a.ctx, "保存文件完成", "info")
	SIGN = 0
}
