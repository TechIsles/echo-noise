<template>
   <div class="min-h-screen flex flex-col">
    <div class="flex-grow mx-auto w-full sm:max-w-2xl px-2">
      <div class="my-4">
        <div v-for="msg in message.messages" :key="msg.id"
          class="w-full h-auto overflow-hidden flex flex-col justify-between">
           <!-- 修改头部布局 -->
           <div class="flex justify-between items-center">
            <!-- 时间部分保持不变 -->
            <div class="flex justify-start items-center h-auto overflow-x-auto whitespace-nowrap hide-scrollbar">
  <div class="w-2 h-2 rounded-full bg-orange-600 mr-2 flex-shrink-0"></div>
  <div class="flex justify-start text-sm">
    <span class="text-orange-500">{{ formatDate(msg.created_at) }}</span>
    <span class="gradient-dot mx-2 flex-shrink-0">@</span>
    <span class="text-orange-500">{{ msg.username || '匿名用户' }}</span>
  </div>
</div>
            <!-- 优化操作按钮组样式 -->
          <div class="message-actions flex justify-end items-center space-x-2 flex-shrink-0 px-3 py-1.5 mr-[9px] -mb-[1px]">
            <!-- ... 按钮内容 ... -->
              <div v-if="msg.private" class="w-5 h-5 flex-shrink-0 transition-transform duration-200 hover:scale-110">
                <UIcon name="i-mdi-lock-outline" class="text-gray-400" />
              </div>
              <div v-if="isLogin" class="w-5 h-5 cursor-pointer flex-shrink-0 transition-all duration-200 hover:scale-110" @click="editMessage(msg)" :title="'编辑内容'">
                <UIcon name="i-mdi-pencil-outline" class="text-gray-400 hover:text-orange-500" />
              </div>
              <div class="w-5 h-5 cursor-pointer flex-shrink-0 transition-all duration-200 hover:scale-110" @click="copyContent(msg.content)" :title="'复制内容'">
                <UIcon name="i-mdi-content-copy" class="text-gray-400 hover:text-orange-500" />
              </div>
              <div class="w-5 h-5 cursor-pointer flex-shrink-0 transition-all duration-200 hover:scale-110" @click="downloadAsImage(msg.id)" :title="'下载为图片'">
                <UIcon name="i-mdi-image-outline" class="text-gray-400 hover:text-orange-500" />
              </div>
              <div class="w-5 h-5 cursor-pointer flex-shrink-0 transition-all duration-200 hover:scale-110" @click="toggleComment(msg.id)" :title="'评论'">
                <UIcon name="i-mdi-comment-outline" class="text-gray-400 hover:text-orange-500" />
              </div>
              <div v-if="isLogin" class="w-5 h-5 cursor-pointer flex-shrink-0 transition-all duration-200 hover:scale-110" @click="deleteMsg(msg.id)" :title="'删除'">
                <UIcon name="i-mdi-paper-roll-outline" class="text-gray-400 hover:text-orange-500" />
              </div>
            </div>
          </div>

          <div class="border-l-2 border-gray-300 p-6 ml-1">
            <div class="content-container" v-if="msg.image_url || msg.content" :data-msg-id="msg.id">
              <!-- 图片内容 -->
              <a v-if="msg.image_url" :href="`${BASE_API}${msg.image_url}`" data-fancybox="uploaded-image">
                <img :src="`${BASE_API}${msg.image_url}`" alt="Image" class="max-w-full object-cover rounded-lg mb-4"
                  loading="lazy" />
              </a>
              <!-- 分隔线 -->
              <div v-if="msg.image_url && msg.content" class="border-t border-gray-600 my-4"></div>
              <!-- 文本内容区域 -->
              <div class="overflow-y-hidden relative" :class="{ 'max-h-[700px]': !isExpanded[msg.id] }">
                <MarkdownRenderer :content="msg.content" />
                <div v-if="shouldShowExpandButton[msg.id] && !isExpanded[msg.id]"
                  class="absolute bottom-0 left-0 right-0 h-32 bg-gradient-to-t from-[rgba(36,43,50,1)] via-[rgba(36,43,50,0.8)] to-transparent">
                </div>
              </div>
              <!-- 展开/折叠按钮 -->
              <div v-if="shouldShowExpandButton[msg.id]" class="text-center mt-2 relative" style="z-index: 9999;">
                <button @click="toggleExpand(msg.id)"
                  class="flex items-center justify-center space-x-1 mx-auto px-4 py-2 text-orange-500 hover:text-orange-600 focus:outline-none transition-colors duration-200">
                  <span>{{ isExpanded[msg.id] ? "收起内容" : "展开全文" }}</span>
                  <UIcon :name="isExpanded[msg.id] ? 'i-mdi-chevron-up' : 'i-mdi-chevron-down'" class="w-5 h-5" />
                </button>
              </div>
            </div>
            <!-- 评论区域 -->
            <div v-show="activeCommentId === msg.id" class="mt-4">
              <div :id="`waline-${msg.id}`"></div>
            </div>
          </div>
        </div>
      </div>
        <!-- 加载更多 -->
        <div v-if="message.hasMore" class="flex justify-center w-full my-4">
        <UButton 
          color="gray" 
          variant="outline" 
          size="sm" 
          class="rounded-full px-6 py-2 bg-[rgba(36,43,50,0.95)] text-white hover:text-white border-none shadow-lg hover:shadow-xl transition-all duration-300 backdrop-blur-sm"
          @click="message.getMessages({ page: message.page + 1, pageSize: 10 })"
        >
          加载更多...
        </UButton>
      </div>
      <!-- 加载完毕提示~ -->
      <div v-else-if="message.messages.length > 0" class="text-center text-gray-500 mt-4">
        <UIcon name="i-fluent-emoji-flat-confetti-ball" size="lg" />
        加载完毕~
      </div>
    </div>
    <!-- 来源信息 - 固定在底部 -->
    <div class="text-center text-xs text-gray-400 py-4">
      来自<a href="https://www.noisework.cn" target="_blank" rel="noopener noreferrer"
        class="text-orange-400 hover:text-orange-500">Noise</a>
      使用<a href="https://github.com/lin-snow/Ech0" target="_blank" rel="noopener noreferrer"
        class="text-orange-400 hover:text-orange-500">Ech0</a>发布
    </div>
  </div>
  <!-- 编辑对话框 -->
  <UModal v-model="showEditModal" :ui="{ width: 'sm:max-w-2xl' }">
    <UCard>
      <template #header>
        <div class="flex justify-between items-center">
          <h3 class="text-lg font-medium">编辑内容</h3>
          <UButton color="gray" variant="ghost" icon="i-mdi-close" class="-my-1" @click="showEditModal = false" />
        </div>
      </template>
      <div class="flex flex-col space-y-4">
        <UTextarea
          v-model="editingContent"
          placeholder="编辑内容..."
          rows="10"
          class="font-mono text-sm"
        />
        <div class="border-t border-gray-200 my-2 pt-2">
          <div class="text-sm text-gray-500 mb-2">预览：</div>
          <!-- 修改预览区域样式 -->
          <div class="p-4 rounded-lg overflow-auto max-h-[300px] bg-[rgba(36,43,50,0.95)]">
            <div class="text-white">
              <MarkdownRenderer :content="editingContent" />
            </div>
          </div>
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end space-x-2">
          <UButton color="gray" variant="outline" @click="showEditModal = false">
            取消
          </UButton>
          <UButton color="orange" @click="saveEditedMessage" :loading="isSaving">
            保存
          </UButton>
        </div>
      </template>
    </UCard>
  </UModal>
</template>

<script setup lang="ts">
import { useMessageStore } from "~/store/message";
import { useUserStore } from "~/store/user";
import MarkdownRenderer from "~/components/index/MarkdownRenderer.vue";

const BASE_API = useRuntimeConfig().public.baseApi;
const { deleteMessage } = useMessage();
const message = useMessageStore();
const isLogin = ref<boolean>(false);
const activeCommentId = ref<number | null>(null);

const deleteMsg = async (id: number) => {
  const confirmDelete = confirm("确定要删除这条消息吗？");
  if (confirmDelete) {
    try {
      await deleteMessage(id);
      // 只在本地移除消息，不重新加载
      message.messages = message.messages.filter(msg => msg.id !== id);
      useToast().add({
        title: '删除成功',
        color: 'green',
        timeout: 2000
      });
    } catch (error) {
      console.error('删除失败:', error);
      useToast().add({
        title: '删除失败',
        color: 'red',
        timeout: 2000
      });
    }
  }
};

const initFancybox = () => {
  if (window.Fancybox) {
    // 先解绑之前的事件，避免冲突
    window.Fancybox.destroy();
    // 统一配置
    const fancyboxOptions = {
      Carousel: {
        infinite: false,
      },
      Toolbar: {
        display: [
          { id: "prev", position: "center" },
          { id: "counter", position: "center" },
          { id: "next", position: "center" },
          "zoom",
          "slideshow",
          "fullscreen",
          "close",
        ],
      },
      Image: {
        zoom: true,
        click: true,
        wheel: "slide",
      },
    };

    // 处理 Markdown 中的图片
    const mdImages = document.querySelectorAll(".markdown-preview img");
    mdImages.forEach((img) => {
      // 移除已存在的包装器
      const parent = img.parentElement;
      if (parent && parent.hasAttribute("data-fancybox")) {
        parent.replaceWith(img);
      }

      const src = img.getAttribute("src");
      const wrapper = document.createElement("a");
      wrapper.href = src;
      wrapper.setAttribute("data-fancybox", "uploaded-image");
      wrapper.style.display = "block";
      img.parentNode.insertBefore(wrapper, img);
      wrapper.appendChild(img);
    });

    // 最后统一绑定事件
    window.Fancybox.bind("[data-fancybox]", fancyboxOptions);
  }
};

const toggleComment = async (msgId: number) => {
  if (activeCommentId.value === msgId) {
    activeCommentId.value = null;
  } else {
    activeCommentId.value = msgId;
    await nextTick();
    // 在评论加载后重新初始化 Fancybox
    initFancybox();
    // 确保 Waline 初始化
    if (window.Waline) {
      // 检查元素是否存在
      const el = document.querySelector(`#waline-${msgId}`);
      if (el) {
        window.Waline.init({
          el: `#waline-${msgId}`,
          serverURL: "https://app-production-80c1.up.railway.app",
          path: `/message/${msgId}`,
          reaction: false,
          meta: ["nick", "mail", "link"],
          requiredMeta: ["mail", "nick"],
          pageview: true,
          search: false,
          wordLimit: 200,
          pageSize: 5,
          avatar: "monsterid",
          emoji: ["https://unpkg.com/@waline/emojis@1.2.0/tieba"],
          imageUploader: false,
          copyright: false,
          dark: 'html[class="dark"]',
        });
      } else {
        console.error(`评论容器 #waline-${msgId} 未找到`);
      }
    } else {
      console.error("Waline 未加载");
    }
  }
};

const formatDate = (dateString: string) => {
  const date = new Date(dateString);
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  const diffInDays = Math.floor(diff / (1000 * 60 * 60 * 24));
  const diffInHours = Math.floor(diff / (1000 * 60 * 60));
  const diffInMinutes = Math.floor(diff / (1000 * 60));

  const diffInSeconds = Math.floor(diff / 1000);
  if (diffInSeconds < 60) {
    return "刚刚";
  } else if (diffInMinutes < 60) {
    return `${diffInMinutes}分钟前`;
  } else if (diffInHours < 24) {
    return `${diffInHours}小时前`;
  } else if (diffInDays < 3) {
    return `${diffInDays}天前`;
  } else {
    return date.toLocaleString();
  }
};
// 添加展开状态管理
const isExpanded = ref<{ [key: number]: boolean }>({});
const shouldShowExpandButton = ref<{ [key: number]: boolean }>({});

// 添加展开/折叠切换函数
const toggleExpand = (msgId: number) => {
  isExpanded.value[msgId] = !isExpanded.value[msgId];
};

// 修改检查内容高度的函数
const checkContentHeight = () => {
  nextTick(() => {
    message.messages.forEach((msg) => {
      const contentEl = document.querySelector(
        `.content-container[data-msg-id="${msg.id}"] .overflow-y-hidden`
      );
      if (contentEl && contentEl.scrollHeight > 700) {
        shouldShowExpandButton.value[msg.id] = true;
        if (isExpanded.value[msg.id] === undefined) {
          isExpanded.value[msg.id] = false;
        }
      }
    });
  });
};

// 确保在内容变化时重新检查高度
watch(
  () => message.messages,
  () => {
    nextTick(async () => {
      await new Promise(resolve => setTimeout(resolve, 100)); // 添加短暂延迟
      checkContentHeight();
      initFancybox();
    });
  },
  { deep: true, immediate: true }
);

onMounted(async () => {
  isLogin.value = useUserStore()?.isLogin;
  
  // 获取消息列表
  await message.getMessages({
    page: 1,
    pageSize: 10,
  });

  // 确保内容加载后检查高度
  await nextTick();
  await new Promise(resolve => setTimeout(resolve, 100));
  checkContentHeight();

  // Waline 加载逻辑
  if (!window.Waline) {
    const script = document.createElement("script");
    script.src = "https://unpkg.com/@waline/client@v2/dist/waline.js";
    script.onload = initFancybox;
    document.head.appendChild(script);
  } else {
    initFancybox();
  }
});

// 简化监听器 - 移除重复的 messages 监听
watch(
  () => message.messages,
  () => {
    nextTick(() => {
      checkContentHeight();
      initFancybox(); // 合并为一个调用
    });
  },
  { deep: true }
);

// 移除组件卸载时的状态重置
onBeforeUnmount(() => {
  if (window.Fancybox) {
    window.Fancybox.destroy();
  }
});
// 添加复制功能
const copyContent = async (content: string) => {
  try {
    await navigator.clipboard.writeText(content);
    // 可以使用 Nuxt 的 toast 提示复制成功
    useToast().add({
      title: '复制成功',
      color: 'green',
      timeout: 2000
    });
  } catch (err) {
    console.error('复制失败:', err);
    useToast().add({
      title: '复制失败',
      color: 'red',
      timeout: 2000
    });
  }
};
// 添加编辑功能
const showEditModal = ref(false);
const editingContent = ref('');
const editingMessageId = ref<number | null>(null);
const isSaving = ref(false);

const editMessage = (msg: any) => {
  editingMessageId.value = msg.id;
  editingContent.value = msg.content;
  showEditModal.value = true;
};

const saveEditedMessage = async () => {
  if (!editingMessageId.value) return;
  
  isSaving.value = true;
  try {
    await deleteMessage(editingMessageId.value);
    
    const oldMsg = message.messages.find(msg => msg.id === editingMessageId.value);
    const newMessage: MessageToSave = {
      content: editingContent.value,
      private: oldMsg?.private || false,
      username: oldMsg?.username || '',
      image_url: oldMsg?.image_url || ''
    };

    const { save } = useMessage();
    const savedMessage = await save(newMessage);
    
    // 在本地更新消息列表
    const index = message.messages.findIndex(msg => msg.id === editingMessageId.value);
    if (index !== -1) {
      message.messages[index] = savedMessage;
    }
    
    showEditModal.value = false;
  } catch (error) {
    console.error('更新消息失败:', error);
  } finally {
    isSaving.value = false;
  }
};
const downloadAsImage = async (msgId: number) => {
  try {
    const element = document.querySelector(`.content-container[data-msg-id="${msgId}"]`);
    if (!element) return;

    // 检查内容类型
    const hasText = element.querySelector('.markdown-preview')?.textContent?.trim();
    const hasImage = element.querySelector('img');
    const hasVideo = element.querySelector('video');
    const hasAudio = element.querySelector('audio');

    // 纯视频或纯音频内容不生成卡片
    if ((!hasText && !hasImage && hasVideo) || (!hasText && !hasImage && hasAudio)) {
      useToast().add({
        title: '此内容不可生成卡片',
        color: 'orange',
        timeout: 2000
      });
      return;
    }

    // 设置超时检测
    const timeout = setTimeout(() => {
      useToast().add({
        title: '生成超时',
        description: '卡片生成时间过长，请稍后重试',
        color: 'red',
        timeout: 3000
      });
    }, 10000);

    // 1. 临时展开内容
    const originalExpanded = isExpanded.value[msgId];
    isExpanded.value[msgId] = true;
    await nextTick();

    // 2. 创建临时容器
    const tempContainer = document.createElement('div');
   tempContainer.style.cssText = `
  padding: 16px;
  background: transparent;
  border-radius: 12px;
  width: ${hasImage ? '640px' : '480px'};
  position: absolute;
  left: -9999px;
  top: 0;
  z-index: -1;
  overflow: visible;
  min-height: fit-content;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  border: none;
`;
    document.body.appendChild(tempContainer);
    
    // 3. 复制并处理内容
    const contentClone = element.cloneNode(true) as HTMLElement;
    
    // 移除所有控制元素和限制
    contentClone.querySelectorAll('.text-center.mt-2, .bg-gradient-to-t').forEach(el => el.remove());
    contentClone.style.cssText = `
      max-height: none;
      overflow: visible;
      padding: 0;
      margin: 0;
      background: transparent;
    `;
    
    // 处理内容区域
    const contentArea = contentClone.querySelector('.overflow-y-hidden');
if (contentArea) {
  contentArea.className = '';
  contentArea.style.cssText = `
    overflow: visible;
    max-height: none !important;
    height: auto !important;
    padding: 12px;
    line-height: 1.6;
    margin-bottom: 0;
    white-space: pre-wrap;
    background: transparent;
    border-radius: 12px;
    font-size: 14px;
    color: #333;
  `;
}

    // 处理媒体元素
    const mediaElements = contentClone.querySelectorAll('video, audio, iframe');
    mediaElements.forEach(media => {
      const placeholder = document.createElement('div');
      placeholder.style.cssText = `
        padding: 15px;
        background: rgba(251, 146, 60, 0.1);
        border: 1px solid rgba(251, 146, 60, 0.3);
        border-radius: 8px;
        color: #fb923c;
        margin: 15px 0;
        word-break: break-all;
      `;
      
      if (media instanceof HTMLVideoElement) {
        placeholder.innerHTML = `🎬 视频链接：${media.src || '未知链接'}`;
      } else if (media instanceof HTMLAudioElement) {
        placeholder.innerHTML = `🎵 音频链接：${media.src || '未知链接'}`;
      } else if (media instanceof HTMLIFrameElement) {
        placeholder.innerHTML = `🔗 嵌入内容链接：${media.src || '未知链接'}`;
      }
      
      media.parentNode?.replaceChild(placeholder, media);
    });

   // 处理图片
   const images = contentClone.querySelectorAll('img');
const processImages = async () => {
  await Promise.all(Array.from(images).map(async (img) => {
    return new Promise<void>((resolve) => {
      const originalSrc = img.src;
      img.crossOrigin = 'anonymous';
      
      // 处理图片路径
      if (originalSrc.startsWith('/')) {
        img.src = `${BASE_API}${originalSrc}`;
      }
      
      if (img.complete) {
        resolve();
      } else {
        img.onload = () => resolve();
        img.onerror = () => {
          console.error('图片加载失败:', originalSrc);
          img.parentElement?.removeChild(img);
          resolve();
        };
      }
    });
  }));
};

await processImages();

    tempContainer.appendChild(contentClone);

    // 添加 footer
    const footer = document.createElement('div');
footer.style.cssText = `
  margin-top: 12px;
  padding-top: 12px;
  text-align: center;
  font-family: -apple-system, BlinkMacSystemFont, sans-serif;
  background: transparent;
`;
footer.innerHTML = `
  <div style="color: #fb923c; font-size: 13px; margin-bottom: 4px; font-weight: 500;">
    Noise·说说·笔记~
  </div>
  <div style="color: rgba(255,255,255,0.5); font-size: 11px;">
    note.noisework.cn
  </div>
`;
    tempContainer.appendChild(footer);

    // 生成图片
    await nextTick();
    const canvas = await html2canvas(tempContainer, {
      backgroundColor: null,
      scale: 2,
      useCORS: true,
      allowTaint: true,
      logging: false,
      width: tempContainer.offsetWidth,
      height: tempContainer.scrollHeight,
      borderRadius: '16px',
      onclone: (clonedDoc) => {
        const clonedElement = clonedDoc.querySelector('.content-container');
        if (clonedElement) {
          clonedElement.style.cssText = `
            overflow: visible !important;
            max-height: none !important;
            height: auto !important;
            padding: 0;
            min-height: ${contentArea?.scrollHeight || 0}px;
            background: transparent;
          `;
        }
      }
    });

    // 清除超时检测
    clearTimeout(timeout);
    // 下载图片
    const link = document.createElement('a');
    link.download = `message-${msgId}.png`;
    link.href = canvas.toDataURL('image/png');
    link.click();

    // 清理临时元素
    document.body.removeChild(tempContainer);
    
    // 恢复原始展开状态
    isExpanded.value[msgId] = originalExpanded;

    useToast().add({
      title: '下载成功',
      color: 'green',
      timeout: 2000
    });
  } catch (error) {
    console.error('下载失败:', error);
    useToast().add({
      title: '下载失败',
      color: 'red',
      timeout: 2000
    });
  }
};
</script>

<style scoped>
/* 修改卡片生成样式 */
.content-container {
  padding: 16px;
  background: transparent;
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  transition: all 0.3s ease;
  border: 1px solid rgba(5, 5, 5, 0.2);
  margin: 8px 0;
  width: 100%;
  box-sizing: border-box;
  position: relative;
  overflow: hidden;
}
.content-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(36, 43, 50, 0.95);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  z-index: -1;
  border-radius: inherit;
}

/* 添加展开/折叠按钮样式 */
button {
  background: rgba(36, 43, 50, 0.95);
  border: 1px solid rgba(251, 146, 60, 0.3);
  border-radius: 20px;
  position: relative;
  z-index: 9999;
}

button:hover {
  background: rgba(46, 53, 60, 0.95);
  border-color: rgba(251, 146, 60, 0.5);
  cursor: pointer;
}

/* 确保内容区域的层级正确 */
.overflow-y-hidden {
  transition: max-height 0.3s ease-in-out;
  position: relative;
  z-index: 1;
}
/* 添加内容过渡动画 */
.overflow-y-hidden {
  transition: max-height 0.3s ease-in-out;
}

/* 修正展开状态下的最大高度限制 */
.content-container .overflow-y-hidden:not(.max-h-\[700px\]) {
  max-height: none;
}
/* 优化图片显示 */
.content-container img {
  width: 100%;
  height: auto;
  min-height: 150px;
  object-fit: cover;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}
/* 修改评论区样式 */
:deep(.wl-comment) {
  background: rgba(36, 43, 50, 0.95) !important;
  border-radius: 8px;
  padding: 12px !important;
  margin-bottom: 12px !important;
}

:deep(.wl-comment .wl-content) {
  color: #fff !important;
  background: transparent !important;
}

:deep(.wl-comment .wl-meta),
:deep(.wl-comment .wl-meta > span),
:deep(.wl-comment .wl-meta > a),
:deep(.wl-comment .wl-meta .wl-time),
:deep(.wl-comment .wl-meta .wl-nick) {
  color: #e5e5e5 !important;
}

:deep(.wl-comment .wl-meta .wl-like),
:deep(.wl-comment .wl-meta .wl-reply) {
  color: #999 !important;
}

:deep(.wl-comment .wl-meta .wl-like:hover),
:deep(.wl-comment .wl-meta .wl-reply:hover) {
  color: #fff !important;
}

/* 确保所有评论相关文本为白色 */
:deep(.wl-comment *) {
  color: #fff !important;
}

:deep(.wl-input) {
  color: #fff !important;
}

:deep(.wl-action) {
  color: #fff !important;
}
:deep(.wl-panel) {
  background: rgba(36, 43, 50, 0.95) !important;
  border: 1px solid rgba(14, 14, 14, 0.2) !important;
}

:deep(.wl-editor) {
  background: rgba(36, 43, 50, 0.95) !important;
}

:deep(.wl-header) {
  border-bottom: 1px solid rgba(14, 14, 14, 0.2) !important;
}

:deep(.wl-card) {
  background: rgba(36, 43, 50, 0.95) !important;
  border: 1px solid rgba(14, 14, 14, 0.2) !important;
}
/* 添加评论框样式 */
:deep(.wl-panel),
:deep(.wl-card) {
  position: relative;
  z-index: 100;
}

/* 确保评论区域不会被遮挡 */
.content-container {
  position: relative;
  z-index: 1;
}
/* 添加评论内容文本颜色 */
:deep(.wl-comment .wl-content) {
  color: #fff !important;
}

:deep(.wl-comment .wl-meta) {
  color: #fff !important;
}

:deep(.wl-comment .wl-meta > span),
:deep(.wl-comment .wl-meta > a) {
  color: #fff !important;
}
/* 移除 markdown 图片的 hover 效果 */
:deep(.markdown-preview img) {
  cursor: pointer;
  transform: none !important; /* 移除 hover 时的缩放效果 */
  transition: none !important; /* 移除过渡效果 */
}

:deep(.markdown-preview img:hover) {
  transform: none !important;
}

/* 确保灯箱层级最高 */
:deep(.fancybox__container) {
  --fancybox-bg: rgba(0, 0, 0, 0.9);
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 9999 !important;
}

:deep(.fancybox__backdrop) {
  z-index: 9998 !important;
}
/* 按钮组样式 */
.message-actions {
  position: relative;
  z-index: 1;
}

/* 按钮悬停效果 */
.message-actions > div {
  position: relative;
  transition: all 0.3s ease;
}

.message-actions > div:hover {
  transform: translateY(-2px);
}

.message-actions > div:hover .text-gray-400 {
  color: #fb923c;
  filter: drop-shadow(0 0 2px rgba(251, 146, 60, 0.3));
}
.gradient-dot {
  /* 添加明亮色彩的动态渐变动画 */
  background: linear-gradient(
    45deg,
    #ff6b6b,
    #ffd93d,
    #ff9a9e,
    #cd4e67,
    #ffb347,
    #ff7eb3,
    #ffa07a
  );
  background-size: 400% 400%;
  animation: rainbow 10s ease infinite;
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  font-weight: bold;
}

@keyframes rainbow {
  0% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0% 50%;
  }
}

/* 隐藏滚动条但保持功能 */
.hide-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
.hide-scrollbar::-webkit-scrollbar {
  display: none;
}
</style>
