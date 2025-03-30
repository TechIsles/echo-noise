<template>
    <div class="fixed inset-0 bg-black overflow-auto">
        <div class="min-h-screen p-4 flex items-center justify-center">
            <div class="w-[800px] max-w-[95%] bg-[#1a1b2e]/80 backdrop-blur-md rounded-lg shadow-xl p-6">
                <h1 class="text-3xl font-bold text-center text-white mb-8">系统管理面板</h1>

                <!-- 系统状态卡片 -->
                <div class="bg-gray-700 rounded-lg p-4 mb-6">
                    <h2 class="text-xl font-semibold text-white mb-4">系统状态</h2>
                    <div class="grid gap-4">
                        <div class="flex justify-between items-center">
                            <span class="text-gray-300">系统管理员</span>
                            <span class="text-white font-medium">{{ userStore?.status?.username }}</span>
                        </div>
                        <div class="flex justify-between items-center">
                            <span class="text-gray-300">当前用户</span>
                            <span class="text-white font-medium">
                                {{ isLogin ? userStore.user?.username : "未登录" }}
                            </span>
                        </div>
                     
                        <div class="flex justify-between items-center">
                            <span class="text-gray-300">笔记总数</span>
                            <span class="text-white font-medium">{{ userStore?.status?.total_messages }} 条</span>
                        </div>
                    </div>
                </div>

                               <!-- 网站配置区域 -->
                <div v-if="isAdmin" class="bg-gray-700 rounded-lg p-4 mb-6">
                    <div class="flex justify-between items-center mb-4">
                        <h2 class="text-xl font-semibold text-white">网站配置</h2>
                    </div>
                    

                    <!-- 配置展示/编辑表单 -->
                    <div class="space-y-4">
                        <div v-for="(label, key) in configLabels" :key="key" class="bg-gray-800 rounded p-3">
                            <div class="flex justify-between items-center mb-2">
                                <span class="text-gray-300">{{ label }}</span>
                                <UButton
                                    size="sm"
                                    @click="editItem[key] = !editItem[key]"
                                    :color="editItem[key] ? 'gray' : 'primary'"
                                    variant="soft"
                                >
                                    {{ editItem[key] ? '取消' : '编辑' }}
                                </UButton>
                            </div>
                            
                            <div v-if="editItem[key]">
                                <template v-if="key === 'backgrounds'">
                                    <div class="space-y-2">
                                        <div v-for="(bg, index) in frontendConfig.backgrounds" 
                                             :key="index" 
                                             class="flex gap-2">
                                            <UInput v-model="frontendConfig.backgrounds[index]" 
                                                   placeholder="输入图片URL" 
                                                   class="flex-1" />
                                            <UButton @click="removeBackground(index)" 
                                                    icon="i-heroicons-trash" 
                                                    color="red" 
                                                    variant="ghost" />
                                        </div>
                                        <div class="flex gap-2">
                                            <UButton @click="addBackground" 
                                                    icon="i-heroicons-plus" 
                                                    variant="ghost" 
                                                    class="mr-2">
                                                添加链接
                                            </UButton>
                                            <UButton @click="triggerFileInput" 
                                                    icon="i-heroicons-cloud-arrow-up" 
                                                    variant="ghost">
                                                上传图片
                                            </UButton>
                                        </div>
                                    </div>
                                </template>
                                <template v-else-if="key === 'subtitleText'">
                                    <UTextarea
                                        v-model="frontendConfig[key]"
                                        :placeholder="`输入${label}`"
                                        class="w-full mb-2"
                                    />
                                </template>
                                <template v-else>
                                    <UInput
                                        v-model="frontendConfig[key]"
                                        :placeholder="`输入${label}`"
                                        class="w-full mb-2"
                                    />
                                </template>
                                <div class="flex justify-end gap-2">
                                    <UButton @click="resetConfigItem(key)" variant="ghost" color="gray">
                                        重置
                                    </UButton>
                                    <UButton @click="saveConfigItem(key)" color="primary">
                                        保存
                                    </UButton>
                                </div>
                            </div>
                            <div v-else>
                                <template v-if="key === 'backgrounds'">
                                    <div class="grid grid-cols-3 gap-2">
                                        <img v-for="(bg, index) in frontendConfig.backgrounds"
                                             :key="index"
                                             :src="bg"
                                             class="w-full h-20 object-cover rounded cursor-pointer"
                                             @click="previewImage(bg)" />
                                    </div>
                                </template>
                                <template v-else>
                                    <p class="text-white break-words">{{ frontendConfig[key] }}</p>
                                </template>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- 保存按钮 -->
                <div v-if="editMode" class="flex justify-end gap-2 mb-6">
                    <UButton @click="resetConfig" variant="ghost" color="gray">
                        重置
                    </UButton>
                    <UButton @click="saveConfig" color="primary">
                        保存所有更改
                    </UButton>
                </div>

                <!-- 底部操作栏 -->
                <div class="flex justify-between items-center">
                    <UButton
                        icon="i-heroicons-arrow-left"
                        variant="ghost"
                        color="white"
                        @click="$router.push('/')"
                        class="text-white hover:text-black"
                    >
                        返回首页
                    </UButton>
                    <div v-if="isLogin">
                        <UButton
                            icon="i-heroicons-power"
                            color="red"
                            variant="ghost"
                            @click="logout"
                        >
                            退出登录
                        </UButton>
                    </div>
                    <div v-else>
                        <UButton
                            color="primary"
                            @click="showLoginModal = true"
                        >
                            登录
                        </UButton>
                    </div>
                </div>
            </div>
        </div>
   

        <!-- 登录模态框 -->
        <UModal v-model="showLoginModal">
            <div class="bg-gray-800 p-6 rounded-lg">
                <h3 class="text-xl font-semibold text-white mb-4">
                    {{ authmode ? '用户登录' : '用户注册' }}
                </h3>
                <UForm :state="authForm" class="space-y-4">
                    <UFormGroup>
                        <UInput
                            v-model="authForm.username"
                            placeholder="用户名"
                            class="w-full"
                        />
                    </UFormGroup>
                    <UFormGroup>
                        <UInput
                            v-model="authForm.password"
                            type="password"
                            placeholder="密码"
                            class="w-full"
                        />
                    </UFormGroup>
                    <div class="flex justify-between items-center">
                        <UButton
                            variant="ghost"
                            @click="authmode = !authmode"
                        >
                            {{ authmode ? '去注册' : '去登录' }}
                        </UButton>
                        <UButton
                            color="primary"
                            @click="authmode ? login(authForm) : register(authForm)"
                        >
                            {{ authmode ? '登录' : '注册' }}
                        </UButton>
                    </div>
                </UForm>
            </div> 
        </UModal>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed, onMounted } from 'vue'
import type { UserToLogin, UserToRegister } from '~/types/models'
import { useUser } from '~/composables/useUser'
import { useUserStore } from '~/store/user'

const userStore = useUserStore()
const { login, register, logout } = useUser()

// 状态变量
const isLogin = computed(() => userStore?.isLogin ?? false)
const isAdmin = computed(() => userStore?.user?.IsAdmin ?? false)
const authmode = ref(true)
const showLoginModal = ref(false)
const editMode = ref(false)
const fileInput = ref<HTMLInputElement | null>(null) 

// 配置相关
const configLabels = {
    siteTitle: '站点标题',
    subtitleText: '欢迎语',
    avatarURL: '头像链接',
    username: '用户名',
    description: '个人描述',
    backgrounds: '背景图片'
}

const frontendConfig = reactive({
    siteTitle: '',
    subtitleText: '',
    avatarURL: '',
    username: '',
    description: '',
    backgrounds: [] as string[]
})

const authForm = reactive<UserToLogin | UserToRegister>({
    username: '',
    password: ''
})

const editItem = reactive({
    siteTitle: false,
    subtitleText: false,
    avatarURL: false,
    username: false,
    description: false,
    backgrounds: false,
})

// 更新默认配置
const defaultConfig = {
    siteTitle: 'Noise的说说笔记',
    subtitleText: '欢迎访问，点击头像可更换封面背景！',
    avatarURL: 'https://s2.loli.net/2025/03/24/HnSXKvibAQlosIW.png',
    username: 'Noise',
    description: '执迷不悟',
    backgrounds: [
        "https://s2.loli.net/2025/03/27/KJ1trnU2ksbFEYM.jpg",
        "https://s2.loli.net/2025/03/27/MZqaLczCvwjSmW7.jpg",
        "https://s2.loli.net/2025/03/27/UMijKXwJ9yTqSeE.jpg",
        "https://s2.loli.net/2025/03/27/WJQIlkXvBg2afcR.jpg",
        "https://s2.loli.net/2025/03/27/oHNQtf4spkq2iln.jpg",
        "https://s2.loli.net/2025/03/27/PMRuX5loc6Uaimw.jpg",
        "https://s2.loli.net/2025/03/27/U2WIslbNyTLt4rD.jpg",
        "https://s2.loli.net/2025/03/27/xu1jZL5Og4pqT9d.jpg",
        "https://s2.loli.net/2025/03/27/OXqwzZ6v3PVIns9.jpg",
        "https://s2.loli.net/2025/03/27/HGuqlE6apgNywbh.jpg",
        "https://s2.loli.net/2025/03/26/d7iyuPYA8cRqD1K.jpg",
        "https://s2.loli.net/2025/03/27/wYy12qDMH6bGJOI.jpg",
        "https://s2.loli.net/2025/03/27/y67m2k5xcSdTsHN.jpg",
    ]
}
// 添加单个配置项保存方法
const saveConfigItem = async (key: string) => {
    try {
        const response = await fetch('/api/settings', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            },
            body: JSON.stringify({
                allowRegistration: true,
                frontendSettings: frontendConfig
            })
        })
        
        const data = await response.json()
        if (data.code === 1) {
            editItem[key] = false
            useToast().add({
                title: '成功',
                description: `${configLabels[key]}已更新`,
                color: 'green'
            })
        } else {
            throw new Error(data.msg)
        }
    } catch (error) {
        useToast().add({
            title: '错误',
            description: error.message || '配置更新失败',
            color: 'red'
        })
    }
}

// 添加单个配置项重置方法
const resetConfigItem = (key: string) => {
    frontendConfig[key] = defaultConfig[key]
    editItem[key] = false
}
// 修改 fetchConfig 方法
const fetchConfig = async () => {
    try {
        const response = await fetch('/api/frontend/config')
        const data = await response.json()
        if (data.code === 1 && data.frontendSettings) {
            Object.assign(frontendConfig, {
                siteTitle: data.frontendSettings.siteTitle || defaultConfig.siteTitle,
                subtitleText: data.frontendSettings.subtitleText || defaultConfig.subtitleText,
                avatarURL: data.frontendSettings.avatarURL || defaultConfig.avatarURL,
                username: data.frontendSettings.username || defaultConfig.username,
                description: data.frontendSettings.description || defaultConfig.description,
                backgrounds: Array.isArray(data.frontendSettings.backgrounds) && data.frontendSettings.backgrounds.length > 0
                    ? [...data.frontendSettings.backgrounds]
                    : [...defaultConfig.backgrounds]
            })
        } else {
            // 如果获取失败，使用默认配置
            Object.assign(frontendConfig, defaultConfig)
        }
    } catch (error) {
        console.error('获取配置失败:', error)
        // 发生错误时也使用默认配置
        Object.assign(frontendConfig, defaultConfig)
    }
}

const saveConfig = async () => {
    try {
        const response = await fetch('/api/settings', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            },
            body: JSON.stringify({
                allowRegistration: true,
                frontendSettings: frontendConfig
            })
        })
        
        const data = await response.json()
        if (data.code === 1) {
            editMode.value = false
            useToast().add({
                title: '成功',
                description: '配置已更新',
                color: 'green'
            })
        } else {
            throw new Error(data.msg)
        }
    } catch (error) {
        useToast().add({
            title: '错误',
            description: error.message || '配置更新失败',
            color: 'red'
        })
    }
}

const resetConfig = () => {
    fetchConfig()
    editMode.value = false
}

const addBackground = () => {
    frontendConfig.backgrounds.push('')
}

const removeBackground = (index: number) => {
    frontendConfig.backgrounds.splice(index, 1)
}

const triggerFileInput = () => {
    fileInput.value?.click()
}

const handleFileUpload = async (event: Event) => {
    const files = (event.target as HTMLInputElement).files
    if (!files) return

    for (const file of Array.from(files)) {
        try {
            const formData = new FormData()
            formData.append('image', file)

            const response = await fetch('/api/images/upload', {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                },
                body: formData
            })

            const data = await response.json()
            if (data.code === 1 && data.data.url) {
                frontendConfig.backgrounds.push(data.data.url)
            }
        } catch (error) {
            console.error('上传失败:', error)
        }
    }

    if (fileInput.value) {
        fileInput.value.value = ''
    }
}

const previewImage = (url: string) => {
    window.open(url, '_blank')
}

// 监听器
watch(() => userStore.isLogin, (newVal) => {
    if (!newVal) {
        userStore.getStatus()
        userStore.getUser()
        userStore.$reset()
    }
})

// 生命周期
onMounted(async () => {
    await userStore.getStatus()
    await userStore.getUser()
    await fetchConfig()
})
</script>

<style scoped>
.hidden {
    display: none;
}
</style>