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
            <!-- 编辑按钮 -->
            <div v-if="isLogin" class="w-5 h-5" @click="deleteMsg(msg.id)">
              <UIcon name="i-mdi-paper-roll-outline" class="text-gray-400" />
            </div>
          </div>

          <div class="border-l-2 border-gray-300 p-6 ml-1">
            <div class="content-container" v-if="msg.image_url || msg.content">
              <!-- 图片内容 -->
              <a v-if="msg.image_url" :href="`${BASE_API}${msg.image_url}`" data-fancybox>
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
      <div v-else-if="message.messages.length > 0" class="text-center text-gray-500 mt-4">
        <UIcon name="i-fluent-emoji-flat-confetti-ball" size="lg" />
        加载完毕~
      </div>
    </div>
        <!-- 来源信息 - 固定在底部 -->
        <div class="text-center text-xs text-gray-400 py-4">
      来自<a href="https://www.noisework.cn" target="_blank" rel="noopener noreferrer" class="text-orange-400 hover:text-orange-500">Noise</a>
      使用<a href="https://github.com/lin-snow/Ech0" target="_blank" rel="noopener noreferrer" class="text-orange-400 hover:text-orange-500">Ech0</a>发布
    </div>
  </div>
</template>

<script setup lang="ts">
import { Fancybox } from "@fancyapps/ui";
import "@fancyapps/ui/dist/fancybox/fancybox.css";
import { useMessageStore } from "~/store/message";
import { useUserStore } from "~/store/user";
import MarkdownRenderer from "~/components/index/MarkdownRenderer.vue";

const BASE_API = useRuntimeConfig().public.baseApi;
const { deleteMessage } = useMessage();
const message = useMessageStore();
const isLogin = ref<boolean>(false);

const deleteMsg = (id: number) => {
  const confirmDelete = confirm("确定要删除这条消息吗？");
  if (confirmDelete) {
    deleteMessage(id);
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

  message.getMessages({
    page: 1,
    pageSize: 10,
  });
  Fancybox.bind("[data-fancybox]", {});
});

onBeforeUnmount(() => {
  Fancybox.destroy();
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
  margin: 6px 0;
  width: 100%;
  box-sizing: border-box;
}

.content-container img {
  width: 100%;
  height: auto;
  min-height: 200px;
  object-fit: cover;
}
</style>