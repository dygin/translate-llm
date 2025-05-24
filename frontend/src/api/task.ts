import axios from 'axios';
import type { Task, PriorityRule, RuleTemplate, RuleGroup, PriorityLog, TaskStats } from '@/types/task';

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    return response.data;
  },
  (error) => {
    if (error.response?.status === 401) {
      // 处理未授权错误
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// 任务管理API
export const taskApi = {
  // 创建任务
  createTask: (data: Partial<Task>) => {
    return api.post<Task>('/tasks', data);
  },

  // 获取任务
  getTask: (id: string) => {
    return api.get<Task>(`/tasks/${id}`);
  },

  // 获取任务列表
  getTasks: (params: any) => {
    return api.get<{ items: Task[]; total: number }>('/tasks', { params });
  },

  // 更新任务状态
  updateTaskStatus: (id: string, status: string) => {
    return api.put(`/tasks/${id}/status`, { status });
  },

  // 删除任务
  deleteTask: (id: string) => {
    return api.delete(`/tasks/${id}`);
  },

  // 重试任务
  retryTask: (id: string) => {
    return api.post(`/tasks/${id}/retry`);
  },

  // 获取任务统计
  getTaskStats: () => {
    return api.get<TaskStats>('/tasks/stats');
  },

  // 更新任务优先级
  updateTaskPriority: (id: string, priority: number) => {
    return api.put(`/tasks/${id}/priority`, { priority });
  },

  // 批量更新任务优先级
  batchUpdatePriority: (ids: string[], priority: number) => {
    return api.put('/tasks/priority/batch', { ids, priority });
  },

  // 根据条件更新任务优先级
  updatePriorityByCondition: (condition: any, priority: number) => {
    return api.put('/tasks/priority/condition', { condition, priority });
  },
};

// 优先级规则管理API
export const priorityRuleApi = {
  // 创建优先级规则
  createPriorityRule: (data: Partial<PriorityRule>) => {
    return api.post<PriorityRule>('/tasks/rules', data);
  },

  // 获取优先级规则
  getPriorityRule: (id: string) => {
    return api.get<PriorityRule>(`/tasks/rules/${id}`);
  },

  // 更新优先级规则
  updatePriorityRule: (id: string, data: Partial<PriorityRule>) => {
    return api.put<PriorityRule>(`/tasks/rules/${id}`, data);
  },

  // 删除优先级规则
  deletePriorityRule: (id: string) => {
    return api.delete(`/tasks/rules/${id}`);
  },

  // 获取优先级规则列表
  getPriorityRules: (params: any) => {
    return api.get<{ items: PriorityRule[]; total: number }>('/tasks/rules', { params });
  },

  // 获取优先级调整日志
  getPriorityLogs: (taskId: string) => {
    return api.get<PriorityLog[]>(`/tasks/${taskId}/priority-logs`);
  },
};

// 规则模板管理API
export const ruleTemplateApi = {
  // 创建规则模板
  createRuleTemplate: (data: Partial<RuleTemplate>) => {
    return api.post<RuleTemplate>('/tasks/templates', data);
  },

  // 获取规则模板
  getRuleTemplate: (id: string) => {
    return api.get<RuleTemplate>(`/tasks/templates/${id}`);
  },

  // 更新规则模板
  updateRuleTemplate: (id: string, data: Partial<RuleTemplate>) => {
    return api.put<RuleTemplate>(`/tasks/templates/${id}`, data);
  },

  // 删除规则模板
  deleteRuleTemplate: (id: string) => {
    return api.delete(`/tasks/templates/${id}`);
  },

  // 获取规则模板列表
  getRuleTemplates: (params: any) => {
    return api.get<{ items: RuleTemplate[]; total: number }>('/tasks/templates', { params });
  },
};

// 规则组管理API
export const ruleGroupApi = {
  // 创建规则组
  createRuleGroup: (data: Partial<RuleGroup>) => {
    return api.post<RuleGroup>('/tasks/groups', data);
  },

  // 获取规则组
  getRuleGroup: (id: string) => {
    return api.get<RuleGroup>(`/tasks/groups/${id}`);
  },

  // 更新规则组
  updateRuleGroup: (id: string, data: Partial<RuleGroup>) => {
    return api.put<RuleGroup>(`/tasks/groups/${id}`, data);
  },

  // 删除规则组
  deleteRuleGroup: (id: string) => {
    return api.delete(`/tasks/groups/${id}`);
  },

  // 获取规则组列表
  getRuleGroups: (params: any) => {
    return api.get<{ items: RuleGroup[]; total: number }>('/tasks/groups', { params });
  },

  // 添加规则到组
  addRuleToGroup: (groupId: string, ruleId: string) => {
    return api.post(`/tasks/groups/${groupId}/rules`, { rule_id: ruleId });
  },

  // 从组中移除规则
  removeRuleFromGroup: (groupId: string, ruleId: string) => {
    return api.delete(`/tasks/groups/${groupId}/rules/${ruleId}`);
  },

  // 获取组中的规则
  getGroupRules: (groupId: string) => {
    return api.get<PriorityRule[]>(`/tasks/groups/${groupId}/rules`);
  },
}; 