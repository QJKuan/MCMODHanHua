<script lang="ts" setup>
import { SelectFile, TranslateJar, SaveJar } from '../../wailsjs/go/main/App';
import {onMounted, onUnmounted, ref} from "vue";
import {EventsOff, EventsOn} from "../../wailsjs/runtime";

const selectAndAnalyzeFile = async () => {
  await SelectFile();
};

const translateFile = async () => {
  await TranslateJar();
};

const saveJar = async () => {
  await SaveJar();
};

// 内嵌页面滚动逻辑
const output = ref<HTMLElement | null>(null);
const scrollToBottom = () => {
  if (output.value) {
    output.value.scrollTop = output.value.scrollHeight;
  }
};
// 监听
const messages = ref<string[]>([]);
onMounted(() => {
  // 订阅后端事件
  EventsOn("trans-msg", (msg: string) => {
    messages.value.push(msg); // 将消息添加到messages数组中
    // console.log(msg);
    scrollToBottom(); // 滚动到底部
  });

  // 清理事件监听器
  onUnmounted(() => {
    EventsOff("trans-msg");
  });
});
</script>

<template>
  <el-container class="app-container">
    <el-header class="app-header">
      <h1 class="app-title">MOD 汉化工具</h1>
    </el-header>
    <el-main class="app-main">
      <el-row :gutter="20">
        <!-- 按钮区域 -->
        <el-col :xs="24" :sm="12" :md="8" :lg="6" :xl="4">
          <el-space direction="vertical" :size="20">
            <el-button
                type="primary"
                size="large"
                @click="selectAndAnalyzeFile"
                style="width: 100%;">
              选择 MOD 文件
            </el-button>
            <el-button
                type="success"
                size="large"
                @click="translateFile"
                style="width: 100%;">
              开始汉化
            </el-button>
            <el-button
                type="warning"
                size="large"
                @click="saveJar"
                style="width: 100%;">
              保存
            </el-button>
          </el-space>
        </el-col>
        <!-- 输出区域 -->
        <el-col :xs="24" :sm="12" :md="16" :lg="18" :xl="20">
          <div ref="output" class="output-box">
            <p v-for="msg in messages" :key="msg">{{ msg }}</p>
          </div>
        </el-col>
      </el-row>
    </el-main>
  </el-container>
</template>

<style scoped>
/* 输出翻译框样式 */
.output-box {
  border: 1px solid #ccc;
  background: #1b2636;
  padding: 10px;
  height: 180px; /* 根据需要调整高度 */
  width: 180px;
  overflow-y: auto;
  margin-top: 0;  /* 离顶部位置 */
  margin-left: 80px;  /* 离开左侧位置 */
  font-size: 10px;
  text-align: left;
}

/* 定义应用容器的样式 */
.app-container {
  height: 100vh; /* 高度占满整个视口 */
  background-image: url('../assets/images/mcbg.jpg'); /* 设置背景图片 */
  background-size: cover; /* 背景图片覆盖整个容器 */
  background-repeat: no-repeat; /* 背景图片不重复 */
  background-position: center; /* 背景图片居中 */
  position: relative; /* 相对定位，用于容纳绝对定位的子元素 */
  overflow: hidden; /* 隐藏溢出内容 */
}

/* 定义应用头部的样式 */
.app-header {
  display: flex; /* 使用Flex布局 */
  align-items: center; /* �垂直居中对齐 */
  justify-content: center; /* 水平居中对齐 */
  min-height: 200px; /* 最小高度为200px */
  font-size: 30px; /* 字体大小 */
}

/* 定义应用标题的样式 */
.app-title {
  color: #dcdbdb; /* 文字颜色 */
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5); /* 文字阴影效果 */
}

/* 定义应用主体的样式 */
.app-main {
  display: flex; /* 使用Flex布局 */
  align-items: center; /* 垂直居中对齐 */
  justify-content: center; /* 水平居中对齐 */
  position: relative; /* 相对定位 */
}

/* 定义应用主体的伪元素:before */
.app-main::before {
  content: ''; /* 伪元素内容为空 */
  position: absolute; /* 绝对定位 */
  top: 0; /* 顶部紧贴父元素 */
  left: 0; /* 左侧紧贴父元素 */
  right: 0; /* 右侧紧贴父元素 */
  bottom: 0; /* 底部紧贴父元素 */
  background-image: inherit; /* 继承背景图片 */
  filter: blur(10px); /* 模糊效果 */
  z-index: -1; /* 低于父元素的内容 */
}

/* 定义应用主体中的行元素 */
.app-main .el-row {
  position: relative; /* 相对定位 */
  z-index: 1; /* 高于伪元素 */
}

/* 定义按钮的样式 */
.el-button {
  width: 100%; /* 宽度占满父容器 */
  font-size: 18px; /* 字体大小 */
  border-radius: 10px; /* 圆角边框 */
}

/* 按钮悬停时的动态效果 */
.el-button:hover {
  transform: translateY(-2px); /* 竖直方向移动-2px */
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2); /* 阴影效果 */
}
</style>