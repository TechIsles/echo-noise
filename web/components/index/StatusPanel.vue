<template>
    <div class="h-3/5 max-w-sm">
        <h1 class="text-5xl italic font-bold text-center text-gray-300 mb-14">后台管理</h1>
        
        <!-- 登录后显示的内容 -->
        <div v-if="showStatus">
            <!-- 当前系统状态 -->
            <!-- 当前系统管理员 -->
            <div class="flex justify-start items-center">
            <span class="text-gray-500 text-md">当前系统管理员:</span>
            <span class="font-bold text-md text-gray-600 ml-2">{{ userStore?.status?.username }}</span>
           </div>
             <!-- 当前登录用户 -->
       <div>
       <span class="text-gray-500 text-md">当前登录用户:</span>
      <span class="font-bold text-md text-gray-600 ml-2">{{ isLogin ? (userStore?.user?.username || userStore?.status?.username) : "当前未登录"}}</span>
      </div>
            <!-- 当前共有笔记 -->
            <div class="flex justify-start items-center">
                <span class="text-gray-500 text-md">当前笔记共有:</span>
                <span class="font-bold text-md text-gray-600 ml-2">{{ userStore?.status?.total_messages }}</span>
                <span class="text-gray-500 text-md ml-2">条</span>
            </div>
            <!-- 操作按钮 -->
             <div class="flex justify-between items-center mt-6">
                <div>
                    <!-- 返回首页 -->
                    <UButton size="sm" icon="i-fluent-arrow-left-16-filled" @click="$router.push('/')" color="gray" variant="ghost"
                        class="text-gray-500 mt-6" />
                </div>
                <div>
                    <div v-if="isLogin">
                        <UButton size="sm" icon="i-mdi-logout" @click="logout" color="gray" variant="solid"
                            class="text-gray-500 mt-6" />
                    </div>
                    <div v-else>
                        <UButton size="sm" icon="i-mdi-account-key-outline" @click="showStatusChange" color="gray"
                            variant="solid" class="text-gray-500 mt-6" />
                    </div>
                </div>
             </div>
            
        </div>
        <div v-else>
            <!-- 登录或注册 -->
            <div class="w-full h-full">
                <div v-if="authmode" class="p-3">
                    <div class="flex items-center justify-between mb-2">
                        <span class="text-gray-500 text-md ml-1">登录</span>
                        <UButton size="sm" icon="i-mdi-account-plus-outline" @click="authmode = false" color="gray"
                            variant="solid" class="text-gray-400" />
                    </div>

                    <UForm :state="authForm" class="flex flex-col gap-2 mt-1">
                        <UInput v-model="authForm.username" placeholder="用户名" />
                        <UInput v-model="authForm.password" type="password" placeholder="密码" />
                        <UButton size="md" @click="login(authForm)" color="gray" variant="solid"
                            class="text-gray-500 mt-6 w-1/3 mx-auto">
                            <span class="mx-auto">
                                登录
                            </span>
                        </UButton>
                    </UForm>
                </div>
                <div v-else class="p-3">
                    <div class="flex items-center justify-between mb-2">
                        <span class="text-gray-500 text-md ml-1">注册</span>
                        <UButton size="sm" icon="i-mdi-account-key-outline" @click="authmode = true" color="gray"
                            variant="solid" class="text-gray-400" />
                    </div>
                    <UForm :state="authForm" class="flex flex-col gap-2 mt-1">
                        <UInput v-model="authForm.username" placeholder="用户名" />
                        <UInput v-model="authForm.password" type="password" placeholder="密码" />
                        <UButton size="md" @click="register(authForm)" color="gray" variant="solid"
                            class="text-gray-500 mt-6 w-1/3 mx-auto">
                            <span class="mx-auto">
                                注册
                            </span>
                        </UButton>
                    </UForm>
                </div>
            </div>
        </div>

    </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed } from 'vue'
import type { UserToLogin, UserToRegister } from '~/types/models';
import { useUser } from '~/composables/useUser';
import { useUserStore } from '~/store/user';

const isLogin = computed(() => userStore?.isLogin ?? false);
const authmode = ref<boolean>(true); // authmode : true 登录， false 注册
const showStatus = ref<boolean>(true); // 是否显示状态

const userStore = useUserStore()
const { login, register, logout } = useUser()

// 添加对登录状态的监听
watch(() => userStore.isLogin, (newVal) => {
    if (!newVal) {
        // 当登出时，重新获取状态并清除用户信息
        userStore.getStatus()
        userStore.getUser()
        userStore.$reset() // 重置 store 状态
    }
})

const showStatusChange = () => {
    showStatus.value = !showStatus.value
}

const authForm = reactive<UserToLogin | UserToRegister>({
    username: '',
    password: ''
})

onMounted(() => {
    showStatus.value = true
    userStore.getStatus()
    userStore.getUser()
})
</script>