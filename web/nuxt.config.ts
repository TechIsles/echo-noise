// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  app: {
    head: {
      link: [
        { rel: 'stylesheet', href: 'https://www.noisework.cn/css/APlayer.min.css' },
      ],
      script: [
        { src: 'https://cdn.jsdelivr.net/npm/jquery@3.2.1/dist/jquery.min.js', body: true },
        { src: 'https://cdn.jsdelivr.net/npm/aplayer@1.10.1/dist/APlayer.min.js', body: true },
        { src: 'https://cdn.jsdelivr.net/npm/meting@2.0.1/dist/Meting.min.js', body: true },
      ],
      title: "Noise·说说·笔记~",
      meta: [
        { name: "viewport", content: "width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no" }
      ]
    }
  },
  compatibilityDate: '2024-11-01',
  devtools: { enabled: true },
  modules: [
    '@nuxt/ui',
    '@nuxt/fonts',
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt',
  ],
  css: [
    '@/assets/fonts/result.css',
  ],
  colorMode: {
    preference: 'light'
  },
  runtimeConfig: {
    public: {
      baseApi: process.env.BASE_API,
    }
  }
})