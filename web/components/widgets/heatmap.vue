<template>
    <UCard class="mx-auto sm:max-w-2xl hover:shadow-md backdrop-blur-sm bg-black/40 shadow-lg text-white">
      <div class="heatmap-header">
        <h3>内容发布日历</h3>
        <div class="heatmap-legend">
          <span>少</span>
          <div class="legend-colors">
            <div v-for="n in 5" :key="n" :style="{ backgroundColor: getColor(n) }"></div>
          </div>
          <span>多</span>
        </div>
      </div>
      
      <div class="calendar-container" ref="calendarContainer">
        <!-- 月份标题 -->
        <div class="month-labels">
          <div class="month-label" v-for="(month, index) in monthLabels" :key="index">{{ month }}</div>
        </div>
        
        <!-- 星期标签 -->
        <div class="weekday-labels">
          <div class="weekday-label" v-for="day in weekdays" :key="day">周{{ day }}</div>
        </div>
        
        <!-- 热力图网格 -->
        <div class="heatmap-grid">
          <div v-for="(week, i) in calendarData" :key="i" class="heatmap-week">
            <div 
              v-for="(day, j) in week" 
              :key="j" 
              class="heatmap-day"
              :style="{ 
                backgroundColor: day.count ? getColor(day.level) : '#ebedf0',
                opacity: day.count ? 1 : 0.4
              }"
              :title="`${day.date}: ${day.count || 0} 条内容`"
            ></div>
          </div>
        </div>
      </div>
    </UCard>
  </template>
  
  <script setup lang="ts">
  import { ref, onMounted, computed } from 'vue'
  
  const rawData = ref([])
  const calendarData = ref([])
  // 添加日历容器引用
const calendarContainer = ref(null)
  // 生成月份标签
  const monthLabels = computed(() => {
    const today = new Date()
    const oneYearAgo = new Date(today)
    oneYearAgo.setFullYear(new Date().getFullYear() - 1)
    
    const labels = []
    let currentDate = new Date(oneYearAgo)
    
    // 生成12个月的标签
    for (let i = 0; i < 12; i++) {
      const month = currentDate.getMonth()
      labels.push(`${month + 1}月`)
      currentDate.setMonth(month + 1)
    }
    
    return labels
  })
  
  // 中文星期
  const weekdays = ['日', '一', '二', '三', '四', '五', '六']
  
  const getColor = (level: number) => {
    const colors = ['#9be9a8', '#40c463', '#30a14e', '#216e39', '#0e4429']
    return colors[Math.min(level - 1, 4)] || '#ebedf0'
  }
  
  const fetchHeatmapData = async () => {
    try {
      console.log('开始获取热力图数据')
      // 先使用测试数据
      generateTestData()
      
      // 尝试从API获取数据
      const response = await fetch('/api/messages/calendar')
      const data = await response.json()
      
      console.log('热力图API返回数据:', data)
      
      if (data && data.code === 1 && data.data && data.data.length > 0) {
        rawData.value = data.data
        generateCalendarData()
      } else {
        console.log('使用测试数据显示热力图')
        // 已经生成了测试数据，不需要再次生成
      }
    } catch (error) {
      console.error('获取热力图数据失败:', error)
      // 已经生成了测试数据，不需要再次生成
    }
  }
  
  // 生成测试数据
  const generateTestData = () => {
    const today = new Date()
    const oneYearAgo = new Date()
    oneYearAgo.setFullYear(today.getFullYear() - 1)
    
    const testData = []
    let currentDate = new Date(oneYearAgo)
    
    while (currentDate <= today) {
      // 随机生成发布数量
      const count = Math.random() > 0.7 ? Math.floor(Math.random() * 10) + 1 : 0
      
      if (count > 0) {
        testData.push({
          date: formatDate(currentDate),
          count: count
        })
      }
      
      currentDate.setDate(currentDate.getDate() + 1)
    }
    
    rawData.value = testData
    generateCalendarData()
  }
  
  // 格式化日期为 YYYY-MM-DD
  const formatDate = (date: Date) => {
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    return `${year}-${month}-${day}`
  }
  
  // 生成日历数据
  const generateCalendarData = () => {
    // 创建一年的日历数据
    const today = new Date()
    const oneYearAgo = new Date()
    oneYearAgo.setFullYear(today.getFullYear() - 1)
    
    // 创建日期映射
    const dateMap = new Map()
    rawData.value.forEach(item => {
      if (item && item.date && item.count !== undefined) {
        dateMap.set(item.date, item.count)
      }
    })
    
    // 生成日历网格
    const calendar = []
    
    // 确保从周日开始
    let currentDate = new Date(oneYearAgo)
    currentDate.setDate(currentDate.getDate() - currentDate.getDay())
    
    // 生成53周的数据
    for (let week = 0; week < 53; week++) {
      const weekData = []
      
      for (let day = 0; day < 7; day++) {
        const dateStr = formatDate(currentDate)
        const count = dateMap.get(dateStr) || 0
        
        weekData.push({
          date: dateStr,
          count: count,
          level: count ? Math.min(Math.ceil(count / 2), 5) : 0
        })
        
        currentDate.setDate(currentDate.getDate() + 1)
      }
      
      calendar.push(weekData)
    }
    
    calendarData.value = calendar
  }
  
  onMounted(() => {
  console.log('热力图组件已挂载')
  fetchHeatmapData()
  
  // 修改自动滚动逻辑
  nextTick(() => {
    if (calendarContainer.value) {
      const today = new Date()
      const weeksToScroll = Math.floor((today.getTime() - oneYearAgo.getTime()) / (7 * 24 * 60 * 60 * 1000))
      const scrollAmount = weeksToScroll * 15 // 15px = 12px(方块宽度) + 3px(间距)
      
      // 设置滚动位置
      const container = calendarContainer.value.querySelector('.heatmap-grid')
      if (container) {
        container.scrollLeft = scrollAmount
      }
    }
  })
})
  </script>
  
  <style scoped>
  .heatmap-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
  }
/* 移除调试相关样式 */
.heatmap-debug {
  position: relative;
  z-index: 100;
  border: 1px solid rgba(255, 255, 255, 0.2); /* 改为白色半透明边框 */
}
.heatmap-grid {
  display: flex;
  gap: 3px;
  overflow-x: auto;
  padding-bottom: 5px;
  scroll-behavior: smooth; /* 添加平滑滚动效果 */
}
.debug-info {
  background-color: red;
  color: white;
  padding: 4px;
  margin-bottom: 8px;
  text-align: center;
  font-weight: bold;
}
  .heatmap-header h3 {
    margin: 0;
    font-size: 16px;
    color: white;
  }
  
  .heatmap-legend {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 12px;
    color: rgba(255, 255, 255, 0.7);
  }
  
  .legend-colors {
    display: flex;
    gap: 2px;
  }
  
  .legend-colors div {
    width: 12px;
    height: 12px;
    border-radius: 2px;
  }
  
  .calendar-container {
    position: relative;
    padding-top: 20px;
    padding-left: 30px;
  }
  
  .month-labels {
    display: flex;
    position: absolute;
    top: 0;
    left: 30px;
    right: 0;
  }
  
  .month-label {
    flex: 1;
    text-align: center;
    font-size: 12px;
    color: rgba(255, 255, 255, 0.7);
    white-space: nowrap;
  }
  
  .weekday-labels {
    position: absolute;
    left: 0;
    top: 20px;
    display: flex;
    flex-direction: column;
    gap: 4px;
  }
  
  .weekday-label {
    height: 12px;
    line-height: 12px;
    font-size: 10px;
    color: rgba(255, 255, 255, 0.7);
    text-align: right;
    padding-right: 4px;
  }
  
  .heatmap-grid {
    display: flex;
    gap: 3px;
    overflow-x: auto;
    padding-bottom: 5px;
  }
  
  .heatmap-week {
    display: flex;
    flex-direction: column;
    gap: 3px;
  }
  
  .heatmap-day {
    width: 12px;
    height: 12px;
    border-radius: 2px;
    transition: all 0.2s ease;
  }
  
  .heatmap-day:hover {
    transform: scale(1.2);
  }
  
  @media screen and (max-width: 768px) {
  .heatmap-day {
    width: 8px;
    height: 8px;
  }
  
  .month-label {
    font-size: 10px;
  }
  
  .weekday-label {
    font-size: 8px;
    height: 8px;      /* 调整高度以匹配日历格子 */
    line-height: 8px; /* 调整行高以确保文字垂直居中 */
    padding-right: 2px; /* 减小右侧padding */
  }

  .weekday-labels {
    gap: 3px;  /* 调整间距以匹配日历格子 */
  }
}
  </style>