<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h3 class="text-lg font-semibold">系统配置</h3>
      <div class="flex gap-2">
        <Button
          variant="outline"
          @click="editMode = editMode === 'json' ? 'form' : 'json'"
        >
          切换到{{ editMode === 'json' ? '表单' : 'JSON' }}编辑
        </Button>
        <Button @click="handleSaveConfig" :disabled="saveLoading">
          {{ saveLoading ? '保存中...' : '保存配置' }}
        </Button>
      </div>
    </div>

    <div v-if="loading" class="text-center py-8">
      <p class="text-gray-500">加载配置中...</p>
    </div>

    <div v-else-if="loadError" class="text-red-600 bg-red-50 p-4 rounded-md">
      {{ loadError }}
    </div>

    <div v-else>
      <!-- JSON 编辑模式 -->
      <Card v-if="editMode === 'json'">
        <CardHeader>
          <div class="flex justify-between items-start">
            <div>
              <CardTitle>JSON 配置编辑器</CardTitle>
              <CardDescription>
                直接编辑配置 JSON，请确保格式正确
              </CardDescription>
            </div>
            <Button
              variant="outline"
              size="sm"
              @click="formatJson"
              :disabled="!configJson.trim()"
            >
              格式化 JSON
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <Textarea
            v-model="configJson"
            class="font-mono text-sm min-h-[500px]"
            :class="{ 'border-red-500': jsonError }"
            placeholder="配置 JSON"
          />
          <div v-if="jsonError" class="mt-2 text-sm text-red-600 bg-red-50 p-2 rounded">
            ❌ {{ jsonError }}
          </div>
          <div v-else-if="isValidJson" class="mt-2 text-sm text-green-600">
            ✓ JSON 格式正确
          </div>
        </CardContent>
      </Card>

      <!-- 表单编辑模式 -->
      <div v-else class="space-y-4">
        <!-- 服务器配置 -->
        <Card>
          <CardHeader>
            <CardTitle>服务器配置</CardTitle>
          </CardHeader>
          <CardContent class="space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div class="space-y-2">
                <Label>端口</Label>
                <Input v-model.number="config.Server.Port" type="number" />
              </div>
              <div class="space-y-2">
                <Label>基础 URL</Label>
                <Input v-model="config.Server.BaseURL" placeholder="http://localhost:8080" />
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- API 配置 -->
        <Card>
          <CardHeader>
            <CardTitle>API 配置</CardTitle>
          </CardHeader>
          <CardContent class="space-y-4">
            <div class="space-y-2">
              <Label>API 模式</Label>
              <Select v-model="config.API.Mode">
                <SelectTrigger>
                  <SelectValue placeholder="选择 API 模式" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="local">本地 (local)</SelectItem>
                  <SelectItem value="redirect">重定向 (redirect)</SelectItem>
                </SelectContent>
              </Select>
              <p class="text-xs text-gray-500">
                local: 直接返回图片流; redirect: 重定向到存储位置
              </p>
            </div>
            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <div class="space-y-0.5">
                  <Label for="api-fallback">启用地区不存在时兜底</Label>
                  <p class="text-xs text-gray-500">
                    如果请求的地区无数据，自动回退到默认地区
                  </p>
                </div>
                <Switch
                  id="api-fallback"
                  v-model="config.API.EnableMktFallback"
                />
              </div>

              <div class="flex items-center justify-between">
                <div class="space-y-0.5">
                  <Label for="api-on-demand">启用按需实时抓取</Label>
                  <p class="text-xs text-gray-500">
                    如果请求的地区无数据，尝试实时从 Bing 抓取
                  </p>
                </div>
                <Switch
                  id="api-on-demand"
                  v-model="config.API.EnableOnDemandFetch"
                />
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- 定时任务配置 -->
        <Card>
          <CardHeader>
            <CardTitle>定时任务配置</CardTitle>
          </CardHeader>
          <CardContent class="space-y-4">
            <div class="flex items-center gap-2">
              <Label for="cron-enabled">启用定时任务</Label>
              <Switch
                id="cron-enabled"
                v-model="config.Cron.Enabled"
              />
            </div>
            <div class="space-y-2">
              <Label>定时表达式 (Cron)</Label>
              <Input v-model="config.Cron.DailySpec" placeholder="0 9 * * *" />
              <p class="text-xs text-gray-500">
                例如: "0 9 * * *" 表示每天 9:00 执行
              </p>
            </div>
          </CardContent>
        </Card>

        <!-- 数据库配置 -->
        <Card>
          <CardHeader>
            <CardTitle>数据库配置</CardTitle>
          </CardHeader>
          <CardContent class="space-y-4">
            <div class="space-y-2">
              <Label>数据库类型</Label>
              <Select v-model="config.DB.Type">
                <SelectTrigger>
                  <SelectValue placeholder="选择数据库类型" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="sqlite">SQLite</SelectItem>
                  <SelectItem value="mysql">MySQL</SelectItem>
                  <SelectItem value="postgres">PostgreSQL</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div class="space-y-2">
              <Label>DSN (数据源名称)</Label>
              <Input 
                v-model="config.DB.DSN" 
                placeholder="数据库连接字符串"
                :class="{ 'border-red-500': dsnError }"
                @blur="validateDSN"
              />
              <p v-if="dsnExamples" class="text-xs text-gray-500">
                💡 示例: {{ dsnExamples }}
              </p>
              <p v-if="dsnError" class="text-xs text-red-600">
                ❌ {{ dsnError }}
              </p>
            </div>
          </CardContent>
        </Card>

        <!-- 存储配置 -->
        <Card>
          <CardHeader>
            <CardTitle>存储配置</CardTitle>
          </CardHeader>
          <CardContent class="space-y-4">
            <div class="space-y-2">
              <Label>存储类型</Label>
              <Select v-model="config.Storage.Type">
                <SelectTrigger>
                  <SelectValue placeholder="选择存储类型" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="local">本地存储</SelectItem>
                  <SelectItem value="s3">S3 存储</SelectItem>
                  <SelectItem value="webdav">WebDAV</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <!-- 本地存储配置 -->
            <div v-if="config.Storage.Type === 'local'" class="space-y-2">
              <Label>本地存储路径</Label>
              <Input v-model="config.Storage.Local.Root" placeholder="./data/images" />
            </div>

            <!-- S3 存储配置 -->
            <div v-if="config.Storage.Type === 's3'" class="space-y-4">
              <div class="grid grid-cols-2 gap-4">
                <div class="space-y-2">
                  <Label>Endpoint</Label>
                  <Input v-model="config.Storage.S3.Endpoint" />
                </div>
                <div class="space-y-2">
                  <Label>Region</Label>
                  <Input v-model="config.Storage.S3.Region" />
                </div>
              </div>
              <div class="space-y-2">
                <Label>Bucket</Label>
                <Input v-model="config.Storage.S3.Bucket" />
              </div>
              <div class="grid grid-cols-2 gap-4">
                <div class="space-y-2">
                  <Label>Access Key</Label>
                  <Input v-model="config.Storage.S3.AccessKey" type="password" />
                </div>
                <div class="space-y-2">
                  <Label>Secret Key</Label>
                  <Input v-model="config.Storage.S3.SecretKey" type="password" />
                </div>
              </div>
              <div class="space-y-2">
                <Label>公开 URL 前缀</Label>
                <Input v-model="config.Storage.S3.PublicURLPrefix" />
              </div>
              <div class="flex items-center gap-2">
                <Label for="s3-force-path">强制路径样式</Label>
                <Switch
                  id="s3-force-path"
                  v-model="config.Storage.S3.ForcePathStyle"
                />
              </div>
            </div>

            <!-- WebDAV 配置 -->
            <div v-if="config.Storage.Type === 'webdav'" class="space-y-4">
              <div class="space-y-2">
                <Label>WebDAV URL</Label>
                <Input v-model="config.Storage.WebDAV.URL" />
              </div>
              <div class="grid grid-cols-2 gap-4">
                <div class="space-y-2">
                  <Label>用户名</Label>
                  <Input v-model="config.Storage.WebDAV.Username" />
                </div>
                <div class="space-y-2">
                  <Label>密码</Label>
                  <Input v-model="config.Storage.WebDAV.Password" type="password" />
                </div>
              </div>
              <div class="space-y-2">
                <Label>公开 URL 前缀</Label>
                <Input v-model="config.Storage.WebDAV.PublicURLPrefix" />
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- 保留策略配置 -->
        <Card>
          <CardHeader>
            <CardTitle>图片保留策略</CardTitle>
          </CardHeader>
          <CardContent class="space-y-4">
            <div class="space-y-2">
              <Label>保留天数</Label>
              <Input v-model.number="config.Retention.Days" type="number" min="1" />
              <p class="text-xs text-gray-500">
                超过指定天数的图片将被自动清理
              </p>
            </div>
          </CardContent>
        </Card>

        <!-- Token 配置 -->
        <Card>
          <CardHeader>
            <CardTitle>Token 配置</CardTitle>
          </CardHeader>
          <CardContent class="space-y-4">
            <div class="space-y-2">
              <Label>默认过期时间 (TTL)</Label>
              <Input v-model="config.Token.DefaultTTL" placeholder="168h" />
              <p class="text-xs text-gray-500">
                例如: 168h (7天), 720h (30天)
              </p>
            </div>
          </CardContent>
        </Card>

        <!-- 日志配置 -->
        <Card>
          <CardHeader>
            <CardTitle>日志配置</CardTitle>
          </CardHeader>
          <CardContent class="space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div class="space-y-2">
                <Label>日志级别</Label>
                <Select v-model="config.Log.Level">
                  <SelectTrigger>
                    <SelectValue placeholder="选择日志级别" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="debug">Debug</SelectItem>
                    <SelectItem value="info">Info</SelectItem>
                    <SelectItem value="warn">Warn</SelectItem>
                    <SelectItem value="error">Error</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div class="space-y-2">
                <Label>日志文件</Label>
                <Input v-model="config.Log.Filename" />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div class="space-y-2">
                <Label>数据库日志文件</Label>
                <Input v-model="config.Log.DBFilename" />
              </div>
              <div class="space-y-2">
                <Label>数据库日志级别</Label>
                <Select v-model="config.Log.DBLogLevel">
                  <SelectTrigger>
                    <SelectValue placeholder="选择数据库日志级别" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="debug">Debug</SelectItem>
                    <SelectItem value="info">Info</SelectItem>
                    <SelectItem value="warn">Warn</SelectItem>
                    <SelectItem value="error">Error</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
            <div class="grid grid-cols-3 gap-4">
              <div class="flex items-center gap-2">
                <Label for="log-console">输出到控制台</Label>
                <Switch
                  id="log-console"
                  v-model="config.Log.LogConsole"
                />
              </div>
              <div class="flex items-center gap-2">
                <Label for="log-show-db">显示数据库日志</Label>
                <Switch
                  id="log-show-db"
                  v-model="config.Log.ShowDBLog"
                />
              </div>
              <div class="flex items-center gap-2">
                <Label for="log-compress">压缩旧日志</Label>
                <Switch
                  id="log-compress"
                  v-model="config.Log.Compress"
                />
              </div>
            </div>
            <div class="grid grid-cols-3 gap-4">
              <div class="space-y-2">
                <Label>单文件大小 (MB)</Label>
                <Input v-model.number="config.Log.MaxSize" type="number" />
              </div>
              <div class="space-y-2">
                <Label>最大文件数</Label>
                <Input v-model.number="config.Log.MaxBackups" type="number" />
              </div>
              <div class="space-y-2">
                <Label>保留天数</Label>
                <Input v-model.number="config.Log.MaxAge" type="number" />
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- 抓取配置 -->
        <Card>
          <CardHeader>
            <CardTitle>抓取配置</CardTitle>
          </CardHeader>
          <CardContent class="space-y-4">
            <div class="space-y-3">
              <Label>抓取地区</Label>
              
              <!-- 已选择的地区 (徽章显示) -->
              <div v-if="config.Fetcher.Regions && config.Fetcher.Regions.length > 0" class="flex flex-wrap gap-2 p-3 bg-gray-50 rounded-md">
                <Badge 
                  v-for="regionValue in config.Fetcher.Regions" 
                  :key="regionValue"
                  variant="secondary"
                  class="px-3 py-1.5 text-sm cursor-pointer hover:bg-red-100 hover:text-red-700 transition-colors group"
                  @click="removeRegion(regionValue)"
                >
                  <span>{{ regionValue }}</span>
                  <span class="ml-1.5 text-xs opacity-60 group-hover:opacity-100">✕</span>
                </Badge>
              </div>
              <div v-else class="text-sm text-gray-500 p-3 bg-gray-50 rounded-md">
                未选择任何地区，默认将抓取 zh-CN
              </div>

              <!-- 手动输入地区代码 -->
              <div class="space-y-2">
                <div class="flex gap-2 items-start">
                  <div class="flex-1 space-y-1">
                    <Input 
                      v-model="regionInput"
                      placeholder="输入地区代码，如: zh-CN, en-US, ja-JP"
                      @keypress.enter="addRegionFromInput"
                      :class="{ 'border-red-500': regionInputError }"
                    />
                    <p v-if="regionInputError" class="text-xs text-red-600">
                      {{ regionInputError }}
                    </p>
                  </div>
                  <Button 
                    @click="addRegionFromInput" 
                    :disabled="!regionInput.trim()"
                    variant="outline"
                  >
                    添加
                  </Button>
                </div>
                
                <!-- 常用地区快速添加 -->
                <div class="space-y-1.5">
                  <p class="text-xs text-gray-600 font-medium">常用地区代码：</p>
                  <div class="flex flex-wrap gap-1.5">
                    <Badge
                      v-for="region in commonRegions"
                      :key="region.code"
                      :variant="config.Fetcher.Regions?.includes(region.code) ? 'default' : 'outline'"
                      class="cursor-pointer text-xs px-2 py-1"
                      :class="{
                        'opacity-50 cursor-not-allowed': config.Fetcher.Regions?.includes(region.code),
                        'hover:bg-primary hover:text-primary-foreground': !config.Fetcher.Regions?.includes(region.code)
                      }"
                      @click="addCommonRegion(region.code)"
                      :title="region.name"
                    >
                      {{ region.code }}
                    </Badge>
                  </div>
                </div>
              </div>

              <p class="text-xs text-gray-500">
                输入 Bing 地区代码（格式：语言-地区，如 zh-CN）。点击上方常用代码可快速添加。点击徽章可移除地区。
              </p>
            </div>
          </CardContent>
        </Card>

        <!-- 功能特性配置 -->
        <Card>
          <CardHeader>
            <CardTitle>功能特性</CardTitle>
          </CardHeader>
          <CardContent class="space-y-4">
            <div class="flex items-center gap-2">
              <Label for="feature-write-daily">写入每日文件</Label>
              <Switch
                id="feature-write-daily"
                v-model="config.Feature.WriteDailyFiles"
              />
            </div>
          </CardContent>
        </Card>

        <!-- Web 配置 -->
        <Card>
          <CardHeader>
            <CardTitle>Web 前端配置</CardTitle>
          </CardHeader>
          <CardContent class="space-y-4">
            <div class="space-y-2">
              <Label>前端静态文件路径</Label>
              <Input v-model="config.Web.Path" placeholder="./webapp/dist" />
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { toast } from 'vue-sonner'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Switch } from '@/components/ui/switch'
import { Badge } from '@/components/ui/badge'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { apiService } from '@/lib/api-service'
import type { Config } from '@/lib/api-types'

const editMode = ref<'json' | 'form'>('form')
const loading = ref(false)
const loadError = ref('')
const saveLoading = ref(false)
const dsnError = ref('')

// 地区输入相关
const regionInput = ref('')
const regionInputError = ref('')

// 常用地区代码列表
const commonRegions = [
  { code: 'zh-CN', name: '中国大陆' },
  { code: 'zh-TW', name: '台湾' },
  { code: 'zh-HK', name: '香港' },
  { code: 'en-US', name: '美国' },
  { code: 'en-GB', name: '英国' },
  { code: 'en-CA', name: '加拿大' },
  { code: 'en-AU', name: '澳大利亚' },
  { code: 'en-IN', name: '印度' },
  { code: 'ja-JP', name: '日本' },
  { code: 'ko-KR', name: '韩国' },
  { code: 'de-DE', name: '德国' },
  { code: 'fr-FR', name: '法国' },
  { code: 'es-ES', name: '西班牙' },
  { code: 'it-IT', name: '意大利' },
  { code: 'pt-BR', name: '巴西' },
  { code: 'ru-RU', name: '俄罗斯' },
  { code: 'ar-SA', name: '沙特阿拉伯' },
  { code: 'th-TH', name: '泰国' },
  { code: 'vi-VN', name: '越南' },
  { code: 'id-ID', name: '印度尼西亚' }
]

const config = ref<Config>({
  Admin: { PasswordBcrypt: '' },
  API: { Mode: 'local', EnableMktFallback: true, EnableOnDemandFetch: false },
  Cron: { Enabled: true, DailySpec: '0 9 * * *' },
  DB: { Type: 'sqlite', DSN: '' },
  Feature: { WriteDailyFiles: true },
  Log: {
    Level: 'info',
    Filename: '',
    DBFilename: '',
    DBLogLevel: 'warn',
    LogConsole: true,
    ShowDBLog: false,
    MaxSize: 10,
    MaxAge: 30,
    MaxBackups: 10,
    Compress: true
  },
  Retention: { Days: 30 },
  Server: { Port: 8080, BaseURL: '' },
  Storage: {
    Type: 'local',
    Local: { Root: './data/images' },
    S3: {
      Endpoint: '',
      AccessKey: '',
      SecretKey: '',
      Bucket: '',
      Region: '',
      ForcePathStyle: false,
      PublicURLPrefix: ''
    },
    WebDAV: {
      URL: '',
      Username: '',
      Password: '',
      PublicURLPrefix: ''
    }
  },
  Token: { DefaultTTL: '168h' },
  Web: { Path: './webapp/dist' },
  Fetcher: { Regions: [] }
})

const configJson = ref('')
const jsonError = ref('')

// 验证地区代码格式 (语言-地区格式，如 zh-CN)
const validateRegionCode = (code: string): boolean => {
  // 基本格式验证：2-3个字母 + 连字符 + 2个字母，如 zh-CN, en-US
  const pattern = /^[a-z]{2,3}-[A-Z]{2}$/
  return pattern.test(code)
}

// 从输入框添加地区
const addRegionFromInput = () => {
  const code = regionInput.value.trim()
  regionInputError.value = ''
  
  if (!code) return
  
  // 验证格式
  if (!validateRegionCode(code)) {
    regionInputError.value = '地区代码格式不正确，应为：语言-地区（如 zh-CN, en-US）'
    return
  }
  
  if (!config.value.Fetcher.Regions) {
    config.value.Fetcher.Regions = []
  }
  
  // 检查是否已存在
  if (config.value.Fetcher.Regions.includes(code)) {
    regionInputError.value = '该地区已存在'
    return
  }
  
  config.value.Fetcher.Regions.push(code)
  regionInput.value = ''
  toast.success(`已添加地区: ${code}`)
}

// 快速添加常用地区
const addCommonRegion = (code: string) => {
  if (!config.value.Fetcher.Regions) {
    config.value.Fetcher.Regions = []
  }
  
  if (config.value.Fetcher.Regions.includes(code)) {
    return
  }
  
  config.value.Fetcher.Regions.push(code)
  const region = commonRegions.find(r => r.code === code)
  toast.success(`已添加地区: ${region?.name || code}`)
}

// 移除地区
const removeRegion = (regionValue: string) => {
  if (!config.value.Fetcher.Regions) return
  config.value.Fetcher.Regions = config.value.Fetcher.Regions.filter(r => r !== regionValue)
  toast.success(`已移除地区: ${regionValue}`)
}

// DSN 示例
const dsnExamples = computed(() => {
  switch (config.value.DB.Type) {
    case 'sqlite':
      return 'data/bing_paper.db 或 file:data/bing_paper.db?cache=shared'
    case 'mysql':
      return 'user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True'
    case 'postgres':
      return 'host=localhost user=postgres password=secret dbname=mydb port=5432 sslmode=disable'
    default:
      return ''
  }
})

// 验证 DSN
const validateDSN = () => {
  dsnError.value = ''
  const dsn = config.value.DB.DSN.trim()
  
  if (!dsn) {
    dsnError.value = 'DSN 不能为空'
    return false
  }
  
  switch (config.value.DB.Type) {
    case 'mysql':
      if (!dsn.includes('@tcp(') && !dsn.includes('://')) {
        dsnError.value = 'MySQL DSN 格式不正确，应包含 @tcp( 或使用 URI 格式'
        return false
      }
      break
    case 'postgres':
      if (!dsn.includes('host=') && !dsn.includes('://')) {
        dsnError.value = 'PostgreSQL DSN 格式不正确，应包含 host= 或使用 URI 格式'
        return false
      }
      break
  }
  
  return true
}

// 格式化 JSON
const formatJson = () => {
  try {
    const parsed = JSON.parse(configJson.value)
    configJson.value = JSON.stringify(parsed, null, 2)
    jsonError.value = ''
    toast.success('JSON 格式化成功')
  } catch (err: any) {
    jsonError.value = 'JSON 格式错误: ' + err.message
    toast.error('JSON 格式错误')
  }
}

// 验证 JSON 是否有效
const isValidJson = computed(() => {
  if (!configJson.value.trim()) return false
  try {
    JSON.parse(configJson.value)
    return true
  } catch {
    return false
  }
})

const fetchConfig = async () => {
  loading.value = true
  loadError.value = ''
  try {
    const data = await apiService.getConfig()
    config.value = data
    configJson.value = JSON.stringify(data, null, 2)
  } catch (err: any) {
    loadError.value = err.message || '获取配置失败'
    console.error('获取配置失败:', err)
  } finally {
    loading.value = false
  }
}

// 监听表单变化更新 JSON
watch(config, (newConfig) => {
  if (editMode.value === 'form') {
    configJson.value = JSON.stringify(newConfig, null, 2)
  }
}, { deep: true })

// 监听 JSON 变化更新表单
watch(configJson, (newJson) => {
  if (editMode.value === 'json') {
    try {
      const parsed = JSON.parse(newJson)
      config.value = parsed
      jsonError.value = ''
    } catch (err: any) {
      jsonError.value = err.message
    }
  }
})

const handleSaveConfig = async () => {
  saveLoading.value = true
  
  try {
    // 如果是 JSON 模式，先验证格式
    if (editMode.value === 'json') {
      if (!isValidJson.value) {
        throw new Error('JSON 格式不正确，请检查语法')
      }
      config.value = JSON.parse(configJson.value)
    } else {
      // 表单模式下验证 DSN
      if (!validateDSN()) {
        throw new Error('DSN 格式不正确: ' + dsnError.value)
      }
    }
    
    await apiService.updateConfig(config.value)
    toast.success('配置保存成功')
    
    // 重新加载配置
    await fetchConfig()
  } catch (err: any) {
    toast.error(err.message || '保存配置失败')
    console.error('保存配置失败:', err)
  } finally {
    saveLoading.value = false
  }
}

onMounted(() => {
  fetchConfig()
})
</script>
