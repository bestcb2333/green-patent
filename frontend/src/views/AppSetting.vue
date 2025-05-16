<script setup lang="ts">
import usePersistedStore from '@/stores/persisted';


const persisted = usePersistedStore()
</script>

<template>
  <div class="h-full p-2">
    <el-card class="h-full">
      <el-form>
        <el-form-item label="后端服务器地址">
          <el-input v-model="persisted.setting.apiAddr" />
        </el-form-item>
        <el-form-item label="专利文件服务器地址">
          <el-input v-model="persisted.setting.patentAddr" />
        </el-form-item>
        <el-form-item label="年报文件服务器地址">
          <el-input v-model="persisted.setting.reportAddr" />
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>
