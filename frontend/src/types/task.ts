// 任务类型定义
export interface Task {
  id: string;
  work_id: string;
  batch_id: string;
  type: TaskType;
  status: TaskStatus;
  priority: number;
  content: string;
  result: string;
  error: string;
  created_at: string;
  updated_at: string;
}

// 任务类型枚举
export enum TaskType {
  ContentGeneration = 'content_generation',
  Translation = 'translation'
}

// 任务状态枚举
export enum TaskStatus {
  Pending = 'pending',
  Processing = 'processing',
  Completed = 'completed',
  Failed = 'failed'
}

// 优先级规则类型定义
export interface PriorityRule {
  id: string;
  name: string;
  description: string;
  conditions: RuleCondition[];
  actions: RuleAction[];
  enabled: boolean;
  created_at: string;
  updated_at: string;
}

// 规则条件类型定义
export interface RuleCondition {
  field: string;
  operator: string;
  value: any;
}

// 规则动作类型定义
export interface RuleAction {
  type: string;
  value: any;
}

// 规则模板类型定义
export interface RuleTemplate {
  id: string;
  name: string;
  description: string;
  conditions: RuleCondition[];
  actions: RuleAction[];
  created_at: string;
  updated_at: string;
}

// 规则组类型定义
export interface RuleGroup {
  id: string;
  name: string;
  description: string;
  rules: PriorityRule[];
  enabled: boolean;
  created_at: string;
  updated_at: string;
}

// 优先级日志类型定义
export interface PriorityLog {
  id: string;
  task_id: string;
  old_priority: number;
  new_priority: number;
  reason: string;
  created_at: string;
}

// 任务统计类型定义
export interface TaskStats {
  total: number;
  pending: number;
  processing: number;
  completed: number;
  failed: number;
  by_type: {
    [key in TaskType]: number;
  };
  by_priority: {
    [key: number]: number;
  };
} 