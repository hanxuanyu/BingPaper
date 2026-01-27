<template>
  <div class="min-h-screen bg-gray-50 p-8">
    <div class="max-w-6xl mx-auto">
      <h1 class="text-4xl font-bold text-center mb-8 text-gray-900">
        Shadcn-Vue 组件展示
      </h1>

      <!-- Tabs 布局 -->
      <Tabs default-value="forms" class="w-full">
        <TabsList class="grid w-full grid-cols-5">
          <TabsTrigger value="forms">表单组件</TabsTrigger>
          <TabsTrigger value="layout">布局组件</TabsTrigger>
          <TabsTrigger value="dialogs">弹窗组件</TabsTrigger>
          <TabsTrigger value="feedback">反馈组件</TabsTrigger>
          <TabsTrigger value="data">数据展示</TabsTrigger>
        </TabsList>

        <!-- 表单组件 -->
        <TabsContent value="forms" class="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>表单组件示例</CardTitle>
              <CardDescription>各种表单输入控件的使用示例</CardDescription>
            </CardHeader>
            <CardContent class="space-y-6">
              <!-- 输入框组件 -->
              <div class="space-y-4">
                <div>
                  <Label for="text-input">文本输入框</Label>
                  <Input
                    id="text-input"
                    v-model="formData.textInput"
                    placeholder="请输入文本"
                    class="mt-1"
                  />
                </div>

                <div>
                  <Label for="textarea">文本域</Label>
                  <Textarea
                    id="textarea"
                    v-model="formData.textArea"
                    placeholder="请输入多行文本"
                    class="mt-1"
                  />
                </div>

                <div>
                  <Label for="select">选择器</Label>
                  <Select v-model="formData.selectValue">
                    <SelectTrigger class="mt-1">
                      <SelectValue placeholder="请选择选项" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="option1">选项1</SelectItem>
                      <SelectItem value="option2">选项2</SelectItem>
                      <SelectItem value="option3">选项3</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <!-- 复选框 -->
                <div class="flex items-center space-x-2">
                  <Checkbox id="checkbox" v-model:checked="formData.checkboxValue" />
                  <Label for="checkbox">复选框选项</Label>
                </div>

                <!-- 开关 -->
                <div class="flex items-center space-x-2">
                  <Switch id="switch" v-model:checked="formData.switchValue" />
                  <Label for="switch">开关选项</Label>
                </div>

                <!-- 单选组 -->
                <div>
                  <Label>单选组</Label>
                  <RadioGroup v-model="formData.radioValue" class="mt-2">
                    <div class="flex items-center space-x-2">
                      <RadioGroupItem value="radio1" id="r1" />
                      <Label for="r1">单选1</Label>
                    </div>
                    <div class="flex items-center space-x-2">
                      <RadioGroupItem value="radio2" id="r2" />
                      <Label for="r2">单选2</Label>
                    </div>
                  </RadioGroup>
                </div>
              </div>

              <Button @click="submitForm" class="w-full">
                提交表单
              </Button>
            </CardContent>
          </Card>
        </TabsContent>

        <!-- 布局组件 -->
        <TabsContent value="layout" class="space-y-6">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- 卡片示例 -->
            <Card>
              <CardHeader>
                <CardTitle>卡片标题</CardTitle>
                <CardDescription>这是一个卡片组件的描述</CardDescription>
              </CardHeader>
              <CardContent>
                <p>卡片内容区域，可以放置任意内容。</p>
              </CardContent>
              <CardFooter>
                <Button variant="outline">卡片操作</Button>
              </CardFooter>
            </Card>

            <!-- 头像组件 -->
            <Card>
              <CardHeader>
                <CardTitle>头像组件</CardTitle>
              </CardHeader>
              <CardContent class="flex flex-col space-y-4">
                <div class="flex items-center space-x-4">
                  <Avatar>
                    <AvatarImage src="https://github.com/shadcn.png" />
                    <AvatarFallback>CN</AvatarFallback>
                  </Avatar>
                  <div>
                    <p class="text-sm font-medium">用户名称</p>
                    <p class="text-xs text-muted-foreground">user@example.com</p>
                  </div>
                </div>
                
                <Separator />
                
                <div class="flex items-center space-x-2">
                  <Badge>标签1</Badge>
                  <Badge variant="secondary">标签2</Badge>
                  <Badge variant="outline">标签3</Badge>
                </div>
              </CardContent>
            </Card>
          </div>

          <!-- 进度条 -->
          <Card>
            <CardHeader>
              <CardTitle>进度组件</CardTitle>
            </CardHeader>
            <CardContent class="space-y-4">
              <div>
                <div class="flex justify-between text-sm mb-2">
                  <span>上传进度</span>
                  <span>{{ progress }}%</span>
                </div>
                <Progress :model-value="progress" />
                <div class="flex space-x-2 mt-2">
                  <Button @click="updateProgress(25)" size="sm" variant="outline">25%</Button>
                  <Button @click="updateProgress(50)" size="sm" variant="outline">50%</Button>
                  <Button @click="updateProgress(75)" size="sm" variant="outline">75%</Button>
                  <Button @click="updateProgress(100)" size="sm" variant="outline">100%</Button>
                </div>
              </div>
            </CardContent>
          </Card>

          <!-- 骨架屏 -->
          <Card>
            <CardHeader>
              <CardTitle>骨架屏组件</CardTitle>
            </CardHeader>
            <CardContent>
              <div class="flex items-center space-x-4">
                <Skeleton class="h-12 w-12 rounded-full" />
                <div class="space-y-2">
                  <Skeleton class="h-4 w-[250px]" />
                  <Skeleton class="h-4 w-[200px]" />
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <!-- 弹窗组件 -->
        <TabsContent value="dialogs" class="space-y-6">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- 对话框 -->
            <Card>
              <CardHeader>
                <CardTitle>对话框组件</CardTitle>
              </CardHeader>
              <CardContent class="space-y-4">
                <Dialog>
                  <DialogTrigger as-child>
                    <Button>打开对话框</Button>
                  </DialogTrigger>
                  <DialogContent class="sm:max-w-[425px]">
                    <DialogHeader>
                      <DialogTitle>对话框标题</DialogTitle>
                      <DialogDescription>
                        这是一个对话框的描述信息，用于向用户说明对话框的用途。
                      </DialogDescription>
                    </DialogHeader>
                    <div class="grid gap-4 py-4">
                      <div class="grid grid-cols-4 items-center gap-4">
                        <Label for="name" class="text-right">
                          用户名
                        </Label>
                        <Input id="name" value="张三" class="col-span-3" />
                      </div>
                      <div class="grid grid-cols-4 items-center gap-4">
                        <Label for="username" class="text-right">
                          邮箱
                        </Label>
                        <Input id="username" value="zhangsan@example.com" class="col-span-3" />
                      </div>
                    </div>
                    <DialogFooter>
                      <Button>保存</Button>
                    </DialogFooter>
                  </DialogContent>
                </Dialog>

                <!-- 警告对话框 -->
                <AlertDialog>
                  <AlertDialogTrigger as-child>
                    <Button variant="destructive">删除操作</Button>
                  </AlertDialogTrigger>
                  <AlertDialogContent>
                    <AlertDialogHeader>
                      <AlertDialogTitle>确认删除？</AlertDialogTitle>
                      <AlertDialogDescription>
                        此操作无法撤销。这将永久删除您的数据。
                      </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                      <AlertDialogCancel>取消</AlertDialogCancel>
                      <AlertDialogAction>删除</AlertDialogAction>
                    </AlertDialogFooter>
                  </AlertDialogContent>
                </AlertDialog>
              </CardContent>
            </Card>

            <!-- 侧边栏 -->
            <Card>
              <CardHeader>
                <CardTitle>侧边栏组件</CardTitle>
              </CardHeader>
              <CardContent>
                <Sheet>
                  <SheetTrigger as-child>
                    <Button variant="outline">打开侧边栏</Button>
                  </SheetTrigger>
                  <SheetContent>
                    <SheetHeader>
                      <SheetTitle>侧边栏标题</SheetTitle>
                      <SheetDescription>
                        这是侧边栏的描述信息。
                      </SheetDescription>
                    </SheetHeader>
                    <div class="grid gap-4 py-4">
                      <div class="grid grid-cols-4 items-center gap-4">
                        <Label for="sidebar-name" class="text-right">
                          名称
                        </Label>
                        <Input id="sidebar-name" value="示例" class="col-span-3" />
                      </div>
                    </div>
                    <SheetFooter>
                      <SheetClose as-child>
                        <Button>保存更改</Button>
                      </SheetClose>
                    </SheetFooter>
                  </SheetContent>
                </Sheet>
              </CardContent>
            </Card>
          </div>

          <!-- 命令面板 -->
          <Card>
            <CardHeader>
              <CardTitle>命令面板</CardTitle>
            </CardHeader>
            <CardContent>
              <CommandDialog :open="commandOpen" @update:open="commandOpen = $event">
                <CommandInput placeholder="输入命令或搜索..." />
                <CommandList>
                  <CommandEmpty>未找到相关命令。</CommandEmpty>
                  <CommandGroup heading="建议">
                    <CommandItem value="calendar">
                      <span>日历</span>
                    </CommandItem>
                    <CommandItem value="search-emoji">
                      <span>搜索表情</span>
                    </CommandItem>
                    <CommandItem value="calculator">
                      <span>计算器</span>
                    </CommandItem>
                  </CommandGroup>
                </CommandList>
              </CommandDialog>
              <Button @click="commandOpen = true">
                打开命令面板 ⌘K
              </Button>
            </CardContent>
          </Card>
        </TabsContent>

        <!-- 反馈组件 -->
        <TabsContent value="feedback" class="space-y-6">
          <!-- 基础 Toast 示例 -->
          <Card>
            <CardHeader>
              <CardTitle>基础 Toast 提示</CardTitle>
              <CardDescription>使用 Sonner 提供各种类型的提示信息</CardDescription>
            </CardHeader>
            <CardContent class="space-y-4">
              <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-2">
                <Button @click="showDefaultToast" variant="outline">默认</Button>
                <Button @click="showSuccessToast" variant="default">成功</Button>
                <Button @click="showInfoToast" variant="secondary">信息</Button>
                <Button @click="showWarningToast" variant="secondary">警告</Button>
                <Button @click="showErrorToast" variant="destructive">错误</Button>
                <Button @click="showPromiseToast" variant="outline">Promise</Button>
              </div>
            </CardContent>
          </Card>

          <!-- 高级 Toast 示例 -->
          <Card>
            <CardHeader>
              <CardTitle>高级 Toast 功能</CardTitle>
              <CardDescription>包含动作按钮、富文本内容和自定义样式的 Toast</CardDescription>
            </CardHeader>
            <CardContent class="space-y-4">
              <div class="grid grid-cols-2 md:grid-cols-4 gap-2">
                <Button @click="showActionToast" variant="outline">带操作</Button>
                <Button @click="showRichToast" variant="outline">富文本</Button>
                <Button @click="showCustomToast" variant="outline">自定义样式</Button>
                <Button @click="showPersistentToast" variant="outline">持久化</Button>
              </div>
            </CardContent>
          </Card>

          <!-- Toast 位置示例 -->
          <Card>
            <CardHeader>
              <CardTitle>Toast 位置设置</CardTitle>
              <CardDescription>
                在不同位置显示 Toast 通知（需要在Toaster组件上配置position属性）
              </CardDescription>
            </CardHeader>
            <CardContent class="space-y-4">
              <div class="grid grid-cols-2 md:grid-cols-4 gap-2">
                <Button @click="() => showPositionToast('top-left')" variant="outline">左上角消息</Button>
                <Button @click="() => showPositionToast('top-right')" variant="outline">右上角消息</Button>
                <Button @click="() => showPositionToast('bottom-left')" variant="outline">左下角消息</Button>
                <Button @click="() => showPositionToast('bottom-right')" variant="outline">右下角消息</Button>
              </div>
              <div class="grid grid-cols-2 gap-2">
                <Button @click="() => showPositionToast('top-center')" variant="outline">顶部中央消息</Button>
                <Button @click="() => showPositionToast('bottom-center')" variant="outline">底部中央消息</Button>
              </div>
              <p class="text-sm text-muted-foreground">
                注意：当前所有Toast都显示在默认位置。要改变位置需要在Toaster组件上设置position属性。
              </p>
            </CardContent>
          </Card>

          <!-- 批量操作示例 -->
          <Card>
            <CardHeader>
              <CardTitle>批量 Toast 操作</CardTitle>
              <CardDescription>管理多个 Toast 通知</CardDescription>
            </CardHeader>
            <CardContent class="space-y-4">
              <div class="flex flex-wrap gap-2">
                <Button @click="showMultipleToasts" variant="outline">显示多个</Button>
                <Button @click="dismissAllToasts" variant="destructive">关闭所有</Button>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <!-- 数据展示 -->
        <TabsContent value="data" class="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>数据表格</CardTitle>
              <CardDescription>展示结构化数据的表格组件</CardDescription>
            </CardHeader>
            <CardContent>
              <div class="rounded-md border">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead class="w-[100px]">ID</TableHead>
                      <TableHead>姓名</TableHead>
                      <TableHead>邮箱</TableHead>
                      <TableHead>状态</TableHead>
                      <TableHead class="text-right">金额</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    <TableRow v-for="user in tableData" :key="user.id">
                      <TableCell class="font-medium">{{ user.id }}</TableCell>
                      <TableCell>{{ user.name }}</TableCell>
                      <TableCell>{{ user.email }}</TableCell>
                      <TableCell>
                        <Badge :variant="user.status === '激活' ? 'default' : 'secondary'">
                          {{ user.status }}
                        </Badge>
                      </TableCell>
                      <TableCell class="text-right">{{ user.amount }}</TableCell>
                    </TableRow>
                  </TableBody>
                </Table>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>

    <!-- Toast 容器 -->
    <Toaster/>
    <!-- <Toaster 
      :close-button="true"
      :rich-colors="true"
      :expand="true"
    /> -->
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { toast } from 'vue-sonner'

// 组件导入
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { 
  Select, 
  SelectContent, 
  SelectItem, 
  SelectTrigger, 
  SelectValue 
} from '@/components/ui/select'
import { Checkbox } from '@/components/ui/checkbox'
import { Switch } from '@/components/ui/switch'
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group'
import { 
  Card, 
  CardContent, 
  CardDescription, 
  CardFooter, 
  CardHeader, 
  CardTitle 
} from '@/components/ui/card'
import { 
  Tabs, 
  TabsContent, 
  TabsList, 
  TabsTrigger 
} from '@/components/ui/tabs'
import { 
  Dialog, 
  DialogContent, 
  DialogDescription, 
  DialogFooter, 
  DialogHeader, 
  DialogTitle, 
  DialogTrigger 
} from '@/components/ui/dialog'
import { 
  AlertDialog, 
  AlertDialogAction, 
  AlertDialogCancel, 
  AlertDialogContent, 
  AlertDialogDescription, 
  AlertDialogFooter, 
  AlertDialogHeader, 
  AlertDialogTitle, 
  AlertDialogTrigger 
} from '@/components/ui/alert-dialog'
import { 
  Sheet, 
  SheetClose, 
  SheetContent, 
  SheetDescription, 
  SheetFooter, 
  SheetHeader, 
  SheetTitle, 
  SheetTrigger 
} from '@/components/ui/sheet'
import { 
  CommandDialog, 
  CommandEmpty, 
  CommandGroup, 
  CommandInput, 
  CommandItem, 
  CommandList 
} from '@/components/ui/command'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { Skeleton } from '@/components/ui/skeleton'
import { Separator } from '@/components/ui/separator'
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableHead, 
  TableHeader, 
  TableRow 
} from '@/components/ui/table'
import { Toaster } from '@/components/ui/sonner'

// 响应式数据
const formData = reactive({
  textInput: '',
  textArea: '',
  selectValue: '',
  checkboxValue: false,
  switchValue: false,
  radioValue: 'radio1'
})

const progress = ref(33)
const commandOpen = ref(false)

const tableData = ref([
  { id: '001', name: '张三', email: 'zhangsan@example.com', status: '激活', amount: '¥2,500' },
  { id: '002', name: '李四', email: 'lisi@example.com', status: '禁用', amount: '¥1,800' },
  { id: '003', name: '王五', email: 'wangwu@example.com', status: '激活', amount: '¥3,200' }
])

// 方法
const submitForm = () => {
  console.log('表单数据:', formData)
  toast('表单提交成功！', {
    description: '您的表单数据已成功提交处理',
  })
}

const updateProgress = (value: number) => {
  progress.value = value
}

// 基础 Toast 方法
const showDefaultToast = () => {
  toast('默认 Toast 通知', {
    description: '这是一个基础的 Toast 通知消息。'
  })
}

const showSuccessToast = () => {
  toast.success('操作成功！', {
    description: '您的操作已经成功完成。',
    action: {
      label: '查看详情',
      onClick: () => console.log('查看详情'),
    },
  })
}

const showErrorToast = () => {
  toast.error('操作失败！', {
    description: '请检查您的输入并重试。',
    action: {
      label: '重试',
      onClick: () => console.log('重试操作'),
    },
  })
}

const showWarningToast = () => {
  toast.warning('警告信息', {
    description: '这是一个警告提示信息。',
    action: {
      label: '了解更多',
      onClick: () => console.log('了解更多'),
    },
  })
}

const showInfoToast = () => {
  toast.info('信息提示', {
    description: '这是一个普通的信息提示。',
    action: {
      label: '关闭',
      onClick: () => console.log('关闭提示'),
    },
  })
}

const showPromiseToast = () => {
  toast.promise<{ name: string }>(
    () =>
      new Promise(resolve =>
        setTimeout(() => resolve({ name: '数据' }), 2000),
      ),
    {
      loading: '正在加载...',
      success: (data: { name: string }) => `${data.name} 加载完成`,
      error: '加载失败',
    },
  )
}

// 高级 Toast 方法
const showActionToast = () => {
  toast('文件上传完成', {
    description: '您的文件已成功上传到云端存储。',
    duration: 5000,
    action: {
      label: '查看文件',
      onClick: () => {
        toast.success('正在打开文件...')
      },
    },
  })
}

const showRichToast = () => {
  toast('新消息通知', {
    description: '您有 3 条未读消息等待处理，请及时查看。',
    duration: 6000,
    action: {
      label: '立即查看',
      onClick: () => console.log('查看消息'),
    },
  })
}

const showCustomToast = () => {
  toast('自定义样式 Toast', {
    description: '这个 Toast 具有自定义的样式和行为。',
    duration: 4000,
    action: {
      label: '自定义操作',
      onClick: () => console.log('自定义操作'),
    },
  })
}

const showPersistentToast = () => {
  toast('持久化通知', {
    description: '这个通知将保持显示，直到用户手动关闭。',
    duration: Infinity,
    action: {
      label: '手动关闭',
      onClick: () => console.log('手动关闭'),
    },
  })
}

// 位置相关方法
const showPositionToast = (position: string) => {
  toast(`${getPositionName(position)} Toast`, {
    description: `这是显示在${getPositionName(position)}的 Toast 通知。`,
    duration: 3000,
    // 注意：vue-sonner 需要在 Toaster 组件上配置 position
    // 这里主要用于展示不同的消息内容
  })
}

const getPositionName = (position: string) => {
  const positionMap: { [key: string]: string } = {
    'top-left': '左上角',
    'top-right': '右上角',
    'bottom-left': '左下角',
    'bottom-right': '右下角',
    'top-center': '顶部中央',
    'bottom-center': '底部中央'
  }
  return positionMap[position] || position
}

// 批量操作方法
const showMultipleToasts = () => {
  setTimeout(() => toast.success('第一个成功消息'), 0)
  setTimeout(() => toast.info('第二个信息消息'), 500)
  setTimeout(() => toast.warning('第三个警告消息'), 1000)
  setTimeout(() => toast.error('第四个错误消息'), 1500)
}

const dismissAllToasts = () => {
  toast.dismiss()
}
</script>

<style>
#app {
  font-family: ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

html, body {
  margin: 0;
  padding: 0;
}

/* 自定义 Toast 样式 */
.custom-toast {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%) !important;
  color: white !important;
  border: none !important;
}

/* Toast 动画增强 */
[data-sonner-toast] {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 不同类型的 Toast 样式增强 */
[data-type="success"] {
  border-left: 4px solid #10b981;
}

[data-type="error"] {
  border-left: 4px solid #ef4444;
}

[data-type="warning"] {
  border-left: 4px solid #f59e0b;
}

[data-type="info"] {
  border-left: 4px solid #3b82f6;
}
</style>