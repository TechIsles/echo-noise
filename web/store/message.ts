import { defineStore } from "pinia";
import { ref } from "vue";
import type { Message, PageQuery, PageQueryResult } from "~/types/models";

export const useMessageStore = defineStore("messageStore", () => {
  // 状态
  const messages = ref<Message[]>([]);
  const total = ref(0);
  const hasMore = ref(true);
  const page = ref<number>(1);
  const pageSize = ref(10);
  const toast = useToast();
  const loading = ref<boolean>(false);

  // 重置状态
  const reset = () => {
    messages.value = [];
    total.value = 0;
    hasMore.value = true;
    page.value = 1;
    loading.value = false;
  };

  // 分页获取笔记列表
  const getMessages = async (query: PageQuery) => {
    if (loading.value) return;
    loading.value = true;

    try {
      const response = await postRequest<PageQueryResult>("messages/page", query);
      
      if (!response) {
        toast.add({
          title: "获取笔记列表失败",
          description: "请稍后重试",
          icon: "i-fluent-error-circle-16-filled",
          color: "red",
          timeout: 2000,
        });
        return null;
      }

      // 根据页码决定是替换还是追加数据
      if (query.page === 1) {
        messages.value = response.data.items;
      } else {
        messages.value = [...messages.value, ...response.data.items];
      }

      total.value = response.data.total;
      page.value = query.page;
      hasMore.value = messages.value.length < total.value;

    } catch (error) {
      console.error("获取笔记列表失败:", error);
      toast.add({
        title: "获取笔记列表失败",
        description: "请稍后重试",
        icon: "i-fluent-error-circle-16-filled",
        color: "red",
        timeout: 2000,
      });
    } finally {
      loading.value = false;
    }
  };

  // 删除笔记
  const deleteMessage = async (id: number) => {
    try {
      const response = await deleteRequest<any>(`messages/${id}`);
      
      if (!response || response.code !== 1) {
        toast.add({
          title: "删除笔记失败",
          description: response?.msg,
          icon: "i-fluent-error-circle-16-filled",
          color: "red",
          timeout: 2000,
        });
        return null;
      }

      // 从本地状态中移除被删除的笔记
      messages.value = messages.value.filter((message) => message.id !== id);
      total.value -= 1;
      
      return response;
    } catch (error) {
      console.error("删除笔记失败:", error);
      throw error;
    }
  };

  // 更新单条笔记
  const updateMessage = (updatedMessage: Message) => {
    const index = messages.value.findIndex(msg => msg.id === updatedMessage.id);
    if (index !== -1) {
      messages.value[index] = updatedMessage;
    }
  };

  return {
    messages,
    total,
    hasMore,
    page,
    pageSize,
    loading,
    reset,
    getMessages,
    deleteMessage,
    updateMessage,
  };
});