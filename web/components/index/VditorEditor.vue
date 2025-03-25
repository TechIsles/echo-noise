<template>
  <div ref="editorContainer" class="vditor-container"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from "vue";
import Vditor from "vditor";
import "vditor/dist/index.css";

const props = defineProps({
  modelValue: {
    type: String,
    default: "",
  },
});

const emit = defineEmits(["update:modelValue"]);

const editorContainer = ref<HTMLElement>();
let vditorInstance: Vditor | null = null;

const editorOptions: IOptions = {
  mode: "ir",
  height: "auto",
  minHeight: 140,
  resize: {
    enable: true,
    position: 'bottom'
  },
  icon: "ant",
  lang: "zh_CN" as keyof II18n,
  theme: "classic",
  toolbar: [
    "emoji",
    "headings", // 标题选项
    "bold",
    "italic",
    "strike",
    "link",
    "|",
    "list",
    "ordered-list",
    "check",
    "|",
    "quote",
    "line",
    "code",
    "inline-code",
    "|",
    "preview",
    "fullscreen"
  ],
  toolbarConfig: {
    pin: true,
  },
  counter: {
    enable: false,
  },
  cache: {
    enable: true,
    id: "vue-vditor",
  },
  input: (content: string) => {
    emit("update:modelValue", content);
  },
  preview: {
    hljs: {
      style: "native",
    },
    actions: [],
  },
  placeholder: "灵感记录~"
};

onMounted(async () => {
  if (!editorContainer.value) return;

  vditorInstance = new Vditor(editorContainer.value, {
    ...editorOptions,
    after: () => {
      vditorInstance?.setValue(props.modelValue);
    },
  });
});

onBeforeUnmount(() => {
  vditorInstance?.destroy();
});

defineExpose({
  clear: () => {
    if (vditorInstance) {
      vditorInstance.setValue('');
    }
  },
});
</script>

<style>
.vditor-container {
  border-radius: 8px;
  margin-bottom: 12px;
  transition: all 0.3s ease;
  border: 1px solid #e9ecef;
  position: relative; /* 添加相对定位 */
}
.vditor-content {
  position: relative;
  z-index: 1;
}
.vditor-container:hover {
  border-color: #90a4ae;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.vditor-toolbar {
  display: flex;
  flex-wrap: nowrap;
  overflow-x: auto;
  white-space: nowrap;
  -webkit-overflow-scrolling: touch;
  scrollbar-width: none;
  -ms-overflow-style: none;
  background-color: #f8f9fab7;
  border-bottom: 1px solid #e9ecef;
  position: relative; /* 改为相对定位 */
  z-index: 2; /* 调整层级 */
}

.vditor-toolbar::-webkit-scrollbar {
  display: none; /* Chrome, Safari and Opera */
}

.vditor-toolbar--pin {
  padding-left: 6px !important;
  background-color: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
}

/* 修改弹出面板样式 */
.vditor-panel--none {
  display: none !important;
}

.vditor-panel {
  position: fixed; /* 改为固定定位 */
  z-index: 9999; /* 确保在工具栏之上 */
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border-radius: 4px;
  border: 1px solid #e9ecef;
}
.vditor-hint {
  position: fixed; /* 改为固定定位 */
  z-index: 999;
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border-radius: 4px;
  border: 1px solid #e9ecef;
}
.vditor-toolbar__item {
  flex-shrink: 0;
  padding: 6px !important;
  transition: all 0.2s ease;
}

.vditor-toolbar__item:hover {
  background-color: #e9ecef;
  border-radius: 4px;
}

.vditor-ir pre.vditor-reset {
  padding: 10px 16px !important;
  color: #1a2634 !important;
  line-height: 1.6;
  font-size: 14px;
}

.vditor-ir pre.vditor-reset:empty:before {
  color: #90a4ae !important;
}

.vditor-resize {
  cursor: ns-resize;
  user-select: none;
  height: 3px;
  background-color: #f1f3f5;
}

.vditor-resize:hover {
  background-color: #e9ecef;
}

@media screen and (max-width: 520px) {
  .vditor-toolbar__item {
    padding: 4px !important;
  }
  
  .vditor-ir pre.vditor-reset {
    padding: 8px 12px !important;
    font-size: 13px;
  }
}
</style>
