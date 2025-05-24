package utils

import (
	"errors"
	"regexp"
	"strings"
)

// ValidateUsername 验证用户名
func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return errors.New("用户名长度必须在3-20个字符之间")
	}

	// 只允许字母、数字和下划线
	matched, err := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	if err != nil || !matched {
		return errors.New("用户名只能包含字母、数字和下划线")
	}

	return nil
}

// ValidatePassword 验证密码
func ValidatePassword(password string) error {
	if len(password) < 6 || len(password) > 20 {
		return errors.New("密码长度必须在6-20个字符之间")
	}

	// 必须包含字母和数字
	hasLetter := false
	hasNumber := false
	for _, char := range password {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			hasLetter = true
		}
		if char >= '0' && char <= '9' {
			hasNumber = true
		}
	}

	if !hasLetter || !hasNumber {
		return errors.New("密码必须包含字母和数字")
	}

	return nil
}

// ValidateEmail 验证邮箱
func ValidateEmail(email string) error {
	if email == "" {
		return nil
	}

	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil || !matched {
		return errors.New("邮箱格式不正确")
	}

	return nil
}

// ValidateFileExt 验证文件扩展名
func ValidateFileExt(filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	if !strings.Contains(AllowedFileExts, ext) {
		return errors.New("不支持的文件类型")
	}
	return nil
}

// ValidateFileSize 验证文件大小
func ValidateFileSize(size int64) error {
	if size > MaxFileSize {
		return errors.New("文件大小超过限制")
	}
	return nil
}

// ValidatePageSize 验证分页大小
func ValidatePageSize(pageSize int) error {
	if pageSize < 1 || pageSize > MaxPageSize {
		return errors.New("分页大小超出范围")
	}
	return nil
}

// ValidateTaskType 验证任务类型
func ValidateTaskType(taskType int) error {
	if taskType != TaskTypeContentGeneration && taskType != TaskTypeTranslation {
		return errors.New("无效的任务类型")
	}
	return nil
}

// ValidateTaskStatus 验证任务状态
func ValidateTaskStatus(status int) error {
	if status < TaskStatusWaiting || status > TaskStatusCanceled {
		return errors.New("无效的任务状态")
	}
	return nil
}

// ValidatePromptType 验证提示词类型
func ValidatePromptType(promptType int) error {
	if promptType != PromptTypeContentGeneration && promptType != PromptTypeTranslation {
		return errors.New("无效的提示词类型")
	}
	return nil
}

// ValidateWorkStatus 验证作品状态
func ValidateWorkStatus(status int) error {
	if status < WorkStatusWaiting || status > WorkStatusCanceled {
		return errors.New("无效的作品状态")
	}
	return nil
}

// ValidateTranslationBatchStatus 验证翻译批次状态
func ValidateTranslationBatchStatus(status int) error {
	if status < TranslationBatchStatusWaiting || status > TranslationBatchStatusCanceled {
		return errors.New("无效的翻译批次状态")
	}
	return nil
}

// ValidateTranslationResultStatus 验证翻译结果状态
func ValidateTranslationResultStatus(status int) error {
	if status < TranslationResultStatusWaiting || status > TranslationResultStatusCanceled {
		return errors.New("无效的翻译结果状态")
	}
	return nil
}

// ValidateUserStatus 验证用户状态
func ValidateUserStatus(status int) error {
	if status != UserStatusActive && status != UserStatusInactive {
		return errors.New("无效的用户状态")
	}
	return nil
}

// ValidateRoleType 验证角色类型
func ValidateRoleType(roleType int) error {
	if roleType != RoleTypeAdmin && roleType != RoleTypeUser {
		return errors.New("无效的角色类型")
	}
	return nil
} 