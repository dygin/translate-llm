<template>
  <el-container class="app-container">
    <el-aside width="200px">
      <el-menu
        :router="true"
        :default-active="$route.path"
        class="el-menu-vertical"
      >
        <el-menu-item index="/tasks">
          <el-icon><Document /></el-icon>
          <span>任务列表</span>
        </el-menu-item>
        <el-menu-item index="/rules">
          <el-icon><Setting /></el-icon>
          <span>规则管理</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header>
        <div class="header-content">
          <h2>AI翻译系统</h2>
          <el-dropdown @command="handleCommand">
            <span class="user-dropdown">
              {{ username }}
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      <el-main>
        <router-view></router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { Document, Setting, ArrowDown } from '@element-plus/icons-vue';

const router = useRouter();
const username = ref('管理员');

const handleCommand = (command: string) => {
  if (command === 'logout') {
    localStorage.removeItem('token');
    router.push('/login');
  }
};
</script>

<style scoped>
.app-container {
  height: 100vh;
}

.el-menu-vertical {
  height: 100%;
  border-right: none;
}

.el-header {
  background-color: #fff;
  border-bottom: 1px solid #dcdfe6;
  padding: 0 20px;
}

.header-content {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.user-dropdown {
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 4px;
}

.el-main {
  background-color: #f5f7fa;
  padding: 20px;
}
</style> 