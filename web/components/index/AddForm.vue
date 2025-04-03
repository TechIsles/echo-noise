<template>
  <UCard class="mx-auto sm:max-w-2xl hover:shadow-md backdrop-blur-sm bg-black/40 shadow-lg text-white">
    <div class="flex justify-between mb-3">
      <div class="flex justify-start items-center gap-2">
        <UIcon name="i-fluent-emoji-flat-alien-monster" class="w-6 h-6" />
        <h2 class="text-lg font-bold italic text-white">说说·笔记~</h2>
      </div>
      <div class="flex gap-2">
        <ClientOnly>
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

  <!-- 预览区域 - 仅在有图片时显示 -->
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
</template>
<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import type { MessageToSave } from "~/types/models";
import { UButton } from "#components";
import { useMessage } from "~/composables/useMessage";
import { useUserStore } from "~/store/user";
import { Fancybox } from "@fancyapps/ui";
import "@fancyapps/ui/dist/fancybox/fancybox.css";
import VditorEditor from './VditorEditor.vue'

const toast = useToast()
const BASE_API = useRuntimeConfig().public.baseApi;
const { save, uploadImage } = useMessage();

const Username = ref("");
const MessageContent = ref("");
const Private = ref<boolean>(false);
const ImageUrl = ref("");
const fileInput = ref<HTMLInputElement | null>(null);
const vditorEditor = ref<InstanceType<typeof VditorEditor> | null>(null);

const privateIcon = computed(() => (Private.value ? 'i-mdi-eye-off-outline' : 'i-mdi-eye-outline'));

const clearForm = () => {
  Username.value = "";
  MessageContent.value = "";
  ImageUrl.value = "";
  Private.value = false;
  if (vditorEditor.value) {
    vditorEditor.value.clear();
  }
};

const addMessage = async () => {
  const message: MessageToSave = {
    username: Username.value,
    content: MessageContent.value,
    private: Private.value,
    image_url: ImageUrl.value,
  };

  const response = await save(message);
  if (response) {
    clearForm();
  }
};

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
});

onBeforeUnmount(() => {
  Fancybox.destroy();
});
</script>