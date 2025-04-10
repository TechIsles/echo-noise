import { defineStore } from 'pinia'
import type { NotifyConfig } from '~/types/models'

export const useNotifyStore = defineStore('notify', {
  state: () => ({
    config: {
      webhookEnabled: false,
      webhookURL: '',
      telegramEnabled: false,
      telegramToken: '',
      telegramChatID: '',
      weworkEnabled: false,
      weworkKey: '',
      feishuEnabled: false,
      feishuWebhook: '',
      feishuSecret: ''
    } as NotifyConfig
  }),
  actions: {
    async fetchConfig() {
      try {
        const response = await fetch('/api/notify/config', {
          credentials: 'include'
        })
        const data = await response.json()
        if (data.code === 1) {
          this.config = data.data
        }
      } catch (error) {
        console.error('获取推送配置失败:', error)
      }
    },

    async saveConfig(config: NotifyConfig) {
      try {
        const response = await fetch('/api/notify/config', {
          method: 'PUT',
          credentials: 'include',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(config)
        })
        const data = await response.json()
        if (data.code === 1) {
          this.config = config
          return true
        }
        return false
      } catch (error) {
        console.error('保存推送配置失败:', error)
        return false
      }
    },

    async testNotify(type: string) {
      try {
        const response = await fetch('/api/notify/test', {
          method: 'POST',
          credentials: 'include',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ type })
        })
        const data = await response.json()
        return data.code === 1
      } catch (error) {
        console.error('测试推送失败:', error)
        return false
      }
    },

    async sendNotify(params: { content: string; images?: string[]; format?: string }) {
      try {
        // 处理图片的 Markdown 格式
        let fullContent = params.content;
        if (params.images && params.images.length > 0) {
          fullContent += '\n\n'; // 添加空行分隔
          params.images.forEach(img => {
            fullContent += `![](${img})\n`; // 使用 Markdown 图片语法
          });
        }

        const response = await fetch('/api/notify/send', {
          method: 'POST',
          credentials: 'include',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            content: fullContent,
            format: params.format || 'markdown',
            config: this.config
          })
        })
        const data = await response.json()
        return data.code === 1
      } catch (error) {
        console.error('发送推送失败:', error)
        return false
      }
    }
  }
})