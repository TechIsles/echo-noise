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
                  <!-- 用户信息配置面板 -->
                <div v-if="isLogin" class="bg-gray-700 rounded-lg p-4 mb-6">
                    <h2 class="text-xl font-semibold text-white mb-4">用户信息配置</h2>
                    <div class="space-y-4">
                        <!-- 用户名修改 -->
                        <div class="bg-gray-800 rounded p-3">
                            <div class="flex justify-between items-center mb-2">
                                <span class="text-gray-300">用户名</span>
                                <UButton
                                    size="sm"
                                    @click="editUserInfo.username = !editUserInfo.username"
                                    :color="editUserInfo.username ? 'gray' : 'primary'"
                                    variant="soft"
                                >
                                    {{ editUserInfo.username ? '取消' : '编辑' }}
                                </UButton>
                            </div>
                            <div v-if="editUserInfo.username">
                                <UInput
                                    v-model="userForm.username"
                                    placeholder="新用户名"
                                    class="w-full mb-2"
                                />
                                <div class="flex justify-end gap-2">
                                    <UButton @click="updateUsername" color="primary">
                                        保存
                                    </UButton>
                                </div>
                            </div>
                            <div v-else>
                                <p class="text-white">{{ userStore.user?.username }}</p>
                            </div>
                        </div>
                         <!-- 在用户信息配置面板中添加 -->
<div class="bg-gray-800 rounded p-3">
    <div class="flex justify-between items-center mb-2">
        <span class="text-gray-300">API Token</span>
        <UButton
            size="sm"
            @click="regenerateToken"
            color="primary"
            variant="soft"
        >
            重新生成
        </UButton>
    </div>
    <div v-if="userToken" class="mb-2">
        <div class="flex items-center gap-2">
            <UInput
                v-model="userToken"
                readonly
                class="font-mono text-sm"
            />
            <UButton
                icon="i-heroicons-clipboard"
                color="gray"
                variant="ghost"
                @click="copyToken"
            />
        </div>
        <p class="text-xs text-gray-400 mt-1">请妥善保管此 Token，它用于 API 访问认证</p>
    </div>
    <div v-else>
        <p class="text-gray-400">暂无 Token</p>
    </div>
</div>
                        <!-- 密码修改 -->
                        <div class="bg-gray-800 rounded p-3">
                            <div class="flex justify-between items-center mb-2">
                                <span class="text-gray-300">修改密码</span>
                                <UButton
                                    size="sm"
                                    @click="editUserInfo.password = !editUserInfo.password"
                                    :color="editUserInfo.password ? 'gray' : 'primary'"
                                    variant="soft"
                                >
                                    {{ editUserInfo.password ? '取消' : '编辑' }}
                                </UButton>
                            </div>
                            <div v-if="editUserInfo.password">
                                <UInput
                                    v-model="userForm.oldPassword"
                                    type="password"
                                    placeholder="当前密码"
                                    class="w-full mb-2"
                                />
                                <UInput
                                    v-model="userForm.newPassword"
                                    type="password"
                                    placeholder="新密码"
                                    class="w-full mb-2"
                                />
                                <div class="flex justify-end gap-2">
                                    <UButton @click="updatePassword" color="primary">
                                        保存
                                    </UButton>
                                </div>
                            </div>
                        </div>

                        <!-- 管理员权限设置（仅管理员可见） -->
                        <div v-if="isAdmin" class="bg-gray-800 rounded p-3">
                            <h3 class="text-white font-semibold mb-3">用户权限管理</h3>
                            <div class="space-y-2">
                                <div v-for="user in userStore?.status?.users" :key="user.userID" 
                                     class="flex justify-between items-center">
                                    <span class="text-gray-300">{{ user.userName }}</span>
                                    <UButton
                                        v-if="user.userID !== 1"
                                        size="sm"
                                        :color="user.isAdmin ? 'yellow' : 'primary'"
                                        @click="toggleAdmin(user.userID)"
                                    >
                                        {{ user.isAdmin ? '取消管理员' : '设为管理员' }}
                                    </UButton>
                                </div>
                            </div>
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
<!-- 数据库管理面板 -->
<div v-if="isAdmin" class="bg-gray-700 rounded-lg p-4 mb-6">
    <h2 class="text-xl font-semibold text-white mb-4">数据库管理</h2>
    <div class="space-y-4">
        <div class="flex gap-4">
            <UButton
                color="primary"
                icon="i-heroicons-arrow-down-tray"
                @click="downloadBackup"
            >
                下载备份
            </UButton>
            <UButton
                color="warning"
                icon="i-heroicons-arrow-up-tray"
                @click="triggerDatabaseUpload"
            >
                恢复数据库
            </UButton>
        </div>
        <div class="text-yellow-400 text-sm max-h-16 overflow-y-auto bg-gray-800/50 rounded p-2">
            🔔：SQLite一键备份恢复，因兼容问题，如果你在使用云端的PostgreSQL/MySQL数据库，可以尝试，但最好前往云服务端来备份和恢复
        </div>
        <input
            type="file"
            ref="databaseFileInput"
            accept=".zip"
            class="hidden"
            @change="handleDatabaseUpload"
        />
    </div>
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
                            @click="handleLogout"
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
        <input
            type="file"
            ref="fileInput"
            accept="image/*"
            multiple
            class="hidden"
            @change="handleFileUpload"
        />
  
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed, onMounted } from 'vue'
import type { UserToLogin, UserToRegister } from '~/types/models'
import { useUser } from '~/composables/useUser'
import { useUserStore } from '~/store/user'
import { useToast } from '#ui/composables/useToast'

const userStore = useUserStore()
const { login, register, logout } = useUser()
const userToken = ref('')

// 重新生成 Token
// 修改 regenerateToken 函数
const regenerateToken = async () => {
    if (!userStore.isLogin) {
        useToast().add({
            title: '错误',
            description: '请先登录',
            color: 'red'
        });
        return;
    }

    try {
        const response = await fetch('/api/user/token/regenerate', {
            method: 'POST',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json'
            }
        });

        const data = await response.json();
        if (!response.ok) {
            throw new Error(data.msg || 'Token生成请求失败');
        }

        if (data.code === 1 && data.data?.token) {
            userToken.value = data.data.token;
            useToast().add({
                title: '成功',
                description: 'Token 已更新',
                color: 'green'
            });
        } else {
            throw new Error(data.msg || 'Token 生成失败');
        }
    } catch (error: any) {
        console.error('Token生成错误:', error);
        useToast().add({
            title: '错误',
            description: error.message || 'Token 生成失败',
            color: 'red'
        });
    }
};

// 复制 Token
const copyToken = async () => {
    try {
        await navigator.clipboard.writeText(userToken.value)
        useToast().add({
            title: '成功',
            description: 'Token 已复制到剪贴板',
            color: 'green'
        })
    } catch (error) {
        useToast().add({
            title: '错误',
            description: '复制失败',
            color: 'red'
        })
    }
}
// 添加退出登录处理函数
const handleLogout = async () => {
    try {
        // 先调用后端退出接口
        const response = await fetch('/api/user/logout', {
            method: 'POST',
            credentials: 'include' // 确保携带cookie
        })
        
        if (!response.ok) {
            throw new Error('退出请求失败')
        }

        // 清除前端状态
        userStore.$reset()
        
        // 强制刷新页面以确保所有状态被清除
        window.location.reload()
        
        useToast().add({
            title: '成功',
            description: '已退出登录',
            color: 'green'
        })
    } catch (error) {
        useToast().add({
            title: '错误',
            description: '退出登录失败',
            color: 'red'
        })
    }
}
// 状态变量
const isLogin = computed(() => userStore?.isLogin ?? false)
const isAdmin = computed(() => {
    return userStore.user?.is_admin ?? false
})
const authmode = ref(true)
const showLoginModal = ref(false)
const editMode = ref(false)
const fileInput = ref<HTMLInputElement | null>(null) 
    const userForm = reactive({
    username: '',
    oldPassword: '',
    newPassword: ''
})
const editUserInfo = reactive({
    username: false,
    password: false
})
const updateUsername = async () => {
    try {
        if (!userForm.username.trim()) {
            throw new Error('用户名不能为空')
        }
        
        const response = await fetch('/api/user/update', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include',
            body: JSON.stringify({ 
                username: userForm.username,
                type: 'username'  // 明确指定更新类型
            })
        })
        const data = await response.json()
        if (data.code === 1) {
            await userStore.getUser()
            editUserInfo.username = false
            userForm.username = ''
            useToast().add({
                title: '成功',
                description: '用户名已更新',
                color: 'green'
            })
        } else {
            throw new Error(data.msg)
        }
    } catch (error) {
        useToast().add({
            title: '错误',
            description: error.message || '更新失败',
            color: 'red'
        })
    }
}

const updatePassword = async () => {
    try {
        // 检查密码是否为空
        if (!userForm.newPassword || !userForm.oldPassword) {
            throw new Error('密码不能为空')
        }

        // 检查新旧密码是否相同
        if (userForm.newPassword === userForm.oldPassword) {
            throw new Error('新密码不能与当前密码相同')
        }

        const response = await fetch('/api/user/change_password', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include',
            body: JSON.stringify({ 
                password: userForm.newPassword,
                oldPassword: userForm.oldPassword 
            })
        })
        const data = await response.json()
        if (data.code === 1) {
            editUserInfo.password = false
            userForm.oldPassword = ''
            userForm.newPassword = ''
            useToast().add({
                title: '成功',
                description: '密码已更新',
                color: 'green'
            })
        } else {
            throw new Error(data.msg)
        }
    } catch (error) {
        useToast().add({
            title: '错误',
            description: error.message || '更新失败',
            color: 'red'
        })
    }
}

const toggleAdmin = async (userId: number) => {
    try {
        const response = await fetch(`/api/user/admin?id=${userId}`, {
            method: 'PUT',
            credentials: 'include'
        })
        const data = await response.json()
        if (data.code === 1) {
            await userStore.getStatus()
            useToast().add({
                title: '成功',
                description: '权限已更新',
                color: 'green'
            })
        } else {
            throw new Error(data.msg)
        }
    } catch (error) {
        useToast().add({
            title: '错误',
            description: error.message || '更新失败',
            color: 'red'
        })
    }
}
// 配置相关
const configLabels = {
    siteTitle: '站点标题',
    subtitleText: '欢迎语',
    avatarURL: '头像链接',
    username: '用户名',
    description: '个人描述',
    backgrounds: '背景图片',
    cardFooterTitle: '卡片页脚标题',
    cardFooterSubtitle: '卡片页脚副标题',
    pageFooterHTML: '页面底部HTML',
    rssTitle: 'RSS 标题',
    rssDescription: 'RSS 描述',
    rssAuthorName: 'RSS 作者',
    rssFaviconURL: 'RSS 图标链接',
    walineServerURL: 'Waline 评论服务器地址' 
}

const frontendConfig = reactive({
    siteTitle: '',
    subtitleText: '',
    avatarURL: '',
    username: '',
    description: '',
    backgrounds: [] as string[],
    cardFooterTitle: '',
    cardFooterSubtitle: '',
    pageFooterHTML: '',
    rssTitle: '',
    rssDescription: '',
    rssAuthorName: '',
    rssFaviconURL: '',
    walineServerURL: '',
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
    cardFooterTitle: false,
    cardFooterSubtitle: false,
    pageFooterHTML: false,
    rssTitle: false,
    rssDescription: false,
    rssAuthorName: false,
    rssFaviconURL: false,
    walineServerURL: false // 新增
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
        ],
        cardFooterTitle: "Noise·说说·笔记~",
    cardFooterSubtitle: "note.noisework.cn",
    pageFooterHTML: `<div class="text-center text-xs text-gray-400 py-4">来自<a href="https://www.noisework.cn" target="_blank" rel="noopener noreferrer" class="text-orange-400 hover:text-orange-500">Noise</a> 使用<a href="https://github.com/lin-snow/Ech0" target="_blank" rel="noopener noreferrer" class="text-orange-400 hover:text-orange-500">Ech0</a>发布</div>`,
    rssTitle: 'Noise的说说笔记',
    rssDescription: '一个说说笔记~',
    rssAuthorName: 'Noise',
    rssFaviconURL: '/favicon.ico',
    walineServerURL: 'https://app-production-80c1.up.railway.app' 
}
// 添加单个配置项保存方法
const saveConfigItem = async (key: string) => {
    try {
        const response = await fetch('/api/settings', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include',
            body: JSON.stringify({
                allow_registration: true,
                frontendSettings: frontendConfig
            })
        });
        
        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.msg || '请求失败');
        }
        
        const data = await response.json();
        if (data.code === 1) {
            editItem[key] = false;
            
            // 更新前端配置
            await fetchConfig();
            
            // 触发前端配置更新事件
            window.dispatchEvent(new CustomEvent('frontend-config-updated'));
            
            // 如果更新的是 RSS 相关配置，则刷新 RSS
            if (key.startsWith('rss')) {
                await fetch('/api/rss/refresh', {
                    method: 'POST',
                    credentials: 'include'
                });
            }
            
            useToast().add({
                title: '成功',
                description: `${configLabels[key]}已更新`,
                color: 'green'
            });
        } else {
            throw new Error(data.msg || '保存失败');
        }
    } catch (error: any) {
        useToast().add({
            title: '错误',
            description: error.message || '配置更新失败',
            color: 'red'
        });
    }
};
// 添加单个配置项重置方法
const resetConfigItem = (key: string) => {
    frontendConfig[key] = defaultConfig[key]
    editItem[key] = false
}
// 修改 fetchConfig 方法
const fetchConfig = async () => {
    try {
        const response = await fetch('/api/frontend/config', {
            credentials: 'include',
            headers: {
                'Cache-Control': 'no-cache',
                'Pragma': 'no-cache'
            }
        });
        
        const data = await response.json();
        console.log('获取到的配置数据:', data);
        
        if (data && data.frontendSettings) {
            // 使用解构合并，保持配置一致性
            Object.assign(frontendConfig, {
                ...defaultConfig,  // 默认值作为基础
                ...data.frontendSettings  // 服务器配置覆盖默认值
            });
            
            // 特殊处理背景图片数组
            if (!data.frontendSettings.backgrounds?.length) {
                frontendConfig.backgrounds = [...defaultConfig.backgrounds];
            }
            
            // 触发前端配置更新事件
            window.dispatchEvent(new CustomEvent('frontend-config-updated'));
        } else {
            // 使用默认配置
            Object.assign(frontendConfig, defaultConfig);
        }
    } catch (error) {
        console.error('获取配置失败:', error);
        // 发生错误时使用默认配置
        Object.assign(frontendConfig, defaultConfig);
    }
};
const saveConfig = async () => {
    try {
        const response = await fetch('/api/settings', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include',
            body: JSON.stringify({
                frontendSettings: frontendConfig
            })
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.msg || '请求失败');
        }

        const data = await response.json();
        if (data.code === 1) {
            editMode.value = false;
            
            // 更新前端配置
            await fetchConfig();
            
            // 触发前端配置更新事件
            window.dispatchEvent(new CustomEvent('frontend-config-updated'));
            
            // 刷新 RSS
            await fetch('/api/rss/refresh', {
                method: 'POST',
                credentials: 'include'
            });
            
            // 刷新用户状态
            await userStore.getStatus();
            
            useToast().add({
                title: '成功',
                description: '配置已更新',
                color: 'green'
            });
        } else {
            throw new Error(data.msg || '保存失败');
        }
    } catch (error: any) {
        useToast().add({
            title: '错误',
            description: error.message || '配置更新失败',
            color: 'red'
        });
    }
};
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

    const allowedTypes = ['image/jpeg', 'image/png', 'image/webp']
    
    for (const file of Array.from(files)) {
        try {
            if (!allowedTypes.includes(file.type)) {
                throw new Error('仅支持 JPG/PNG/WEBP 格式')
            }

            const formData = new FormData()
            formData.append('image', file)

            const toast = useToast().add({
                title: `正在上传 ${file.name}`,
                description: '请稍候...',
                color: 'blue',
                timeout: 0
            })

            const response = await fetch('/api/images/upload', {
                method: 'POST',
                credentials: 'include',
                body: formData
            })

            const data = await response.json()
            
            if (!response.ok || data.code !== 1) {
                throw new Error(data.msg || '上传失败')
            }

            if (data.code === 1 && data.data) {
                const fullImageUrl = data.data.startsWith('http') 
                    ? data.data 
                    : `${window.location.origin}${data.data}`
                
                if (!frontendConfig.backgrounds.includes(fullImageUrl)) {
                    frontendConfig.backgrounds.push(fullImageUrl)
                    await saveConfig()
                }

                useToast().add({
                    title: '上传成功',
                    description: `${file.name} 已添加`,
                    color: 'green'
                })
            }
        } catch (error: any) {
            useToast().add({
                title: `上传失败: ${file.name}`,
                description: error.message || '文件上传失败',
                color: 'red'
            })
        } finally {
            if (toast) {
                useToast().remove(toast.id)
            }
        }
    }

    new Audio('/sound/success.mp3').play().catch(() => {})
    
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
const isLoading = ref(false) // 新增加载状态

onMounted(async () => {
    try {
        isLoading.value = true;
        
        // 先获取用户状态和配置
        await Promise.all([
            userStore.getStatus(),
            userStore.getUser(),
            fetchConfig()
        ]);

        // 如果用户已登录，再获取 token
        if (userStore.isLogin) {
            const response = await fetch('/api/user/token', {
                method: 'GET',
                credentials: 'include',
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            const data = await response.json();
            if (data.code === 1 && data.data?.token) {
                userToken.value = data.data.token;
            }
        }
    } catch (error) {
        console.error('初始化失败:', error);
        useToast().add({
            title: '初始化失败',
            description: '请刷新页面重试',
            color: 'red',
            timeout: 3000
        });
    } finally {
        isLoading.value = false;
    }
});
const databaseFileInput = ref<HTMLInputElement | null>(null)

const downloadBackup = async () => {
    try {
        const response = await fetch('/api/backup/download', {
            credentials: 'include'
        })
        
        if (!response.ok) {
            throw new Error('下载失败')
        }

        const blob = await response.blob()
        const url = window.URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        a.download = `noise_backup_${new Date().toISOString().slice(0,10)}.zip`
        document.body.appendChild(a)
        a.click()
        window.URL.revokeObjectURL(url)
        document.body.removeChild(a)
    } catch (error) {
        useToast().add({
            title: '错误',
            description: error.message || '备份下载失败',
            color: 'red'
        })
    }
}

const triggerDatabaseUpload = () => {
    databaseFileInput.value?.click()
}
const emit = defineEmits(['restore-success'])
const handleDatabaseUpload = async (event: Event) => {
    const files = (event.target as HTMLInputElement).files
    if (!files || !files[0]) return

    try {
        const formData = new FormData()
        formData.append('database', files[0])

        const response = await fetch('/api/backup/restore', {
            method: 'POST',
            credentials: 'include',
            body: formData
        })

        const data = await response.json()
        if (data.code === 1) {
            useToast().add({
                title: '成功',
                description: '数据库恢复成功',
                color: 'green'
            })
            emit('restore-success')
            // 添加成功后刷新页面
            setTimeout(() => {
                window.location.reload()
            }, 1500)
        } else {
            throw new Error(data.msg)
        }
    } catch (error) {
        useToast().add({
            title: '错误',
            description: error.message || '数据库恢复失败',
            color: 'red'
        })
    }

    if (databaseFileInput.value) {
        databaseFileInput.value.value = ''
    }
}
</script>

<style scoped>
.hidden {
    display: none;
}
</style>