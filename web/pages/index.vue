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
    <h1 class="header-title">{{ frontendConfig.siteTitle }}</h1>
    <div ref="subtitleEl" class="header-subtitle">{{ frontendConfig.subtitleText }}</div>
    
    <div class="profile-info">
        <img class="avatar" 
             @click="changeBackground" 
             :src="frontendConfig.avatarURL" 
             :alt="frontendConfig.username">
        <div class="profile-text">
            <div class="title">{{ frontendConfig.username }}</div>
            <div class="description">{{ frontendConfig.description }}</div>
        </div>
    </div>
          </div>
        </div>
        <AddForm />
        <!-- 确保 MessageList 有足够的 z-index -->
        <MessageList class="message-list-container" />
      </UContainer>
      <Notification />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'  // 添加 onUnmounted
import AddForm from '@/components/index/AddForm.vue'
import MessageList from '@/components/index/MessageList.vue'
import Notification from '~/components/widgets/Notification.vue';

// 添加前端配置的响应式对象
const frontendConfig = ref({
    siteTitle: '',
    subtitleText: '',
    avatarURL: '',
    username: '',
    description: '',
    backgrounds: [] as string[]
})

const backgroundStyle = computed(() => ({
    '--bg-image': `url(${currentImage.value || frontendConfig.value.backgrounds[0]})`
}))
// 添加 headerImageStyle 计算属性
const headerImageStyle = computed(() => ({
    'background-image': `url(${currentImage.value || frontendConfig.value.backgrounds[0]})`,
    'background-size': 'cover',
    'background-position': 'center'
}))
// 修改 fetchConfig 方法

const fetchConfig = async () => {
    try {
        const response = await fetch('/api/frontend/config');
        if (!response.ok) throw new Error('Network response was not ok');
        
        const data = await response.json();
        console.log('获取到的配置数据:', data); // 添加调试日志
        
        if (data && data.frontendSettings) {
            const settings = data.frontendSettings;
            frontendConfig.value = {
                siteTitle: settings.siteTitle || '默认标题',
                subtitleText: settings.subtitleText || '默认副标题',
                avatarURL: settings.avatarURL || 'https://s2.loli.net/2025/03/24/HnSXKvibAQlosIW.png',
                username: settings.username || 'Noise',
                description: settings.description || '执迷不悟',
                backgrounds: settings.backgrounds || [
                'https://s2.loli.net/2025/03/27/KJ1trnU2ksbFEYM.jpg',
                'https://s2.loli.net/2025/03/27/MZqaLczCvwjSmW7.jpg',
                'https://s2.loli.net/2025/03/27/UMijKXwJ9yTqSeE.jpg',
                'https://s2.loli.net/2025/03/27/WJQIlkXvBg2afcR.jpg',
                'https://s2.loli.net/2025/03/27/oHNQtf4spkq2iln.jpg',
                'https://s2.loli.net/2025/03/27/PMRuX5loc6Uaimw.jpg',
                'https://s2.loli.net/2025/03/27/U2WIslbNyTLt4rD.jpg',
                'https://s2.loli.net/2025/03/27/xu1jZL5Og4pqT9d.jpg',
                'https://s2.loli.net/2025/03/27/OXqwzZ6v3PVIns9.jpg',
                'https://s2.loli.net/2025/03/27/HGuqlE6apgNywbh.jpg',
                'https://s2.loli.net/2025/03/26/d7iyuPYA8cRqD1K.jpg',
                'https://s2.loli.net/2025/03/27/wYy12qDMH6bGJOI.jpg',
                'https://s2.loli.net/2025/03/27/y67m2k5xcSdTsHN.jpg',
                ]
            };
            
            // 设置初始背景图并标记为已加载
            if (frontendConfig.value.backgrounds.length > 0) {
                currentImage.value = frontendConfig.value.backgrounds[0];
                const img = new Image();
                img.src = currentImage.value;
                img.onload = () => {
                    isLoaded.value = true;
                };
            } else {
                isLoaded.value = true; // 如果没有背景图也要设置为已加载
            }
        } else {
            throw new Error('Invalid response format');
        }
    } catch (error) {
        console.error('获取配置失败:', error);
        // 设置默认值
        frontendConfig.value = {
            siteTitle: 'Noise的说说笔记',
            subtitleText: '欢迎访问，点击头像可更换封面背景！',
            avatarURL: 'https://s2.loli.net/2025/03/24/HnSXKvibAQlosIW.png',
            username: 'Noise',
            description: '执迷不悟',
            backgrounds: [
                'https://s2.loli.net/2025/03/27/KJ1trnU2ksbFEYM.jpg',
                'https://s2.loli.net/2025/03/27/MZqaLczCvwjSmW7.jpg',
                'https://s2.loli.net/2025/03/27/UMijKXwJ9yTqSeE.jpg',
                'https://s2.loli.net/2025/03/27/WJQIlkXvBg2afcR.jpg',
                'https://s2.loli.net/2025/03/27/oHNQtf4spkq2iln.jpg',
                'https://s2.loli.net/2025/03/27/PMRuX5loc6Uaimw.jpg',
                'https://s2.loli.net/2025/03/27/U2WIslbNyTLt4rD.jpg',
                'https://s2.loli.net/2025/03/27/xu1jZL5Og4pqT9d.jpg',
                'https://s2.loli.net/2025/03/27/OXqwzZ6v3PVIns9.jpg',
                'https://s2.loli.net/2025/03/27/HGuqlE6apgNywbh.jpg',
                'https://s2.loli.net/2025/03/26/d7iyuPYA8cRqD1K.jpg',
                'https://s2.loli.net/2025/03/27/7Zck3y6XTzhYPs5.jpg',
                'https://s2.loli.net/2025/03/27/y67m2k5xcSdTsHN.jpg',
            ]
        };
        isLoaded.value = true; // 错误时也要设置为已加载
    }
};


const currentImage = ref('')
const isLoaded = ref(false)


const changeBackground = () => {
    const newIndex = Math.floor(Math.random() * frontendConfig.value.backgrounds.length)
    const newImage = frontendConfig.value.backgrounds[newIndex]
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


const subtitleEl = ref<HTMLElement | null>(null)

// 修改打字效果函数
const startTypeEffect = () => {
  if (!subtitleEl.value) return
  
  let index = 0
  let isDeleting = false
  let isWaiting = false
  
  const typeInterval = setInterval(() => {
    if (isWaiting) return

    if (!isDeleting) {
      // 打字过程
      subtitleEl.value!.textContent = frontendConfig.value.subtitleText.slice(0, index + 1)
      index++
      
      if (index >= frontendConfig.value.subtitleText.length) {
        isWaiting = true
        setTimeout(() => {
          isDeleting = true
          isWaiting = false
        }, 2000)
      }
    } else {
      // 删除过程
      index--
      subtitleEl.value!.textContent = frontendConfig.value.subtitleText.slice(0, index)
      
      if (index <= 0) {
        isWaiting = true
        subtitleEl.value!.textContent = ''
        setTimeout(() => {
          isDeleting = false
          isWaiting = false
          index = 0
        }, 1000)
      }
    }
  }, 100)

  return typeInterval
}

onMounted(async () => {
    await fetchConfig();
    
    if (frontendConfig.value.backgrounds.length > 0) {
        currentImage.value = frontendConfig.value.backgrounds[Math.floor(Math.random() * frontendConfig.value.backgrounds.length)]
        const img = new Image()
        img.src = currentImage.value
        img.onload = () => {
            isLoaded.value = true
        }
    }
    const typeInterval = startTypeEffect()
    onUnmounted(() => {
        if (typeInterval) {
            clearInterval(typeInterval)
        }
    })
})
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
.header-subtitle {
  position: absolute;
  top: calc(50% + 60px);  /* 增加间距 */
  left: 50%;
  transform: translate(-50%, -50%);
  color: white;
  font-size: 1rem;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.5);
  white-space: nowrap;
}

@media screen and (max-width: 768px) {
  .header-subtitle {
    font-size: 0.9rem;
    top: calc(50% + 45px);  /* 保持移动端的适配 */
  }
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
  border-radius: 18px;
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
  flex-direction: row-reverse;  /* 改变方向，头像在右侧 */
  align-items: center;
  gap: 10px;
  max-width: 80%;  /* 限制最大宽度 */
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
  text-align: left;  /* 改为左对齐 */
  min-width: 0;  /* 允许内容收缩 */
  overflow-x: auto;  /* 允许横向滚动 */
  scrollbar-width: none;  /* 隐藏滚动条 (Firefox) */
  -ms-overflow-style: none;  /* 隐藏滚动条 (IE/Edge) */
  padding: 5px 0;
}
.profile-text .title {
  font-size: 1.2rem;
  font-weight: bold;
}

.profile-text .description {
  font-size: 0.9rem;
  opacity: 0.9;
}
.profile-text::-webkit-scrollbar {
  display: none;
}

.profile-text .title {
font-size: 1.2rem;
font-weight: bold;
color: #fcfafb;  
text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
white-space: nowrap;  /* 防止换行 */

}

.profile-text .description {
  font-size: 0.9rem;
  color: #fcfafb; 
  text-shadow: 1px 1px 3px rgba(0, 0, 0, 0.5);
  white-space: nowrap;
  opacity: 0.95;
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
/* 添加新的样式 */
.message-list-container {
  position: relative;
  z-index: 10;
}

.container-fixed {
  backdrop-filter: blur(4px);
  border-radius: 8px;
  margin: 0 auto;
  max-width: 1200px;
  width: 100%;
  position: relative;
  /* 降低 container 的 z-index，确保不会遮挡评论框 */
  z-index: 1;
  box-sizing: border-box;
  overflow: visible; /* 修改为 visible，允许评论框溢出 */
}

/* 确保背景不会遮挡评论框 */
.background-container::before {
  z-index: -1;
}
</style>