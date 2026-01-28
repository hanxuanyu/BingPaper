// 假期API类型定义
export interface HolidayDay {
  /** 节日名称 */
  name: string;
  /** 日期, ISO 8601 格式 */
  date: string;
  /** 是否为休息日 */
  isOffDay: boolean;
}

export interface Holidays {
  /** 完整年份, 整数。*/
  year: number;
  /** 所用国务院文件网址列表 */
  papers: string[];
  days: HolidayDay[];
}

// 假期数据缓存
const holidayCache = new Map<number, Holidays>();

/**
 * 获取指定年份的假期数据
 */
export async function getHolidaysByYear(year: number): Promise<Holidays | null> {
  // 检查缓存
  if (holidayCache.has(year)) {
    return holidayCache.get(year)!;
  }

  try {
    const response = await fetch(`https://api.coding.icu/cnholiday/${year}.json`);
    
    if (!response.ok) {
      console.warn(`获取${year}年假期数据失败: ${response.status}`);
      return null;
    }

    const data: Holidays = await response.json();
    
    // 缓存数据
    holidayCache.set(year, data);
    
    return data;
  } catch (error) {
    console.error(`获取${year}年假期数据出错:`, error);
    return null;
  }
}

/**
 * 获取指定日期的假期信息
 */
export function getHolidayByDate(holidays: Holidays | null, dateStr: string): HolidayDay | null {
  if (!holidays) return null;
  
  return holidays.days.find(day => day.date === dateStr) || null;
}

/**
 * 清除假期缓存
 */
export function clearHolidayCache() {
  holidayCache.clear();
}
