export default defineNuxtPlugin(() => {
  return {
    provide: {
      initAplayer: () => {
        // APlayer 的初始化逻辑将在客户端执行
        return window.APlayer
      }
    }
  }
})