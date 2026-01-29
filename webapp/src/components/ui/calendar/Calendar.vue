<template>
  <div class="fixed inset-0 z-40">
    <div 
      ref="calendarPanel"
      class="fixed bg-gradient-to-br from-black/30 via-black/20 to-black/30 backdrop-blur-xl rounded-3xl p-3 sm:p-4 w-[calc(100%-1rem)] sm:w-full max-w-[95vw] sm:max-w-[420px] shadow-2xl border border-white/10 cursor-move select-none"
      :style="{ left: panelPos.x + 'px', top: panelPos.y + 'px' }"
      @mousedown="startDrag"
      @touchstart="startDrag"
      @click.stop
    >
      <!-- 拖动手柄指示器 -->
      <div class="absolute top-2 left-1/2 -translate-x-1/2 w-12 h-1 bg-white/20 rounded-full"></div>
      
      <!-- 头部 -->
      <div class="flex items-center justify-between mb-3 sm:mb-4 mt-2">
        <div class="flex items-center gap-1.5 sm:gap-2 flex-1">
          <button
            @click.stop="previousMonth"
            :disabled="!canGoPrevious"
            class="p-1 sm:p-1.5 hover:bg-white/20 rounded-lg transition-colors disabled:opacity-30 disabled:cursor-not-allowed"
          >
            <svg class="w-3.5 h-3.5 sm:w-4 sm:h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
            </svg>
          </button>
          
          <div class="text-center flex-1">
            <!-- 年月选择器 -->
            <div class="flex items-center justify-center gap-1 sm:gap-1.5 mb-0.5">
              <!-- 年份选择 -->
              <Select v-model="currentYearString" @update:modelValue="onYearChange">
                <SelectTrigger 
                  class="w-[90px] sm:w-[105px] h-6 sm:h-7 bg-white/10 text-white border-white/20 hover:bg-white/20 backdrop-blur-md font-bold text-xs sm:text-sm px-1.5 sm:px-2"
                  @click.stop
                  @mousedown.stop
                >
                  <SelectValue>{{ currentYear }}年</SelectValue>
                </SelectTrigger>
                <SelectContent class="max-h-[300px] bg-gray-900/95 backdrop-blur-xl border-white/20">
                  <SelectItem
                    v-for="year in yearOptions"
                    :key="year"
                    :value="String(year)"
                    class="text-white hover:bg-white/20 focus:bg-white/20 cursor-pointer"
                  >
                    {{ year }}年
                  </SelectItem>
                </SelectContent>
              </Select>
              
              <!-- 月份选择 -->
              <Select v-model="currentMonthString" @update:modelValue="onMonthChange">
                <SelectTrigger 
                  class="w-[65px] sm:w-[75px] h-6 sm:h-7 bg-white/10 text-white border-white/20 hover:bg-white/20 backdrop-blur-md font-bold text-xs sm:text-sm px-1.5 sm:px-2"
                  @click.stop
                  @mousedown.stop
                >
                  <SelectValue>{{ currentMonth + 1 }}月</SelectValue>
                </SelectTrigger>
                <SelectContent class="bg-gray-900/95 backdrop-blur-xl border-white/20">
                  <SelectItem
                    v-for="month in 12"
                    :key="month"
                    :value="String(month - 1)"
                    class="text-white hover:bg-white/20 focus:bg-white/20 cursor-pointer"
                  >
                    {{ month }}月
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>
            
            <div class="text-[10px] sm:text-xs text-white/60 drop-shadow-md font-['Microsoft_YaHei_UI','Microsoft_YaHei',sans-serif] leading-relaxed">
              {{ lunarMonthYear }}
            </div>
          </div>
          
          <button
            @click.stop="nextMonth"
            :disabled="!canGoNext"
            class="p-1 sm:p-1.5 hover:bg-white/20 rounded-lg transition-colors disabled:opacity-30 disabled:cursor-not-allowed"
          >
            <svg class="w-3.5 h-3.5 sm:w-4 sm:h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
            </svg>
          </button>
        </div>
        
        <button
          @click.stop="$emit('close')"
          class="p-1 sm:p-1.5 hover:bg-white/20 rounded-lg transition-colors ml-1.5 sm:ml-2"
        >
          <svg class="w-3.5 h-3.5 sm:w-4 sm:h-4 text-white drop-shadow-lg" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
        </button>
      </div>

      <!-- 星期标题 -->
      <div class="grid grid-cols-7 gap-1 sm:gap-1.5 mb-1.5 sm:mb-2 pointer-events-none">
        <div 
          v-for="(day, idx) in weekDays" 
          :key="day"
          class="text-center text-[11px] sm:text-[13px] font-medium py-1 sm:py-1.5 drop-shadow-md leading-none"
          :class="idx === 0 || idx === 6 ? 'text-red-300/80' : 'text-white/70'"
        >
          {{ day }}
        </div>
      </div>

      <!-- 日期格子 -->
      <div class="grid grid-cols-7 gap-1 sm:gap-1.5">
        <div
          v-for="(day, index) in calendarDays"
          :key="index"
          class="relative aspect-square flex flex-col items-center justify-center rounded-lg transition-opacity pointer-events-none py-0.5 sm:py-1"
          :class="[
            day.isCurrentMonth && !day.isFuture ? 'text-white' : 'text-white/25',
            day.isToday ? 'bg-blue-400/40 ring-2 ring-blue-300/50' : '',
            day.isSelected ? 'bg-white/30 ring-1 ring-white/40' : '',
            day.isFuture ? 'opacity-40' : '',
            day.isWeekend && day.isCurrentMonth ? 'text-red-200/90' : '',
            (day.apiHoliday?.isOffDay || (!day.apiHoliday && day.isWeekend)) ? 'text-red-300' : ''
          ]"
        >
          <!-- 休息/上班标记 (API优先，其次周末) - 使用圆形SVG -->
          <div 
            v-if="day.isCurrentMonth && (day.apiHoliday || day.isWeekend)"
            class="absolute top-0 right-0 w-[14px] h-[14px] sm:w-4 sm:h-4"
          >
            <svg viewBox="0 0 20 20" class="w-full h-full drop-shadow-md">
              <circle 
                cx="10" 
                cy="10" 
                r="9" 
                :fill="day.apiHoliday ? (day.apiHoliday.isOffDay ? '#ef4444' : '#3b82f6') : '#ef4444'"
                opacity="0.65"
              />
              <text 
                x="9.8" 
                y="10.5" 
                text-anchor="middle" 
                dominant-baseline="middle" 
                fill="white" 
                font-size="11" 
                font-weight="bold"
                font-family="'Microsoft YaHei UI','Microsoft YaHei','PingFang SC','Hiragino Sans GB',sans-serif"
              >
                {{ day.apiHoliday ? (day.apiHoliday.isOffDay ? '休' : '班') : '休' }}
              </text>
            </svg>
          </div>
          
          <!-- 公历日期 -->
          <div 
            class="text-[13px] sm:text-[15px] font-medium drop-shadow-md font-['Helvetica','Arial',sans-serif] leading-none mb-0.5 sm:mb-1"
            :class="(day.apiHoliday?.isOffDay || (!day.apiHoliday && day.isWeekend)) ? 'text-red-300 font-bold' : ''"
          >
            {{ day.day }}
          </div>
          
          <!-- 农历/节日/节气 (不显示API节假日名称) -->
          <div 
            class="text-[9px] sm:text-[10px] leading-tight drop-shadow-sm font-['Microsoft_YaHei_UI','Microsoft_YaHei',sans-serif] text-center px-0.5"
            :class="[
              day.festival || day.solarTerm || day.lunarFestival ? 'text-red-300 font-semibold' : 'text-white/60'
            ]"
          >
            {{ day.festival || day.solarTerm || day.lunarFestival || day.lunarDay }}
          </div>
        </div>
      </div>

      <!-- 今日按钮 -->
      <div class="mt-3 sm:mt-4 flex justify-center">
        <button
          @click.stop="goToToday"
          class="px-4 sm:px-5 py-1 sm:py-1.5 bg-white/15 hover:bg-white/30 text-white rounded-lg text-[11px] sm:text-xs font-medium transition-all hover:scale-105 active:scale-95 drop-shadow-lg"
        >
          回到今天
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { Solar } from 'lunar-javascript'
import { 
  Select, 
  SelectContent, 
  SelectItem, 
  SelectTrigger, 
  SelectValue 
} from '@/components/ui/select'
import { getHolidaysByYear, getHolidayByDate, type Holidays, type HolidayDay } from '@/lib/holiday-service'

interface CalendarDay {
  day: number
  isCurrentMonth: boolean
  isToday: boolean
  isSelected: boolean
  isFuture: boolean
  isWeekend: boolean
  isHoliday: boolean
  holidayName: string
  apiHoliday: HolidayDay | null  // API返回的假期信息
  lunarDay: string
  festival: string
  lunarFestival: string
  solarTerm: string
  date: Date
}

const props = defineProps<{
  selectedDate: string // YYYY-MM-DD
}>()

const emit = defineEmits<{
  close: []
}>()

const weekDays = ['日', '一', '二', '三', '四', '五', '六']

// 日历面板位置
const calendarPanel = ref<HTMLElement | null>(null)
const panelPos = ref({ x: 0, y: 0 })
const isDragging = ref(false)
const dragStart = ref({ x: 0, y: 0 })

// 计算图片实际显示区域（与ImageView保持一致）
const getImageDisplayBounds = () => {
  const windowWidth = window.innerWidth
  const windowHeight = window.innerHeight
  
  // 必应图片通常是16:9或类似宽高比
  const imageAspectRatio = 16 / 9
  const windowAspectRatio = windowWidth / windowHeight
  
  let displayWidth: number
  let displayHeight: number
  let offsetX: number
  let offsetY: number
  
  if (windowAspectRatio > imageAspectRatio) {
    // 窗口更宽，图片上下占满，左右留黑边
    displayHeight = windowHeight
    displayWidth = displayHeight * imageAspectRatio
    offsetX = (windowWidth - displayWidth) / 2
    offsetY = 0
  } else {
    // 窗口更高，图片左右占满，上下留黑边
    displayWidth = windowWidth
    displayHeight = displayWidth / imageAspectRatio
    offsetX = 0
    offsetY = (windowHeight - displayHeight) / 2
  }
  
  return {
    left: offsetX,
    top: offsetY,
    right: offsetX + displayWidth,
    bottom: offsetY + displayHeight,
    width: displayWidth,
    height: displayHeight
  }
}

// 初始化面板位置（移动端居中，桌面端右上角，限制在图片显示区域内）
const initPanelPosition = () => {
  if (typeof window !== 'undefined') {
    const bounds = getImageDisplayBounds()
    const isMobile = window.innerWidth < 640 // sm breakpoint
    
    if (isMobile) {
      // 移动端：在图片区域内居中显示
      const panelWidth = Math.min(bounds.width - 16, window.innerWidth - 16)
      const panelHeight = 580 // 估计高度
      panelPos.value = {
        x: Math.max(bounds.left, bounds.left + (bounds.width - panelWidth) / 2),
        y: Math.max(bounds.top + 8, bounds.top + (bounds.height - panelHeight) / 2)
      }
    } else {
      // 桌面端：在图片区域右上角
      const panelWidth = Math.min(420, bounds.width * 0.9)
      const panelHeight = 600
      panelPos.value = {
        x: bounds.right - panelWidth - 40,
        y: Math.max(bounds.top + 80, bounds.top + (bounds.height - panelHeight) / 2)
      }
    }
  }
}

const currentYear = ref(new Date().getFullYear())
const currentMonth = ref(new Date().getMonth())
const isChangingMonth = ref(false)

// 假期数据
const holidaysData = ref<Map<number, Holidays | null>>(new Map())
const loadingHolidays = ref(false)

// 字符串版本的年月（用于Select组件）
const currentYearString = computed({
  get: () => String(currentYear.value),
  set: (val: string) => {
    currentYear.value = Number(val)
  }
})

const currentMonthString = computed({
  get: () => String(currentMonth.value),
  set: (val: string) => {
    currentMonth.value = Number(val)
  }
})

// 年份改变处理
const onYearChange = (value: any) => {
  if (value !== null && value !== undefined) {
    currentYear.value = Number(value)
  }
}

// 月份改变处理
const onMonthChange = (value: any) => {
  if (value !== null && value !== undefined) {
    currentMonth.value = Number(value)
  }
}

// 生成年份选项（从2009年到当前年份+10年）
const yearOptions = computed(() => {
  const currentYearValue = new Date().getFullYear()
  const years: number[] = []
  for (let year = currentYearValue - 30; year <= currentYearValue + 10; year++) {
    years.push(year)
  }
  return years
})

// 计算是否可以切换月份（不限制）
const canGoPrevious = computed(() => {
  return !isChangingMonth.value
})

const canGoNext = computed(() => {
  return !isChangingMonth.value
})

// 初始化为选中的日期
watch(() => props.selectedDate, (newDate) => {
  if (newDate) {
    const date = new Date(newDate)
    currentYear.value = date.getFullYear()
    currentMonth.value = date.getMonth()
  }
}, { immediate: true })

// 初始化位置
initPanelPosition()

// 加载假期数据
const loadHolidaysForYear = async (year: number) => {
  if (holidaysData.value.has(year)) {
    return
  }
  
  loadingHolidays.value = true
  try {
    const data = await getHolidaysByYear(year)
    holidaysData.value.set(year, data)
  } catch (error) {
    console.error(`加载${year}年假期数据失败:`, error)
    holidaysData.value.set(year, null)
  } finally {
    loadingHolidays.value = false
  }
}

// 组件挂载时加载当前年份的假期数据
onMounted(() => {
  const currentYearValue = currentYear.value
  loadHolidaysForYear(currentYearValue)
  // 预加载前后一年的数据
  loadHolidaysForYear(currentYearValue - 1)
  loadHolidaysForYear(currentYearValue + 1)
})

// 监听年份变化，加载对应的假期数据
watch(currentYear, (newYear) => {
  loadHolidaysForYear(newYear)
  // 预加载前后一年
  loadHolidaysForYear(newYear - 1)
  loadHolidaysForYear(newYear + 1)
})

// 开始拖动
const startDrag = (e: MouseEvent | TouchEvent) => {
  const target = e.target as HTMLElement
  // 如果点击的是按钮或其子元素，不触发拖拽
  if (target.closest('button') || target.closest('[class*="grid"]')) {
    return
  }
  
  e.preventDefault()
  isDragging.value = true
  
  const clientX = e instanceof MouseEvent ? e.clientX : e.touches[0].clientX
  const clientY = e instanceof MouseEvent ? e.clientY : e.touches[0].clientY
  
  dragStart.value = {
    x: clientX - panelPos.value.x,
    y: clientY - panelPos.value.y
  }
  
  document.addEventListener('mousemove', onDrag)
  document.addEventListener('mouseup', stopDrag)
  document.addEventListener('touchmove', onDrag, { passive: false })
  document.addEventListener('touchend', stopDrag)
}

// 拖动中
const onDrag = (e: MouseEvent | TouchEvent) => {
  if (!isDragging.value) return
  
  if (e instanceof TouchEvent) {
    e.preventDefault()
  }
  
  const clientX = e instanceof MouseEvent ? e.clientX : e.touches[0].clientX
  const clientY = e instanceof MouseEvent ? e.clientY : e.touches[0].clientY
  
  const newX = clientX - dragStart.value.x
  const newY = clientY - dragStart.value.y
  
  // 限制在图片实际显示区域内
  if (calendarPanel.value) {
    const rect = calendarPanel.value.getBoundingClientRect()
    const bounds = getImageDisplayBounds()
    
    const minX = bounds.left
    const maxX = bounds.right - rect.width
    const minY = bounds.top
    const maxY = bounds.bottom - rect.height
    
    panelPos.value = {
      x: Math.max(minX, Math.min(newX, maxX)),
      y: Math.max(minY, Math.min(newY, maxY))
    }
  }
}

// 停止拖动
const stopDrag = () => {
  isDragging.value = false
  document.removeEventListener('mousemove', onDrag)
  document.removeEventListener('mouseup', stopDrag)
  document.removeEventListener('touchmove', onDrag)
  document.removeEventListener('touchend', stopDrag)
}

// 农历月份年份
const lunarMonthYear = computed(() => {
  const solar = Solar.fromDate(new Date(currentYear.value, currentMonth.value, 15))
  const lunar = solar.getLunar()
  return `${lunar.getYearInChinese()}年${lunar.getMonthInChinese()}月`
})

// 获取日历天数
const calendarDays = computed<CalendarDay[]>(() => {
  const year = currentYear.value
  const month = currentMonth.value
  
  // 当月第一天
  const firstDay = new Date(year, month, 1)
  const firstDayWeek = firstDay.getDay()
  
  // 当月最后一天
  const lastDay = new Date(year, month + 1, 0)
  const lastDate = lastDay.getDate()
  
  // 上月最后几天
  const prevLastDay = new Date(year, month, 0)
  const prevLastDate = prevLastDay.getDate()
  
  const days: CalendarDay[] = []
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  
  // 填充上月日期
  for (let i = firstDayWeek - 1; i >= 0; i--) {
    const day = prevLastDate - i
    const date = new Date(year, month - 1, day)
    days.push(createDayObject(date, false))
  }
  
  // 填充当月日期
  for (let day = 1; day <= lastDate; day++) {
    const date = new Date(year, month, day)
    days.push(createDayObject(date, true))
  }
  
  // 填充下月日期
  const remainingDays = 42 - days.length // 6行7列
  for (let day = 1; day <= remainingDays; day++) {
    const date = new Date(year, month + 1, day)
    days.push(createDayObject(date, false))
  }
  
  return days
})

// 创建日期对象
const createDayObject = (date: Date, isCurrentMonth: boolean): CalendarDay => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  
  const selectedDate = new Date(props.selectedDate)
  selectedDate.setHours(0, 0, 0, 0)
  
  // 转换为农历
  const solar = Solar.fromDate(date)
  const lunar = solar.getLunar()
  
  // 获取节日
  const festivals = solar.getFestivals()
  const festival = festivals.length > 0 ? festivals[0] : ''
  
  // 获取农历节日
  const lunarFestivals = lunar.getFestivals()
  const lunarFestival = lunarFestivals.length > 0 ? lunarFestivals[0] : ''
  
  // 获取节气
  const solarTerm = lunar.getJieQi()
  
  // 获取API假期数据 - 使用本地时间避免时区偏移
  const dateStr = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
  const yearHolidays = holidaysData.value.get(date.getFullYear())
  const apiHoliday = getHolidayByDate(yearHolidays || null, dateStr)
  
  // 检查是否为假期（使用lunar-javascript的节日信息）
  let isHoliday = false
  let holidayName = ''
  
  try {
    if (festival || lunarFestival) {
      // 常见法定节假日
      const legalHolidays = ['元旦', '春节', '清明', '劳动节', '端午', '中秋', '国庆']
      const holidayNames = [festival, lunarFestival].filter(Boolean)
      
      for (const name of holidayNames) {
        if (legalHolidays.some(legal => name.includes(legal))) {
          isHoliday = true
          holidayName = name
          break
        }
      }
    }
  } catch (e) {
    console.debug('假期信息获取失败:', e)
  }
  
  // 判断是否为周末（周六或周日）
  const isWeekend = date.getDay() === 0 || date.getDay() === 6
  
  // 农历日期显示
  let lunarDay = lunar.getDayInChinese()
  if (lunar.getDay() === 1) {
    lunarDay = lunar.getMonthInChinese() + '月'
  }
  
  return {
    day: date.getDate(),
    isCurrentMonth,
    isToday: date.getTime() === today.getTime(),
    isSelected: date.getTime() === selectedDate.getTime(),
    isFuture: date > today,
    isWeekend,
    isHoliday,
    holidayName,
    apiHoliday,
    lunarDay,
    festival,
    lunarFestival,
    solarTerm,
    date
  }
}

// 上一月
const previousMonth = () => {
  if (!canGoPrevious.value) return
  
  if (currentMonth.value === 0) {
    currentMonth.value = 11
    currentYear.value--
  } else {
    currentMonth.value--
  }
}

// 下一月
const nextMonth = () => {
  if (!canGoNext.value) return
  
  if (currentMonth.value === 11) {
    currentMonth.value = 0
    currentYear.value++
  } else {
    currentMonth.value++
  }
}

// 回到今天
const goToToday = () => {
  const today = new Date()
  currentYear.value = today.getFullYear()
  currentMonth.value = today.getMonth()
}

// 不再支持点击日期选择
// 日历仅作为台历展示功能

// 清理
import { onUnmounted } from 'vue'
onUnmounted(() => {
  document.removeEventListener('mousemove', onDrag)
  document.removeEventListener('mouseup', stopDrag)
  document.removeEventListener('touchmove', onDrag)
  document.removeEventListener('touchend', stopDrag)
})
</script>
