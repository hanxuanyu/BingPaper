declare module 'lunar-javascript' {
  export class Solar {
    static fromDate(date: Date): Solar
    getLunar(): Lunar
    getFestivals(): string[]
  }

  export class Lunar {
    getYearInChinese(): string
    getMonthInChinese(): string
    getDayInChinese(): string
    getDay(): number
    getJieQi(): string
    getFestivals(): string[]
  }

  export class HolidayUtil {
    // Add methods if needed
  }
}
