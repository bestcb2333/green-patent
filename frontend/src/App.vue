<script setup lang="ts">
import {BellFilled, Hide, View} from '@element-plus/icons-vue';
import {ref} from 'vue';

const isCollapse = ref(false)
</script>

<template>
  <div class="h-screen flex flex-col">

    <el-menu mode="horizontal" :ellipsis="false">
      <el-menu-item>
        企业绿色专利与创新信息系统
      </el-menu-item>
      <el-menu-item class="!ml-auto">
        <el-icon>
          <BellFilled />
        </el-icon>
      </el-menu-item>
    </el-menu>

    <div class="grow min-h-0 container mx-auto flex">

      <el-menu router class="shrink-0 flex flex-col" :collapse="isCollapse"
        :default-active="$route.meta.path as string"
      >
        <el-menu-item>
          <template #title>
            主菜单
          </template>
        </el-menu-item>
        <el-menu-item v-for="item in $router.getRoutes().filter(route => route.meta.path)"
          :key="item.name" :index="item.meta.path as string"
        >
          <el-icon>
            <component :is="item.meta.icon" />
          </el-icon>
          <template #title>
            {{item.meta.name}}
          </template>
        </el-menu-item>
        <el-menu-item class="mt-auto" @click="isCollapse=!isCollapse">
          <el-icon>
            <Hide v-if="isCollapse" />
            <View v-else />
          </el-icon>
          <template #title>
            折叠菜单栏
          </template>
        </el-menu-item>
      </el-menu>

      <div class="grow min-w-0">
        <router-view />
      </div>

    </div>

  </div>
</template>
