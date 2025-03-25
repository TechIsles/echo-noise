<template>
  <div ref="previewElement" class="markdown-preview"></div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';
import Vditor from 'vditor';

// 定义正则表达式
const BILIBILI_REG = /https:\/\/www\.bilibili\.com\/video\/(BV[\w]+)\/?$/;
const YOUTUBE_REG = /https:\/\/(?:www\.)?youtube\.com\/watch\?v=([\w-]+)|https:\/\/youtu\.be\/([\w-]+)/;
const NETEASE_MUSIC_REG = /https:\/\/music\.163\.com(?:\/#)?\/song\?id=(\d+)/;
const QQMUSIC_REG = /https:\/\/y\.qq\.com\/n\/yqq\/song(\w+)\.html/;
const QQVIDEO_REG = /https:\/\/v\.qq\.com\/x\/cover\/\w+\/(\w+)\.html/;
const SPOTIFY_REG = /https:\/\/open\.spotify\.com\/(track|album|playlist)\/([a-zA-Z0-9]+)/;
const YOUKU_REG = /https:\/\/v\.youku\.com\/v_show\/id_([a-zA-Z0-9]+)\.html/;

const previewElement = ref<HTMLDivElement | null>(null);

const props = defineProps({
  content: {
    type: String,
    required: true,
  },
});

const processMediaLinks = (content: string): string => {
  return content
    .replace(BILIBILI_REG, "<div class='video-wrapper'><iframe src='https://www.bilibili.com/blackboard/html5mobileplayer.html?bvid=$1&as_wide=1&high_quality=1&danmaku=0' scrolling='no' border='0' frameborder='no' framespacing='0' allowfullscreen='true' style='position:absolute;height:100%;width:100%'></iframe></div>")
    .replace(YOUTUBE_REG, "<div class='video-wrapper'><iframe src='https://www.youtube.com/embed/$1$2' title='YouTube video player' frameborder='0' allow='accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture' allowfullscreen></iframe></div>")
    .replace(NETEASE_MUSIC_REG, "<div class='music-wrapper'><meting-js auto='https://music.163.com/#/song?id=$1'></meting-js></div>")
    .replace(QQMUSIC_REG, "<meting-js auto='https://y.qq.com/n/yqq/song$1.html'></meting-js>")
    .replace(QQVIDEO_REG, "<div class='video-wrapper'><iframe src='//v.qq.com/iframe/player.html?vid=$1' allowFullScreen='true' frameborder='no'></iframe></div>")
    .replace(SPOTIFY_REG, "<div class='spotify-wrapper'><iframe style='border-radius:12px' src='https://open.spotify.com/embed/$1/$2?utm_source=generator&theme=0' width='100%' frameBorder='0' allowfullscreen='' allow='autoplay; clipboard-write; encrypted-media; fullscreen; picture-in-picture' loading='lazy'></iframe></div>")
    .replace(YOUKU_REG, "<div class='video-wrapper'><iframe src='https://player.youku.com/embed/$1' frameborder=0 'allowfullscreen'></iframe></div>");
};

const renderMarkdown = async (markdown: string) => {
  if (!previewElement.value) return;

  try {
    if (typeof Vditor === 'undefined') {
      console.error('Vditor is not loaded.');
      return;
    }

    const processedContent = processMediaLinks(markdown);

    Vditor.preview(previewElement.value, processedContent, {
      mode: 'light',
      lang: 'zh_CN',
      theme: {
        current: 'dark'
      },
      hljs: {
        style: 'github-dark',
        lineNumber: true,
        enable: true
      },
      after: () => {
        const links = previewElement.value?.getElementsByTagName('a');
        if (links) {
          Array.from(links).forEach(link => {
            link.setAttribute('target', '_blank');
            link.setAttribute('rel', 'noopener noreferrer');
          });
        }
        console.log('Rendering complete.');
      }
    });
  } catch (error) {
    console.error("Error rendering markdown:", error);
    previewElement.value.innerHTML = '';
  }
};

watch(
  () => props.content,
  async (newContent) => {
    await renderMarkdown(newContent);
  },
  { immediate: true }
);

onMounted(() => {
  renderMarkdown(props.content);
// 确保 MetingJS 正确初始化
if (window.APlayer && window.MetingJSElement) {
    console.log('MetingJS is ready');
  } else {
    console.error('MetingJS or APlayer is not loaded properly');
  }
});
</script>

<style>
.markdown-preview {
  font-family: "LXGW WenKai Screen";

}

.markdown-preview h1,
.markdown-preview h2,
.markdown-preview h3,
.markdown-preview h4,
.markdown-preview h5,
.markdown-preview h6 {
  color: rgb(251, 247, 247);
}

.markdown-preview p {
  color: rgb(227, 220, 220);
}

.markdown-preview table thead tr {
  background-color: rgba(223, 226, 229, 0.49) !important;
}

.markdown-preview table tbody tr {
  background-color: rgba(232, 232, 237, 0.39) !important;
}

.video-wrapper {
  position: relative;
  width: 100%;
  padding-bottom: 56.25%; /* 16:9 宽高比 */
  margin: 1em 0;
}

.video-wrapper iframe {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.music-wrapper {
  width: 100%;
  margin: 1em 0;
}

.spotify-wrapper {
  width: 100%;
  margin: 1em 0;
}

.spotify-wrapper iframe {
  width: 100%;
  height: 352px;
}

.markdown-preview :deep(img) {
  max-width: 100%;
  height: auto;
  display: block;
  margin: 1em auto;
}


.markdown-preview :deep(pre) {
  overflow-x: auto;
  background-color: #0d1117;
  border-radius: 6px;
  padding: 16px;
  margin: 1em 0;
  border: 1px solid #30363d;
  max-width: 100%;
  white-space: pre-wrap;
  word-wrap: break-word;
  box-sizing: border-box;
}


.markdown-preview :deep(.hljs) {
  background-color: transparent;
  color: #c9d1d9;
  padding: 0;
}

.markdown-preview :deep(.hljs-keyword) {
  color: #ff7b72;
}

.markdown-preview :deep(.hljs-string) {
  color: #a5d6ff;
}

.markdown-preview :deep(.hljs-comment) {
  color: #8b949e;
  font-style: italic;
}

.markdown-preview :deep(.hljs-function) {
  color: #d2a8ff;
}

.markdown-preview :deep(.hljs-number) {
  color: #79c0ff;
}

.markdown-preview :deep(.hljs-operator) {
  color: #ff7b72;
}

.markdown-preview :deep(.hljs-class) {
  color: #ffa657;
}

.markdown-preview :deep(.hljs-variable) {
  color: #ffa657;
}

.markdown-preview :deep(.hljs-line-numbers) {
  border-right: 1px solid #30363d;
  padding-right: 1em;
  margin-right: 1em;
  color: #6e7681;
  -webkit-user-select: none;
  user-select: none;
}

.markdown-preview :deep(blockquote) {
  border-left: 4px solid #14141484;
  margin: 1em 0;
  padding: 0.5em 1em;
  background-color: rgba(0, 0, 0, 0.05);
}

.markdown-preview :deep(a) {
  color: #0366d6;
  text-decoration: none;
}

.markdown-preview :deep(a:hover) {
  text-decoration: underline;
}

.markdown-preview :deep(table) {
  border-collapse: collapse;
  width: 100%;
  margin: 1em 0;
}

.markdown-preview :deep(th),
.markdown-preview :deep(td) {
  border: 1px solid #ddd;
  padding: 8px;
  text-align: left;
}

.markdown-preview :deep(ul),
.markdown-preview :deep(ol) {
  padding-left: 2em;
}

.markdown-preview :deep(hr) {
  border: none;
  border-top: 1px solid #ddd;
  margin: 1em 0;
}
.music-wrapper {
  width: 100%;
  margin: 1em 0;
  max-width: 800px;
  margin-left: auto;
  margin-right: auto;
}

.aplayer {
  box-shadow: 0 0 10px rgba(0,0,0,0.1);
  border-radius: 4px;
  margin: 1em 0 !important;
}
</style>