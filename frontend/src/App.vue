<script setup lang="ts">
import {Hide, View} from '@element-plus/icons-vue';
import {reactive, ref} from 'vue';
import useSessionStore from './stores/session';
import {request} from './axios';
import usePersistedStore from './stores/persisted';
import CryptoJS from 'crypto-js';

const persisted = usePersistedStore()
const session = useSessionStore()

const isCollapse = ref(false)
const isDialogOpen = ref(false)
const currentTab = ref('login')

const loginForm = reactive({
  username: '',
  password: '',
})

async function login() {
  try {
    persisted.token = await request.post<any, string>('/login', {
      username: loginForm.username,
      password: CryptoJS.SHA256(loginForm.password).toString(CryptoJS.enc.Hex),
    })
    await session.loadUser()
    isDialogOpen.value = false
  } catch {}
}

const signupForm = reactive({
  email: '',
  authcode: '',
  username: '',
  password: '',
})

async function signup() {
  try {
    persisted.token = await request.post<any, string>('/signup', {
      username: signupForm.username,
      password: CryptoJS.SHA256(signupForm.password).toString(CryptoJS.enc.Hex),
      email: signupForm.email,
      authcode: signupForm.authcode,
    })
    await session.loadUser()
    isDialogOpen.value = false
  } catch {}
}

const retrieveForm = reactive({
  email: '',
  authcode: '',
  password: '',
})

async function retrieve() {
  try {
    persisted.token = await request.post<any, string>('/retrieve', {
      email: retrieveForm.email,
      authcode: retrieveForm.authcode,
      password: CryptoJS.SHA256(retrieveForm.password).toString(CryptoJS.enc.Hex),
    })
    await session.loadUser()
    isDialogOpen.value = false
  } catch {}
}

async function logout() {
  persisted.token = null
  session.user = null
}
</script>

<template>
  <div class="h-screen flex flex-col">

    <el-menu mode="horizontal" :ellipsis="false">
      <el-menu-item>
        企业绿色专利与创新信息系统
      </el-menu-item>
      <el-menu-item class="!ml-auto">
        <el-button v-if="session.user" @click="logout">
          退出登录
        </el-button>
        <el-button v-else type="primary" @click="isDialogOpen=true">
          登录/注册
        </el-button>
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
  <el-dialog v-model="isDialogOpen" title="欢迎使用企业绿色专利系统">
    <el-tabs v-model="currentTab">
      <el-tab-pane label="登录" name="login">
        <el-form :model="loginForm" label-width="auto">
          <el-form-item label="用户名">
            <el-input v-model="loginForm.username" />
          </el-form-item>
          <el-form-item label="密码">
            <el-input v-model="loginForm.password" />
          </el-form-item>
          <el-form-item>
            <el-button class="ms-2" type="primary" @click="login">
              登录
            </el-button>
            <el-button @click="isDialogOpen=false">
              返回
            </el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>
      <el-tab-pane label="注册" name="signup">
        <el-form :model="signupForm" label-width="auto">
          <el-form-item label="邮箱">
            <el-input v-model="signupForm.email" />
          </el-form-item>
          <el-form-item label="邮箱验证码">
            <el-input v-model="signupForm.authcode" />
          </el-form-item>
          <el-form-item label="用户名">
            <el-input v-model="signupForm.username" />
          </el-form-item>
          <el-form-item label="密码">
            <el-input v-model="signupForm.password" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" class="ms-auto" @click="signup">
              注册
            </el-button>
            <el-button @click="isDialogOpen=false">
              返回
            </el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>
      <el-tab-pane label="找回密码" name="retrieve">
        <el-form :model="retrieveForm" label-width="auto">
          <el-form-item label="邮箱">
            <el-input v-model="retrieveForm.email" />
          </el-form-item>
          <el-form-item label="验证码">
            <el-input v-model="retrieveForm.authcode" />
          </el-form-item>
          <el-form-item label="新密码">
            <el-input v-model="retrieveForm.password" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="retrieve" class="ms-auto">
              找回密码
            </el-button>
            <el-button @click="isDialogOpen=false">
              返回
            </el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>
    </el-tabs>
  </el-dialog>
</template>
