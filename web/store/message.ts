import { defineStore } from "pinia";
import { ref } from "vue";
import type { Message, PageQuery, PageQueryResult } from "~/types/models";

export const useMessageStore = defineStore("messageStore", () => {
  // 状态
  const messages = ref<Message[]>([]);
  const total = ref(0);
  const hasMore = ref(true);
  const page = ref<number>(0);
  const pageSize = ref(5);
  const toast = useToast();
  const loading = ref<boolean>(true); // 添加loading状态

  // 获取笔记列表
//   const getAllMessages = async () => {
//     try {
//       const response = await getRequest<Message[]>("messages");
//       if (!response) {
//         console.error("获取笔记列表失败");
//         toast.add({
//             title: "获取笔记列表失败",
//             description: "请稍后重试",
//             icon: "i-fluent-error-circle-16-filled",
//             color: "red",
//             timeout: 2000,
//             });
//         return null;
//       }

//       messages.value = response.data;
//     } catch (error) {
//       console.error(error);
//     }
//   };

  // 分页获取笔记列表
  const getMessages = async (query: PageQuery) => {
    try {
      if (query.page < page.value) return;
      const response = await postRequest<PageQueryResult>(
        "messages/page",
        query
      );
      if (!response) {
        console.error("获取笔记列表失败");
        toast.add({
            title: "获取笔记列表失败",
            description: "请稍后重试",
            icon: "i-fluent-error-circle-16-filled",
            color: "red",
            timeout: 2000,
            });
        return null;
      }

      // 将返回的笔记数据追加到 messages 数组中
      if (messages.value.length < response.data.total) {
        messages.value.push(...response.data.items);
      }
      total.value = response.data.total;
      page.value += 1; // 更新当前页码

      // 更新分页状态
      hasMore.value = messages.value.length < total.value;
    } catch (error) {
      console.error(error);
    } finally {
      loading.value = false; // 数据加载完成时设置为 false
    }
  };

  const deleteMessage = async (id: number) => {
    try {
      const response = await deleteRequest<any>(`messages/${id}`);
      if (!response || response.code !== 1) {
        console.error("删除笔记失败");
        toast.add({
            title: "删除笔记失败",
            description: response?.msg,
            icon: "i-fluent-error-circle-16-filled",
            color: "red",
            timeout: 2000,
            });
        return null;
      }

      // 从 messages 中删除对应的笔记
      messages.value = messages.value.filter((message) => message.id !== id);
      return response;
    } catch (error) {
      console.error(error);
    }
  }

  return {
    page,
    pageSize,
    messages,
    total,
    hasMore,
    loading,
    getMessages,
    deleteMessage,
  };
});
