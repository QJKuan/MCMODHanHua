package translate

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Trans 存储结果
type Trans struct {
	Index     int
	Translate string
	Err       error
}

var Cancels []context.CancelFunc

// TranslateDeepL 使用DeepL翻译
func TranslateDeepL(urls map[int]string, appCtx context.Context) []Trans {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Headless,              // run headless mode
		chromedp.NoFirstRun,            // skip first run tasks
		chromedp.NoDefaultBrowserCheck, // disable default browser check
		chromedp.DisableGPU,            // disable GPU usage for lower resources
	)
	allocatorCtx, cancelAllocator := chromedp.NewExecAllocator(context.Background(), opts...)
	Cancels = append(Cancels, cancelAllocator)

	//ctx, cancel := chromedp.NewContext(allocatorCtx, chromedp.WithLogf(log.Printf), chromedp.WithDebugf(log.Printf))
	ctx, cancel := chromedp.NewContext(allocatorCtx)
	Cancels = append(Cancels, cancel)

	// 正常释放资源
	defer TransClose()

	// 记录开始时间
	//startTime := time.Now().Unix()

	// 定义要获取文本的节点信息
	nodeInfo := "/html[1]/body[1]/div[1]/div[1]/div[1]/div[3]/div[2]/div[1]/div[2]/div[1]/main[1]/div[2]/nav[1]/div[1]/div[2]/div[1]/div[1]/div[1]/div[1]/div[1]/div[1]/section[1]/div[1]/div[2]/div[3]/section[1]/div[1]/d-textarea[1]/div[1]/p[1]/span[1]"
	//nodeInfo := "//*[@id=\"textareasContainer\"]/div[3]/section/div[1]/d-textarea/div/p/span"

	var results []Trans
	// 遍历每个URL并执行任务
	runtime.EventsEmit(appCtx, "trans-msg", "正在和翻译服务建立连接.......")
	for i, url := range urls {
		var res Trans
		res.Index = i

		//TODO 处理现阶段无法翻译的字符
		if url == "" {
			// 判空返回
			res.Translate = ""
			results = append(results, res)
			continue
		}
		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.Text(nodeInfo, &res.Translate, chromedp.BySearch),
		)
		if err != nil {
			res.Err = err
			fmt.Printf("第 %v 行翻译异常: %v\n", i, err)
			continue
		}
		runtime.EventsEmit(appCtx, "trans-msg", fmt.Sprintf("翻译: %s\n", res.Translate))
		results = append(results, res)
	}
	runtime.EventsEmit(appCtx, "trans-msg", "翻译完成,点击保存可保存汉化后文件")
	return results
	// 计算并输出执行时间
	//fmt.Printf("执行时间: %d 秒\n", time.Now().Unix()-startTime)
}

// TransClose 释放资源
func TransClose() {
	for _, cancel := range Cancels {
		cancel()
	}
}
