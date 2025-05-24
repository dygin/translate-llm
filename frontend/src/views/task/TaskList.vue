<template>
  <div class="task-list">
    <!-- 搜索和过滤 -->
    <div class="filter-container">
      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="任务类型">
          <el-select v-model="filterForm.type" placeholder="请选择任务类型" clearable>
            <el-option
              v-for="type in taskTypes"
              :key="type.value"
              :label="type.label"
              :value="type.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="任务状态">
          <el-select v-model="filterForm.status" placeholder="请选择任务状态" clearable>
            <el-option
              v-for="status in taskStatuses"
              :key="status.value"
              :label="status.label"
              :value="status.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="优先级">
          <el-select v-model="filterForm.priority" placeholder="请选择优先级" clearable>
            <el-option
              v-for="priority in priorities"
              :key="priority.value"
              :label="priority.label"
              :value="priority.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 批量操作 -->
    <div class="batch-actions">
      <el-button type="primary" @click="handleBatchUpdatePriority">批量更新优先级</el-button>
      <el-button type="danger" @click="handleBatchDelete">批量删除</el-button>
    </div>

    <!-- 任务列表 -->
    <el-table
      v-loading="loading"
      :data="tasks"
      @selection-change="handleSelectionChange"
      border
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="id" label="ID" width="220" />
      <el-table-column prop="work_id" label="工作ID" width="220" />
      <el-table-column prop="batch_id" label="批次ID" width="220" />
      <el-table-column prop="type" label="任务类型" width="120">
        <template #default="{ row }">
          <el-tag :type="getTaskTypeTag(row.type)">
            {{ getTaskTypeLabel(row.type) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusTag(row.status)">
            {{ getStatusLabel(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="priority" label="优先级" width="100">
        <template #default="{ row }">
          <el-tag :type="getPriorityTag(row.priority)">
            {{ row.priority }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180" />
      <el-table-column label="操作" fixed="right" width="200">
        <template #default="{ row }">
          <el-button-group>
            <el-button size="small" @click="handleView(row)">查看</el-button>
            <el-button
              size="small"
              type="primary"
              @click="handleUpdatePriority(row)"
            >更新优先级</el-button>
            <el-button
              size="small"
              type="warning"
              @click="handleRetry(row)"
              :disabled="row.status !== 'failed'"
            >重试</el-button>
            <el-button
              size="small"
              type="danger"
              @click="handleDelete(row)"
            >删除</el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div class="pagination-container">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- 更新优先级对话框 -->
    <el-dialog
      v-model="priorityDialogVisible"
      :title="isBatch ? '批量更新优先级' : '更新优先级'"
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
import { ref, reactive, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import type { Task, TaskType, TaskStatus } from '@/types/task';
import { taskApi } from '@/api/task';

// 状态定义
const loading = ref(false);
const tasks = ref<Task[]>([]);
const total = ref(0);
const currentPage = ref(1);
const pageSize = ref(20);
const selectedTasks = ref<Task[]>([]);
const priorityDialogVisible = ref(false);
const isBatch = ref(false);
const currentTask = ref<Task | null>(null);

// 过滤表单
const filterForm = reactive({
  type: '',
  status: '',
  priority: '',
});

// 优先级表单
const priorityForm = reactive({
  priority: 0,
});

// 任务类型选项
const taskTypes = [
  { label: '内容生成', value: TaskType.ContentGeneration },
  { label: '翻译', value: TaskType.Translation },
];

// 任务状态选项
const taskStatuses = [
  { label: '待处理', value: TaskStatus.Pending },
  { label: '处理中', value: TaskStatus.Processing },
  { label: '已完成', value: TaskStatus.Completed },
  { label: '失败', value: TaskStatus.Failed },
];

// 优先级选项
const priorities = [
  { label: '低', value: 0 },
  { label: '中', value: 1 },
  { label: '高', value: 2 },
  { label: '紧急', value: 3 },
];

// 获取任务列表
const fetchTasks = async () => {
  loading.value = true;
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      ...filterForm,
    };
    const response = await taskApi.getTasks(params);
    tasks.value = response.items;
    total.value = response.total;
  } catch (error) {
    ElMessage.error('获取任务列表失败');
  } finally {
    loading.value = false;
  }
};

// 搜索
const handleSearch = () => {
  currentPage.value = 1;
  fetchTasks();
};

// 重置过滤
const resetFilter = () => {
  Object.keys(filterForm).forEach(key => {
    filterForm[key] = '';
  });
  handleSearch();
};

// 选择变化
const handleSelectionChange = (selection: Task[]) => {
  selectedTasks.value = selection;
};

// 查看任务
const handleView = (task: Task) => {
  // TODO: 实现查看任务详情
};

// 更新优先级
const handleUpdatePriority = (task: Task) => {
  isBatch.value = false;
  currentTask.value = task;
  priorityForm.priority = task.priority;
  priorityDialogVisible.value = true;
};

// 批量更新优先级
const handleBatchUpdatePriority = () => {
  if (selectedTasks.value.length === 0) {
    ElMessage.warning('请选择要更新的任务');
    return;
  }
  isBatch.value = true;
  priorityForm.priority = 0;
  priorityDialogVisible.value = true;
};

// 提交优先级更新
const submitPriorityUpdate = async () => {
  try {
    if (isBatch.value) {
      const ids = selectedTasks.value.map(task => task.id);
      await taskApi.batchUpdatePriority(ids, priorityForm.priority);
      ElMessage.success('批量更新优先级成功');
    } else if (currentTask.value) {
      await taskApi.updateTaskPriority(currentTask.value.id, priorityForm.priority);
      ElMessage.success('更新优先级成功');
    }
    priorityDialogVisible.value = false;
    fetchTasks();
  } catch (error) {
    ElMessage.error('更新优先级失败');
  }
};

// 重试任务
const handleRetry = async (task: Task) => {
  try {
    await taskApi.retryTask(task.id);
    ElMessage.success('重试任务成功');
    fetchTasks();
  } catch (error) {
    ElMessage.error('重试任务失败');
  }
};

// 删除任务
const handleDelete = async (task: Task) => {
  try {
    await ElMessageBox.confirm('确定要删除该任务吗？', '提示', {
      type: 'warning',
    });
    await taskApi.deleteTask(task.id);
    ElMessage.success('删除任务成功');
    fetchTasks();
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除任务失败');
    }
  }
};

// 批量删除
const handleBatchDelete = async () => {
  if (selectedTasks.value.length === 0) {
    ElMessage.warning('请选择要删除的任务');
    return;
  }
  try {
    await ElMessageBox.confirm(`确定要删除选中的 ${selectedTasks.value.length} 个任务吗？`, '提示', {
      type: 'warning',
    });
    await Promise.all(selectedTasks.value.map(task => taskApi.deleteTask(task.id)));
    ElMessage.success('批量删除成功');
    fetchTasks();
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败');
    }
  }
};

// 分页处理
const handleSizeChange = (val: number) => {
  pageSize.value = val;
  fetchTasks();
};

const handleCurrentChange = (val: number) => {
  currentPage.value = val;
  fetchTasks();
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
  return taskTypes.find(t => t.value === type)?.label || type;
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
  return taskStatuses.find(s => s.value === status)?.label || status;
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
  fetchTasks();
});
</script>

<style scoped>
.task-list {
  padding: 20px;
}

.filter-container {
  margin-bottom: 20px;
}

.batch-actions {
  margin-bottom: 20px;
}

.pagination-container {
  margin-top: 20px;
  text-align: right;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style> 