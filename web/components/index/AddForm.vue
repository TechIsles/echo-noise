<template>
  <UCard class="mx-auto sm:max-w-2xl hover:shadow-md backdrop-blur-sm bg-black/40 shadow-lg text-white">
    <div class="flex justify-between mb-3">
      <div class="flex justify-start items-center gap-2">
        <UIcon 
          name="i-heroicons-pencil-square" 
          class="w-6 h-6 transition-all duration-300 hover:scale-110 hover:text-blue-400 animate-pulse" 
        />
        <h2 class="text-lg font-bold italic text-white">说说·笔记~</h2>
      </div>
      <div class="flex gap-2">
        <ClientOnly>
          <!-- 添加搜索按钮 -->
         <div @click="showSearchModal=true" class="cursor-pointer flex">
         <UIcon name="i-heroicons-magnifying-glass" class="w-5 h-5 text-gray-200" />
        </div>
          <!-- 添加热力图图标 -->
          <button @click="toggleHeatmap">
            <UIcon name="i-mdi-calendar-month" class="w-5 h-5 text-gray-200" />
          </button>
          <a href="/rss" target="_blank">
            <UIcon name="i-mdi-rss" class="w-5 h-5 text-gray-200" />
          </a>
        </ClientOnly>
        <NuxtLink to="/status">
          <UIcon name="i-mdi-server-outline" class="w-5 h-5 text-gray-200" />
        </NuxtLink>
      </div>
    </div>

    <div>
      <VditorEditor ref="vditorEditor" v-model="MessageContent" />
      <div class="flex justify-between items-center">
        <div class="flex items-center justify-start gap-2">
          <!-- <UInput size="sm" color="gray" :trailing="true" placeholder="你の名字" v-model="Username" class="w-24" /> -->
          <input
            id="file-input"
            ref="fileInput"
            type="file"
            accept="image/*"
            @change="addImage"
            class="hidden"
            placeholder="选择图片"
          />
          <!-- 使用 NuxtUI 的 UButton 触发文件选择 -->
          <UButton
            color="gray"
            variant="solid"
            class="cursor-pointer"
            size="sm"
            icon="i-fluent-image-20-regular"
            @click="triggerFileInput"
          />
                  <!-- 是否设为私密？ -->
                  <UButton
             color="gray"
             variant="solid"
             size="sm"
             @click="Private = !Private"
             :icon="privateIcon"
             :title="Private ? '设为公开' : '设为私密'"
             :ui="{ tooltip: { text: Private ? '设为公开' : '设为私密' } }"
            />
       
<!-- 推送开关 -->
<UButton
  color="gray"
  variant="solid"
  size="sm"
  @click="toggleNotify"
  :icon="enableNotify ? 'i-mdi-bell' : 'i-mdi-bell-off'"
  :title="enableNotify ? '关闭推送' : '开启推送'"
  :ui="{ tooltip: { text: enableNotify ? '关闭推送' : '开启推送' } }"
  class="text-gray-600" 
/>
</div>
        <div class="flex gap-2">
          <!-- 清空表单 -->
          <UButton
            icon="i-fluent-broom-16-regular"
            variant="solid"
            color="gray"
            size="sm"
            @click="clearForm"
          />

          <!-- 添加笔记 -->
          <UButton
            icon="i-fluent-add-12-filled"
            variant="solid"
            color="gray"
            size="sm"
            @click="addMessage"
          />
        </div>
      </div>
    </div>
  </UCard>

  <!-- 图片预览区域 - 仅在有图片时显示 -->
  <div
    v-if="ImageUrl"
    class="mx-auto sm:max-w-2xl mt-5 backdrop-blur-sm bg-black/40 p-4 rounded-md"
  >
    <!-- 图片预览 -->
    <div class="shadow-md rounded-md mx-auto overflow-hidden">
      <a :href="`${BASE_API}${ImageUrl}`" data-fancybox>
        <img :src="`${BASE_API}${ImageUrl}`" alt="" loading="lazy" class="w-full h-auto" />
      </a>
    </div>

    <!-- 分割线 -->
    <hr class="border-gray-600 my-4">

    <!-- 文本预览 -->
    <div v-if="MessageContent" class="prose prose-invert max-w-none">
      <div v-html="MessageContent"></div>
    </div>
  </div>
  <!-- 添加搜索模态框组件 -->
  <SearchMode 
    v-model="showSearchModal" 
    @search-result="handleSearchResult" 
  />
</template>
<script setup lang="ts">
import { ref, computed, inject, onMounted, onUnmounted, watch } from 'vue'
import type { MessageToSave } from "~/types/models";
import { UButton } from "#components";
import { useMessage } from "~/composables/useMessage";
import { useUserStore } from '~/store/user'
import { Fancybox } from '@fancyapps/ui'
import '@fancyapps/ui/dist/fancybox/fancybox.css'
import VditorEditor from './VditorEditor.vue'
// 导入搜索模式组件
import SearchMode from './Searchmode.vue'
import { useMessageStore } from '~/store/message'
import { useNotifyStore } from '~/store/notify'

// 添加搜索相关变量和函数
const showSearchModal = ref(false);
const emit = defineEmits(['search-result']);
// 处理搜索结果
const handleSearchResult = (result: any) => {
  // 将搜索结果传递给父组件
  emit('search-result', result);
};

const toast = useToast()
const BASE_API = useRuntimeConfig().public.baseApi;
const { save, uploadImage } = useMessage();

// 添加热力图控制变量
const showHeatmap = inject('showHeatmap') as Ref<boolean>;

// 提供给父组件的方法，用于控制热力图显示
provide('showHeatmap', showHeatmap);

// 切换热力图显示状态
const toggleHeatmap = () => {
  console.log('热力图按钮被点击，当前状态:', showHeatmap.value);
  showHeatmap.value = !showHeatmap.value;
  console.log('切换后状态:', showHeatmap.value);
};

const Username = ref("");
const MessageContent = ref("");
const Private = ref<boolean>(false);
const ImageUrl = ref("");
const fileInput = ref<HTMLInputElement | null>(null);
const vditorEditor = ref<InstanceType<typeof VditorEditor> | null>(null);

const privateIcon = computed(() => (Private.value ? 'i-mdi-eye-off-outline' : 'i-mdi-eye-outline'));

// 在顶部 script 区域添加 enableNotify 的定义
const notifyStore = useNotifyStore()
const enableNotify = ref(localStorage.getItem('enableNotify') === 'true')

// 保留这一个 clearForm 函数，删除另一个重复的定义
const clearForm = () => {
  Username.value = "";
  MessageContent.value = "";
  ImageUrl.value = "";
  Private.value = false;
  // 不重置推送状态，保持用户的选择
  if (vditorEditor.value) {
    vditorEditor.value.clear();
  }
};

const addMessage = async () => {
  // 添加基本验证
  if (!MessageContent.value.trim() && !ImageUrl.value) {
    toast.add({
      title: '错误',
      description: '请输入内容或上传图片',
      color: 'red',
      timeout: 2000
    });
    return;
  }

  const message: MessageToSave = {
    username: Username.value,
    content: MessageContent.value,
    private: Private.value,
    image_url: ImageUrl.value,
    // 修改推送逻辑：私密内容不进行推送
    notify: Private.value ? false : enableNotify.value,
  };

  try {
    const response = await save(message);
    if (response) {
      toast.add({
        title: '成功',
        description: '发布成功',
        color: 'green',
        timeout: 2000
      });
      clearForm();
    }
  } catch (error: any) {
    console.error('发布错误:', error);
    toast.add({
      title: '错误',
      description: error.message || '发布失败',
      color: 'red',
      timeout: 2000
    });
  }
};

// 删除这里重复的 clearForm 函数定义

const triggerFileInput = () => {
  const input = document.getElementById("file-input");
  if (input) {
    input.click();
  }
};

const addImage = async (event: Event) => {
  const input = event.target as HTMLInputElement;
  const file = input.files ? input.files[0] : null;

  const userStore = useUserStore();
  if (!userStore.isLogin) {
    toast.add({
      title: '错误',
      description: '请先登录',
      color: 'red',
      timeout: 2000
    });
    return;
  }

  if (!file) {
    toast.add({
      title: '错误',
      description: '没有选择文件',
      color: 'red',
      timeout: 2000
    });
    return;
  }

  // 检查文件类型和扩展名
  const allowedTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'];
  const allowedExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.webp'];
  
  const fileExtension = file.name.toLowerCase().substring(file.name.lastIndexOf('.'));
  
  if (!allowedTypes.includes(file.type) || !allowedExtensions.includes(fileExtension)) {
    toast.add({
      title: '错误',
      description: '仅支持 JPG、PNG、GIF、WEBP 格式的图片',
      color: 'red',
      timeout: 2000
    });
    return;
  }

  // 检查文件大小（限制为 5MB）
  const maxSize = 5 * 1024 * 1024; // 5MB
  if (file.size > maxSize) {
    toast.add({
      title: '错误',
      description: '图片大小不能超过 5MB',
      color: 'red',
      timeout: 2000
    });
    return;
  }

  try {
    // 创建 FormData 对象
    const formData = new FormData();
    formData.append('image', file);

    // 发送请求
    const response = await fetch(`${BASE_API}/images/upload`, {
      method: 'POST',
      body: formData,
      credentials: 'include'
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.msg || '图片上传失败');
    }

    const data = await response.json();
    if (data.code === 1 && data.data) {
      ImageUrl.value = data.data;
      // 直接在编辑器中插入图片
      if (vditorEditor.value?.vditor) {
        const imageMarkdown = `\n![](${BASE_API}${data.data})\n`;
        vditorEditor.value.vditor.insertValue(imageMarkdown);
        // 确保编辑器内容更新
        MessageContent.value = vditorEditor.value.vditor.getValue();
      }
      toast.add({
        title: '成功',
        description: '图片上传成功',
        color: 'green',
        timeout: 2000
      });
    } else {
      throw new Error(data.msg || '图片上传失败');
    }
  } catch (error: any) {
    console.error('上传错误:', error);
    toast.add({
      title: '错误',
      description: error.message || '图片上传失败',
      color: 'red',
      timeout: 2000
    });
  } finally {
    if (fileInput.value) {
      fileInput.value.value = '';
    }
  }
};

onMounted(() => {
  Fancybox.bind("[data-fancybox]", {});
  console.log('组件加载时推送开关状态:', enableNotify.value);
});

onBeforeUnmount(() => {
  Fancybox.destroy();
});
// 添加切换推送状态的方法
const toggleNotify = () => {
  enableNotify.value = !enableNotify.value;
  localStorage.setItem('enableNotify', enableNotify.value.toString());
  console.log('推送开关状态:', enableNotify.value);
};
</script>