package server

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// MessageDialogHandle 文件信息处理
func MessageDialogHandle(ctx context.Context, msg string, level string) string {
	var err error
	// TODO WIN 环境不允许自定义按钮
	switch level {
	case "info":
		_, err = runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
			Type:          runtime.InfoDialog,
			Title:         "很对",
			Message:       msg,
			DefaultButton: "Yes",
		})
	case "error":
		_, err = runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
			Type:          runtime.ErrorDialog,
			Title:         "不对",
			Message:       msg,
			DefaultButton: "Yes",
		})
	case "warning":
		_, err = runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
			Type:          runtime.WarningDialog,
			Title:         "不太对",
			Message:       msg,
			DefaultButton: "Yes",
		})
	case "question":
		var ret string
		ret, err = runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
			Type:          runtime.QuestionDialog,
			Title:         "有疑问",
			Message:       msg,
			DefaultButton: "Yes",
		})
		return ret
	default:
		runtime.LogError(ctx, "参数 level 输入有误")
	}

	if err != nil {
		runtime.LogError(ctx, "对话框异常 : "+err.Error())
	}
	return ""
}
