package utils

// 任务类型
const (
	TaskTypeContentGeneration = 1 // 内容生成
	TaskTypeTranslation       = 2 // 翻译
)

// 任务状态
const (
	TaskStatusWaiting  = 0 // 等待中
	TaskStatusRunning  = 1 // 运行中
	TaskStatusPaused   = 2 // 已暂停
	TaskStatusFailed   = 3 // 失败
	TaskStatusSuccess  = 4 // 成功
	TaskStatusCanceled = 5 // 已取消
)

// 提示词类型
const (
	PromptTypeContentGeneration = 1 // 内容生成
	PromptTypeTranslation       = 2 // 翻译
)

// 作品状态
const (
	WorkStatusWaiting  = 0 // 等待中
	WorkStatusRunning  = 1 // 运行中
	WorkStatusSuccess  = 2 // 成功
	WorkStatusFailed   = 3 // 失败
	WorkStatusCanceled = 4 // 已取消
)

// 翻译批次状态
const (
	TranslationBatchStatusWaiting  = 0 // 等待中
	TranslationBatchStatusRunning  = 1 // 运行中
	TranslationBatchStatusSuccess  = 2 // 成功
	TranslationBatchStatusFailed   = 3 // 失败
	TranslationBatchStatusCanceled = 4 // 已取消
)

// 翻译结果状态
const (
	TranslationResultStatusWaiting  = 0 // 等待中
	TranslationResultStatusSuccess  = 1 // 成功
	TranslationResultStatusFailed   = 2 // 失败
	TranslationResultStatusCanceled = 3 // 已取消
)

// 用户状态
const (
	UserStatusActive   = 1 // 正常
	UserStatusInactive = 0 // 禁用
)

// 角色类型
const (
	RoleTypeAdmin = 1 // 管理员
	RoleTypeUser  = 2 // 普通用户
)

// 缓存键前缀
const (
	CacheKeyUserToken = "user:token:" // 用户token缓存键前缀
	CacheKeyTaskLock  = "task:lock:"  // 任务锁缓存键前缀
)

// 缓存过期时间（秒）
const (
	CacheExpireUserToken = 86400 // 用户token过期时间：24小时
	CacheExpireTaskLock  = 300   // 任务锁过期时间：5分钟
)

// 分页默认值
const (
	DefaultPageSize = 10  // 默认每页数量
	MaxPageSize     = 100 // 最大每页数量
)

// 文件上传限制
const (
	MaxFileSize     = 100 << 20 // 最大文件大小：100MB
	AllowedFileExts = ".mp4,.srt,.txt" // 允许的文件扩展名
)

// 重试配置
const (
	MaxRetryCount = 3    // 最大重试次数
	RetryInterval = 5000 // 重试间隔（毫秒）
) 