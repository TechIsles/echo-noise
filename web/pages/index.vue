<template>
  <div class="background-container" :style="backgroundStyle">
    <div class="loading" v-if="!isLoaded">
      <div class="rainbow-spinner"></div>
      <div class="loading-text">加载中...</div>
    </div>
    <div class="content-wrapper">
      <UContainer class="container-fixed py-2 pb-4 my-4">
        <div class="moments-header">
          <div class="header-image" :style="headerImageStyle">
            <h1 class="header-title">Noise的说说笔记</h1>
            <div class="profile-info">
              <img class="avatar" 
                   @click="changeBackground" 
                   src="https://s2.loli.net/2025/03/24/HnSXKvibAQlosIW.png" 
                   alt="Noise">
              <div class="profile-text">
                <div class="title">Noise</div>
                <div class="description">执迷不悟</div>
              </div>
            </div>
          </div>
        </div>
        <AddForm />
        <MessageList />
      </UContainer>
      <Notification />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import AddForm from '@/components/index/AddForm.vue'
import MessageList from '@/components/index/MessageList.vue'
import Notification from '~/components/widgets/Notification.vue';

const backgroundImages = [
  'https://files.catbox.moe/8lsc79.webp',
  'https://files.catbox.moe/vql8zy.webp',
  'https://files.catbox.moe/1n77zd.webp',
  'https://files.catbox.moe/9539a2.webp',
  'https://files.catbox.moe/dttmse.webp',
  'https://files.catbox.moe/wkwh2y.webp',
  'https://files.catbox.moe/1n77zd.webp',
];

const currentImage = ref('')
const isLoaded = ref(false)

const changeBackground = () => {
  const newIndex = Math.floor(Math.random() * backgroundImages.length)
  const newImage = backgroundImages[newIndex]
  if (newImage === currentImage.value) {
    changeBackground()
    return
  }
  const img = new Image()
  img.src = newImage
  img.onload = () => {
    currentImage.value = newImage
  }
}

onMounted(() => {
  currentImage.value = backgroundImages[Math.floor(Math.random() * backgroundImages.length)]
  const img = new Image()
  img.src = currentImage.value
  img.onload = () => {
    isLoaded.value = true
  }
})

const backgroundStyle = computed(() => ({
  '--bg-image': `url(${currentImage.value})`
}))

const headerImageStyle = computed(() => ({
  backgroundImage: `url(${currentImage.value})`
}))
</script>

<style>
html, body {
  margin: 0;
  padding: 0;
  width: 100%;
  overflow-x: hidden;
  min-height: 100vh;
  overflow-y: overlay;
}

.background-container {
  min-height: 100vh;
  width: 100vw;
  position: relative;
  overflow-x: hidden;
}

.background-container::before {
  content: '';
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  width: 100vw;
  height: 100vh;
  background-image: var(--bg-image);
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  filter: blur(8px);
  z-index: -1;
  transition: background-image 0.5s ease;
}

.content-wrapper {
  position: relative;
  width: 100%;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
  padding-right: calc(100vw - 100%);
  overflow-x: hidden;
}

.container-fixed {
  backdrop-filter: blur(4px);
  border-radius: 8px;
  margin: 0 auto;
  max-width: 1200px;
  width: 100%;
  position: relative;
  z-index: 1;
  box-sizing: border-box;
  overflow: hidden;
}

.moments-header {
  margin-bottom: 20px;
}

.header-image {
  position: relative;
  width: 100%;
  height: 300px;
  background-size: cover;
  background-position: center;
  border-radius: 8px;
  overflow: hidden;
  transition: background-image 0.5s ease;
}

.header-title {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: white;
  font-size: 2.5rem;
  font-weight: bold;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
  margin: 0;
  white-space: nowrap;
  transition: font-size 0.3s ease;
}

@media screen and (max-width: 768px) {
  .header-title {
    font-size: 1.8rem;
  }
}

.profile-info {
  position: absolute;
  bottom: 20px;
  right: 20px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  border: 2px solid white;
  object-fit: cover;
  cursor: pointer;
  transition: transform 0.3s ease;
}

.avatar:hover {
  transform: scale(1.1);
}

.profile-text {
  text-align: right;
  color: white;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.5);
}

.profile-text .title {
  font-size: 1.2rem;
  font-weight: bold;
}

.profile-text .description {
  font-size: 0.9rem;
  opacity: 0.9;
}

.u-container {
  backdrop-filter: blur(4px);
  border-radius: 8px;
  margin: 0 auto;
  max-width: 1200px;
  width: 100%;
  position: relative;
  z-index: 1;
  box-sizing: border-box;
  overflow-x: hidden;
}

.loading {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background: rgba(255, 255, 255, 0.8);
  z-index: 1000;
  gap: 15px;
}

.loading-text {
  font-size: 16px;
  color: #333;
}

.rainbow-spinner {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  border: 4px solid transparent;
  border-top: 4px solid #FF0000;
  border-right: 4px solid #00FF00;
  border-bottom: 4px solid #0000FF;
  border-left: 4px solid #FF00FF;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>