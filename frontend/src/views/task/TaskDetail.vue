<template>
  <div class="task-detail">
    <div class="header">
      <h2>任务详情</h2>
      <div class="actions">
        <el-button-group>
          <el-button
            type="primary"
            @click="handleUpdatePriority"
          >更新优先级</el-button>
          <el-button
            type="warning"
            @click="handleRetry"
            :disabled="task?.status !== 'failed'"
          >重试</el-button>
          <el-button
            type="danger"
            @click="handleDelete"
          >删除</el-button>
        </el-button-group>
      </div>
    </div>

    <el-descriptions
      v-if="task"
      :column="2"
      border
    >
      <el-descriptions-item label="任务ID">{{ task.id }}</el-descriptions-item>
      <el-descriptions-item label="工作ID">{{ task.work_id }}</el-descriptions-item>
      <el-descriptions-item label="批次ID">{{ task.batch_id }}</el-descriptions-item>
      <el-descriptions-item label="任务类型">
        <el-tag :type="getTaskTypeTag(task.type)">
          {{ getTaskTypeLabel(task.type) }}
        </el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="任务状态">
        <el-tag :type="getStatusTag(task.status)">
          {{ getStatusLabel(task.status) }}
        </el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="优先级">
        <el-tag :type="getPriorityTag(task.priority)">
          {{ task.priority }}
        </el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="创建时间">{{ task.created_at }}</el-descriptions-item>
      <el-descriptions-item label="更新时间">{{ task.updated_at }}</el-descriptions-item>
    </el-descriptions>

    <div class="content-section">
      <h3>任务内容</h3>
      <el-card class="content-card">
        <pre>{{ task?.content }}</pre>
      </el-card>
    </div>

    <div class="content-section">
      <h3>处理结果</h3>
      <el-card class="content-card">
        <pre>{{ task?.result }}</pre>
      </el-card>
    </div>

    <div v-if="task?.error" class="content-section">
      <h3>错误信息</h3>
      <el-card class="content-card error">
        <pre>{{ task.error }}</pre>
      </el-card>
    </div>

    <div class="content-section">
      <h3>优先级调整日志</h3>
      <el-table
        :data="priorityLogs"
        border
      >
        <el-table-column prop="old_priority" label="原优先级" width="100">
          <template #default="{ row }">
            <el-tag :type="getPriorityTag(row.old_priority)">
              {{ row.old_priority }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="new_priority" label="新优先级" width="100">
          <template #default="{ row }">
            <el-tag :type="getPriorityTag(row.new_priority)">
              {{ row.new_priority }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="reason" label="调整原因" />
        <el-table-column prop="created_at" label="调整时间" width="180" />
      </el-table>
    </div>

    <!-- 更新优先级对话框 -->
    <el-dialog
      v-model="priorityDialogVisible"
      title="更新优先级"
      width="400px"
    >
      <el-form :model="priorityForm" label-width="80px">
        <el-form-item label="优先级">
          <el-select v-model="priorityForm.priority" placeholder="请选择优先级">
            <el-option
              v-for="priority in priorities"
              :key="priority.value"
              :label="priority.label"
              :value="priority.value"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="priorityDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitPriorityUpdate">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import type { Task, TaskType, TaskStatus, PriorityLog } from '@/types/task';
import { taskApi, priorityRuleApi } from '@/api/task';

const route = useRoute();
const router = useRouter();

// 状态定义
const task = ref<Task | null>(null);
const priorityLogs = ref<PriorityLog[]>([]);
const priorityDialogVisible = ref(false);
const priorityForm = ref({
  priority: 0,
});

// 优先级选项
const priorities = [
  { label: '低', value: 0 },
  { label: '中', value: 1 },
  { label: '高', value: 2 },
  { label: '紧急', value: 3 },
];

// 获取任务详情
const fetchTaskDetail = async () => {
  const taskId = route.params.id as string;
  try {
    const [taskResponse, logsResponse] = await Promise.all([
      taskApi.getTask(taskId),
      priorityRuleApi.getPriorityLogs(taskId),
    ]);
    task.value = taskResponse;
    priorityLogs.value = logsResponse;
  } catch (error) {
    ElMessage.error('获取任务详情失败');
    router.push('/tasks');
  }
};

// 更新优先级
const handleUpdatePriority = () => {
  if (!task.value) return;
  priorityForm.value.priority = task.value.priority;
  priorityDialogVisible.value = true;
};

// 提交优先级更新
const submitPriorityUpdate = async () => {
  if (!task.value) return;
  try {
    await taskApi.updateTaskPriority(task.value.id, priorityForm.value.priority);
    ElMessage.success('更新优先级成功');
    priorityDialogVisible.value = false;
    fetchTaskDetail();
  } catch (error) {
    ElMessage.error('更新优先级失败');
  }
};

// 重试任务
const handleRetry = async () => {
  if (!task.value) return;
  try {
    await taskApi.retryTask(task.value.id);
    ElMessage.success('重试任务成功');
    fetchTaskDetail();
  } catch (error) {
    ElMessage.error('重试任务失败');
  }
};

// 删除任务
const handleDelete = async () => {
  if (!task.value) return;
  try {
    await ElMessageBox.confirm('确定要删除该任务吗？', '提示', {
      type: 'warning',
    });
    await taskApi.deleteTask(task.value.id);
    ElMessage.success('删除任务成功');
    router.push('/tasks');
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除任务失败');
    }
  }
};

// 工具函数
const getTaskTypeTag = (type: TaskType) => {
  switch (type) {
    case TaskType.ContentGeneration:
      return 'success';
    case TaskType.Translation:
      return 'info';
    default:
      return '';
  }
};

const getTaskTypeLabel = (type: TaskType) => {
  switch (type) {
    case TaskType.ContentGeneration:
      return '内容生成';
    case TaskType.Translation:
      return '翻译';
    default:
      return type;
  }
};

const getStatusTag = (status: TaskStatus) => {
  switch (status) {
    case TaskStatus.Pending:
      return 'info';
    case TaskStatus.Processing:
      return 'warning';
    case TaskStatus.Completed:
      return 'success';
    case TaskStatus.Failed:
      return 'danger';
    default:
      return '';
  }
};

const getStatusLabel = (status: TaskStatus) => {
  switch (status) {
    case TaskStatus.Pending:
      return '待处理';
    case TaskStatus.Processing:
      return '处理中';
    case TaskStatus.Completed:
      return '已完成';
    case TaskStatus.Failed:
      return '失败';
    default:
      return status;
  }
};

const getPriorityTag = (priority: number) => {
  switch (priority) {
    case 0:
      return 'info';
    case 1:
      return '';
    case 2:
      return 'warning';
    case 3:
      return 'danger';
    default:
      return '';
  }
};

// 初始化
onMounted(() => {
  fetchTaskDetail();
});
</script>

<style scoped>
.task-detail {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.content-section {
  margin-top: 20px;
}

.content-card {
  margin-top: 10px;
}

.content-card pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.content-card.error {
  background-color: #fef0f0;
  color: #f56c6c;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style> 