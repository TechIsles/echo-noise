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

    <!-- 显示图片 -->
    <div
      v-if="ImageUrl"
      class="w-full h-auto shadow-md rounded-md mx-auto mt-5 overflow-hidden"
    >
      <a :href="`${BASE_API}${ImageUrl}`" data-fancybox>
        <img :src="`${BASE_API}${ImageUrl}`" alt="" loading="lazy" />
      </a>
    </div>
  </UCard>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import type { MessageToSave } from "~/types/models";
import { UButton } from "#components";
import { useMessage } from "~/composables/useMessage";
import { Fancybox } from "@fancyapps/ui";
import "@fancyapps/ui/dist/fancybox/fancybox.css";
import VditorEditor from './VditorEditor.vue'

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

  if (!file) {
    toast.add({
      title: '错误',
      description: '没有选择文件',
      color: 'red'
    })
    return;
  }

  const imageUrl = await uploadImage(file);

  if (!imageUrl) {
    toast.add({
      title: '错误',
      description: '图片上传失败，请稍后重试',
      color: 'red'
    })
    return;
  }

  ImageUrl.value = imageUrl;
  toast.add({
    title: '成功',
    description: '图片上传成功',
    color: 'green'
  })
};

onMounted(() => {
  Fancybox.bind("[data-fancybox]", {});
});

onBeforeUnmount(() => {
  Fancybox.destroy();
});
</script>