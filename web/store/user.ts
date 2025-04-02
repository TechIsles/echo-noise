import type { User, Status, UserToLogin, UserToRegister, Response } from "~/types/models"

export const useUserStore = defineStore("userStore", () => {
    // 状态
    const user = ref<User | null>(null);
    const status = ref<Status | null>(null);
    const isLogin = ref<boolean>(false);
    const toast = useToast()

    // 注册
    const register = async (userToRegister: UserToRegister) => {
        const response = await postRequest<any>("register", userToRegister, {
            credentials: 'include'
        });
        if (!response || response.code !== 1) {
            console.log("注册失败");
            toast.add({
                title: "注册失败",
                description: response?.msg,
                icon: "i-fluent-error-circle-16-filled",
                color: "red",
                timeout: 2000,
            });
            return false;
        }

        return response.code === 1;
    };

    // 登录
    const login = async (userToLogin: UserToLogin) => {
        const response = await postRequest<User>("login", userToLogin, {
            credentials: 'include'
        });
        if (!response || response.code !== 1) {
            console.log("登录失败");
            toast.add({
                title: "登录失败",
                description: response?.msg,
                icon: "i-fluent-error-circle-16-filled",
                color: "red",
                timeout: 2000,
            });
            return false;
        }

        if (response && response.code === 1 && response.data) {
            user.value = response.data;
            isLogin.value = true;
            getStatus();
            return true;
        }

        return false;
    }

    // 获取状态
    const getStatus = async () => {
        const response = await getRequest<Status>("status", {
            credentials: 'include'
        });
        if (!response || response.code !== 1) {
            console.log("获取系统信息失败");
            toast.add({
                title: "获取系统信息失败",
                description: response?.msg,
                icon: "i-fluent-error-circle-16-filled",
                color: "red",
                timeout: 2000,
            });
            return false;
        }

        if (response && response.code === 1 && response.data) {
            status.value = response.data;
            return true;
        }
    }

    // 获取当前登录用户信息
    const getUser = async () => {
        const response = await getRequest<User>("user", {
            credentials: 'include'
        });
        if (!response || response.code !== 1) {
            console.log("获取用户信息失败");
            toast.add({
                title: "获取用户信息失败",
                description: response?.msg,
                icon: "i-fluent-error-circle-16-filled",
                color: "red",
                timeout: 2000,
            });
            return false;
        }

        if (response && response.code === 1 && response.data) {
            user.value = response.data;
            isLogin.value = true;
            return true;
        }
    }

    // 退出登录
    const logout = async () => {
        const response = await postRequest("logout", {}, {
            credentials: 'include'
        });
        
        isLogin.value = false;
        user.value = null;
        status.value = null;

        return true;
    }

    return {
        user,
        status,
        isLogin,
        register,
        login,
        getStatus,
        logout,
        getUser,
    }
})