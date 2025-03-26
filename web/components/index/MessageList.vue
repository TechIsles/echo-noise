<template>
  <div class="min-h-screen flex flex-col">
    <div class="flex-grow mx-auto w-full sm:max-w-2xl px-2">
      <div class="my-4">
        <div
          v-for="msg in message.messages"
          :key="msg.id"
          class="w-full h-auto overflow-hidden flex flex-col justify-between"
        >
          <div class="flex justify-between items-center">
            <!-- 时间 -->
            <div class="flex justify-start items-center h-auto">
              <div class="w-2 h-2 rounded-full bg-orange-600 mr-2"></div>
              <!-- 保留年月日 -->
              <div class="flex justify-start text-sm text-orange-500">
                {{ formatDate(msg.created_at) }}
              </div>
            </div>
            <!-- 操作按钮 -->
            <div class="flex justify-end items-center space-x-2">
              <!-- 显示是否为私密 -->
              <div v-if="msg.private" class="w-5 h-5">
                <UIcon name="i-mdi-lock-outline" class="text-gray-400" />
              </div>
              <!-- 评论按钮 -->
              <div
                class="w-5 h-5 cursor-pointer"
                @click="toggleComment(msg.id)"
              >
                <UIcon
                  name="i-mdi-comment-outline"
                  class="text-gray-400 hover:text-orange-500"
                />
              </div>
              <!-- 删除按钮 -->
              <div
                v-if="isLogin"
                class="w-5 h-5 cursor-pointer"
                @click="deleteMsg(msg.id)"
              >
                <UIcon
                  name="i-mdi-paper-roll-outline"
                  class="text-gray-400 hover:text-orange-500"
                />
              </div>
            </div>
          </div>

          <div class="border-l-2 border-gray-300 p-6 ml-1">
            <div class="content-container" v-if="msg.image_url || msg.content">
              <!-- 图片内容 -->
              <a
                v-if="msg.image_url"
                :href="`${BASE_API}${msg.image_url}`"
                data-fancybox="uploaded-image"
              >
                <img
                  :src="`${BASE_API}${msg.image_url}`"
                  alt="Image"
                  class="max-w-full object-cover rounded-lg mb-4"
                  loading="lazy"
                />
              </a>
              <!-- 分隔线 -->
              <div
                v-if="msg.image_url && msg.content"
                class="border-t border-gray-600 my-4"
              ></div>
              <!-- 文本内容 -->
              <div class="overflow-y-auto">
                <MarkdownRenderer :content="msg.content" />
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
          class="rounded-full border-gray-200"
          @click="message.getMessages({ page: message.page + 1, pageSize: 10 })"
        >
          加载更多...
        </UButton>
      </div>
      <!-- 加载完毕提示~ -->
      <div
        v-else-if="message.messages.length > 0"
        class="text-center text-gray-500 mt-4"
      >
        <UIcon name="i-fluent-emoji-flat-confetti-ball" size="lg" />
        加载完毕~
      </div>
    </div>
    <!-- 来源信息 - 固定在底部 -->
    <div class="text-center text-xs text-gray-400 py-4">
      来自<a
        href="https://www.noisework.cn"
        target="_blank"
        rel="noopener noreferrer"
        class="text-orange-400 hover:text-orange-500"
        >Noise</a
      >
      使用<a
        href="https://github.com/lin-snow/Ech0"
        target="_blank"
        rel="noopener noreferrer"
        class="text-orange-400 hover:text-orange-500"
        >Ech0</a
      >发布
    </div>
  </div>
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

const deleteMsg = (id: number) => {
  const confirmDelete = confirm("确定要删除这条消息吗？");
  if (confirmDelete) {
    deleteMessage(id);
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
      }
    };

    // 处理 Markdown 中的图片
    const mdImages = document.querySelectorAll(".markdown-preview img");
    mdImages.forEach((img) => {
      // 移除已存在的包装器
      const parent = img.parentElement;
      if (parent && parent.hasAttribute('data-fancybox')) {
        parent.replaceWith(img);
      }
      
      const src = img.getAttribute("src");
      const wrapper = document.createElement("a");
      wrapper.href = src;
      wrapper.setAttribute("data-fancybox", "uploaded-image");
      wrapper.style.display = 'block';
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

onMounted(() => {
  isLogin.value = useUserStore()?.isLogin;

 // 确保 Waline 脚本加载完成
  if (!window.Waline) {
    console.warn("Waline 未加载，尝试加载 Waline 脚本");
    const script = document.createElement("script");
    script.src = "https://unpkg.com/@waline/client@v2/dist/waline.js";
    script.onload = () => {
      console.log("Waline 脚本加载成功");
      // Waline 加载完成后初始化 Fancybox
      initFancybox();
    };
    document.head.appendChild(script);
  }

  // 添加定时器确保图片加载完成后初始化
  const initTimer = setInterval(() => {
    const images = document.querySelectorAll('.markdown-preview img');
    if (images.length > 0) {
      initFancybox();
      clearInterval(initTimer);
    }
  }, 1000);

  if (document.readyState === "complete") {
    message.getMessages({
      page: 1,
      pageSize: 10,
    });
  } else {
    window.addEventListener("load", () => {
      message.getMessages({
        page: 1,
        pageSize: 10,
      });
    });
  }
});

// 监听消息变化，确保新加载的内容也能正确绑定 Fancybox
watch(
  () => message.messages,
  () => {
    nextTick(() => {
      initFancybox();
    });
  },
  { deep: true }
);

onBeforeUnmount(() => {
  if (window.Fancybox) {
    window.Fancybox.destroy();
  }
});
</script>

<style scoped>
.content-container {
  padding: 8px;
  background: rgba(36, 43, 50, 0.95);
  border-radius: 12px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.593);
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid rgba(14, 14, 14, 0.2);
  margin: 4px 0;
  width: 100%;
  box-sizing: border-box;
}

.content-container img {
  width: 100%;
  height: auto;
  min-height: 200px;
  object-fit: cover;
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
</style>
