<template>
  <div 
    class="fixed z-50 bg-black/80 backdrop-blur-sm rounded-lg shadow-lg p-4 text-white"
    :style="{ top: `${position.y}px`, left: `${position.x}px` }"
  >
    <div class="flex justify-between items-center mb-3">
      <h3 class="text-lg font-medium">图床上传</h3>
      <UButton 
        icon="i-heroicons-x-mark" 
        color="gray" 
        variant="ghost" 
        size="xs" 
        @click="$emit('close')" 
      />
    </div>
    
    <div class="mb-4">
      <div class="flex gap-2 mb-3">
        <UButton 
          :color="selectedHost === 'smms' ? 'blue' : 'gray'" 
          variant="solid" 
          size="sm" 
          @click="selectedHost = 'smms'"
        >
          SM.MS
        </UButton>
        <UButton 
          :color="selectedHost === 'github' ? 'blue' : 'gray'" 
          variant="solid" 
          size="sm" 
          @click="selectedHost = 'github'"
        >
          GitHub
        </UButton>
      </div>
      
      <!-- GitHub 配置 -->
      <div v-if="selectedHost === 'github'" class="space-y-2 mb-3">
        <UInput 
          v-model="githubToken" 
          placeholder="GitHub Token" 
          size="sm"
          type="password"
        />
        <UInput 
          v-model="githubRepo" 
          placeholder="仓库名 (用户名/仓库)" 
          size="sm"
        />
        <UInput 
          v-model="githubBranch" 
          placeholder="分支名 (默认: main)" 
          size="sm"
        />
        <UInput 
          v-model="githubPath" 
          placeholder="存储路径 (默认: images/)" 
          size="sm"
        />
      </div>
      
      <!-- SM.MS 配置 -->
      <div v-if="selectedHost === 'smms'" class="space-y-2 mb-3">
        <UInput 
          v-model="smmsToken" 
          placeholder="SM.MS API Token (可选)" 
          size="sm"
          type="password"
        />
      </div>
      
      <!-- 保存配置按钮 - 移到这里，让它在任何配置下都显示 -->
      <div class="flex justify-end mt-2 mb-3">
        <UButton 
          color="green" 
          variant="solid" 
          size="sm" 
          icon="i-heroicons-check" 
          @click="saveConfig"
        >
          保存配置
        </UButton>
      </div>
      
      <div class="flex flex-col items-center justify-center border-2 border-dashed border-gray-500 rounded-lg p-4 cursor-pointer hover:border-blue-400 transition-colors"
           @click="triggerFileInput"
           @dragover.prevent="isDragging = true"
           @dragleave.prevent="isDragging = false"
           @drop.prevent="handleFileDrop"
           :class="{ 'border-blue-400': isDragging }"
      >
        <input
          ref="fileInput"
          type="file"
          accept="image/*"
          @change="handleFileSelect"
          class="hidden"
        />
        <UIcon 
          v-if="!isUploading && !previewUrl" 
          name="i-heroicons-cloud-arrow-up" 
          class="w-10 h-10 text-gray-400 mb-2" 
        />
        <img 
          v-if="previewUrl && !isUploading" 
          :src="previewUrl" 
          class="max-h-32 max-w-full mb-2 rounded" 
        />
        <UProgress 
          v-if="isUploading" 
          :value="uploadProgress" 
          color="blue" 
          class="w-full mb-2" 
        />
        <p v-if="!isUploading && !previewUrl" class="text-sm text-gray-300">
          点击或拖拽图片到此处上传
        </p>
        <p v-if="isUploading" class="text-sm text-gray-300">
          上传中... {{ uploadProgress }}%
        </p>
      </div>
      
      <!-- 添加明确的上传按钮 -->
      <div v-if="previewUrl && !isUploading && !uploadedUrl" class="flex justify-end mt-2">
        <UButton 
          color="blue" 
          variant="solid" 
          size="sm" 
          icon="i-heroicons-cloud-arrow-up" 
          @click="startUpload"
        >
          开始上传
        </UButton>
      </div>
      
      <div v-if="errorMessage" class="mt-2 text-red-400 text-sm">
        {{ errorMessage }}
      </div>
      
      <div v-if="uploadedUrl" class="mt-2">
        <div class="flex items-center gap-2 bg-gray-800 p-2 rounded">
          <input 
            type="text" 
            :value="uploadedUrl" 
            readonly 
            class="bg-transparent flex-1 text-sm outline-none"
          />
          <UButton 
            icon="i-heroicons-clipboard" 
            color="gray" 
            variant="ghost" 
            size="xs" 
            @click="copyToClipboard(uploadedUrl)" 
          />
        </div>
        <div class="flex justify-end mt-2">
          <UButton 
            color="blue" 
            variant="solid" 
            size="sm" 
            @click="insertImage" 
          >
            插入图片
          </UButton>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import { useToast, useRuntimeConfig } from '#imports';

const props = defineProps({
  position: {
    type: Object,
    required: true,
    default: () => ({ x: 0, y: 0 })
  }
});

const emit = defineEmits(['close', 'upload-success']);
const toast = useToast();

// 状态变量
const selectedHost = ref('smms');
const fileInput = ref<HTMLInputElement | null>(null);
const isDragging = ref(false);
const isUploading = ref(false);
const uploadProgress = ref(0);
const previewUrl = ref('');
const uploadedUrl = ref('');
const errorMessage = ref('');

// GitHub 配置
const githubToken = ref('');
const githubRepo = ref('');
const githubBranch = ref('main');
const githubPath = ref('images/');

// SM.MS 配置
const smmsToken = ref('');

// 触发文件选择
const triggerFileInput = () => {
  fileInput.value?.click();
};

// 处理文件拖放
const handleFileDrop = (event: DragEvent) => {
  isDragging.value = false;
  const files = event.dataTransfer?.files;
  if (files && files.length > 0) {
    handleFile(files[0]);
  }
};

// 处理文件选择
const handleFileSelect = (event: Event) => {
  const input = event.target as HTMLInputElement;
  const files = input.files;
  if (files && files.length > 0) {
    handleFile(files[0]);
  }
};

// 添加一个变量存储当前选择的文件
const selectedFile = ref<File | null>(null);

// 处理文件
// 修改文件处理方法
const handleFile = async (file: File) => {
  // 检查文件类型
  if (!file.type.startsWith('image/')) {
    errorMessage.value = '请选择图片文件';
    return;
  }
  
  // 检查文件大小 (限制为 5MB)
  const maxSize = 5 * 1024 * 1024;
  if (file.size > maxSize) {
    errorMessage.value = '图片大小不能超过 5MB';
    return;
  }
  
  // 清除错误信息
  errorMessage.value = '';
  
  // 创建预览
  const reader = new FileReader();
  reader.onload = (e) => {
    previewUrl.value = e.target?.result as string;
  };
  reader.readAsDataURL(file);
  
  // 直接开始上传
  try {
    if (selectedHost.value === 'smms') {
      await uploadToSMMS(file);
    } else if (selectedHost.value === 'github') {
      await uploadToGitHub(file);
    }
  } catch (error: any) {
    errorMessage.value = error.message || '上传失败';
    isUploading.value = false;
  }
};

// 添加复制到剪贴板方法
const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text);
    toast.add({
      title: '成功',
      description: '已复制到剪贴板',
      color: 'green',
      timeout: 2000
    });
  } catch (error) {
    toast.add({
      title: '错误',
      description: '复制失败',
      color: 'red',
      timeout: 2000
    });
  }
};

// 添加开始上传方法
const startUpload = () => {
  if (selectedFile.value) {
    uploadFile(selectedFile.value);
  } else {
    errorMessage.value = '请先选择图片';
  }
};

// 修改 SM.MS 上传方法
const uploadToSMMS = async (file: File) => {
  const formData = new FormData();
  formData.append('smfile', file);
  
  // 显示上传进度
  isUploading.value = true;
  uploadProgress.value = 10;
  
  try {
    // 使用 imgtp.com 作为替代图床 (SM.MS 的替代方案)
    const response = await fetch('https://imgtp.com/api/upload', {
      method: 'POST',
      body: formData
    });
    
    uploadProgress.value = 90;
    
    if (!response.ok) {
      throw new Error(`HTTP 错误: ${response.status}`);
    }
    
    const data = await response.json();
    
    if (data.code === 200) {
      uploadedUrl.value = data.data.url;
      uploadProgress.value = 100;
      isUploading.value = false;
    } else {
      throw new Error(data.msg || '图片上传失败');
    }
  } catch (error: any) {
    console.error('图片上传错误:', error);
    errorMessage.value = `上传失败: ${error.message}`;
    isUploading.value = false;
  }
};

// 插入图片 - 修改为插入 Markdown 格式链接
const insertImage = () => {
  // 创建 Markdown 格式的图片链接
  const markdownLink = `![](${uploadedUrl.value})`;
  emit('upload-success', markdownLink);
  // 关闭弹窗
  emit('close');
};

// 保存配置到 localStorage - 修复保存功能
const saveConfig = () => {
  try {
    const config = {
      selectedHost: selectedHost.value,
      github: {
        token: githubToken.value,
        repo: githubRepo.value,
        branch: githubBranch.value,
        path: githubPath.value
      },
      smms: {
        token: smmsToken.value
      }
    };
    
    localStorage.setItem('imageHostingConfig', JSON.stringify(config));
    
    toast.add({
      title: '成功',
      description: '配置已保存',
      color: 'green',
      timeout: 2000
    });
  } catch (error) {
    console.error('保存配置失败:', error);
    toast.add({
      title: '错误',
      description: '保存配置失败',
      color: 'red',
      timeout: 2000
    });
  }
};

// 从 localStorage 加载配置
const loadConfig = () => {
  const configStr = localStorage.getItem('imageHostingConfig');
  if (configStr) {
    try {
      const config = JSON.parse(configStr);
      selectedHost.value = config.selectedHost || 'smms';
      
      if (config.github) {
        githubToken.value = config.github.token || '';
        githubRepo.value = config.github.repo || '';
        githubBranch.value = config.github.branch || 'main';
        githubPath.value = config.github.path || 'images/';
      }
      
      if (config.smms) {
        smmsToken.value = config.smms.token || '';
      }
    } catch (error) {
      console.error('加载配置失败:', error);
    }
  }
};

// 监听配置变化并自动保存
watch([selectedHost, githubToken, githubRepo, githubBranch, githubPath, smmsToken], 
  () => {
    saveConfig();
  },
  { deep: true }
);

// 组件挂载时加载配置
onMounted(() => {
  loadConfig();
});
</script>