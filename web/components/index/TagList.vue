<template>
  <div class="tags-container mx-auto w-full sm:max-w-2xl">
    <div class="tags-wrapper">
      <div class="tags-scroll">
        <span
          v-for="tag in tags"
          :key="tag.name"
          class="tag-item"
          @click="handleTagClick(tag.name)"
        >
          #{{ tag.name }}
          <span class="tag-count">({{ tag.count }})</span>
        </span>
      </div>
    </div>
    <div class="scroll-fade"></div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const emit = defineEmits(['tagClick'])

const props = defineProps({
  tags: {
    type: Array,
    default: () => []
  }
})
// 过滤掉非正常标签
const filteredTags = computed(() => {
  return props.tags.filter(tag => {
    // 过滤掉包含特殊字符的标签
    const invalidChars = /[/?=&]/;
    // 过滤掉数字ID和媒体链接
    const isMediaLink = /^(song|video|playlist)\?id=\d+$/;
    return !invalidChars.test(tag.name) && !isMediaLink.test(tag.name);
  });
});
const handleTagClick = (tagName: string) => {
  emit('tagClick', tagName)
}
</script>

<style scoped>
.tags-container {
  width: 100%;
  margin: 0.5rem auto;
  padding: 0 0.5rem;
  position: relative;
  background: transparent;
}

.tags-wrapper {
  position: relative;
  overflow: hidden;
  background: transparent;
}

.tags-scroll {
  display: flex;
  flex-wrap: nowrap;
  gap: 0.5rem;
  overflow-x: auto;
  padding: 0.5rem 0;
  -webkit-overflow-scrolling: touch;
  scrollbar-width: none;
  -ms-overflow-style: none;
  background: transparent;
}

.tags-scroll::-webkit-scrollbar {
  display: none;
}

.tag-item {
  display: inline-flex;
  align-items: center;
  padding: 0.25rem 0.75rem;
  color: #c7cace;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 0.875rem;
  white-space: nowrap;
  flex-shrink: 0;
}

.tag-item:hover {
  color: #fb923c;
  transform: translateY(-1px);
}

.tag-count {
  margin-left: 0.25rem;
  font-size: 0.75rem;
  opacity: 0.8;
}

.scroll-fade {
  position: absolute;
  right: 0;
  top: 0;
  height: 100%;
  width: 32px;
  background: linear-gradient(to left, transparent, transparent);
  pointer-events: none;
}
</style>