<script setup lang="ts">
import {request} from '@/axios';
import type {User} from '@/tables';
import {formatDate} from '@/utils';
import dayjs from 'dayjs';
import {ref, watch} from 'vue';

const users = ref<User[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
watch(([page, pageSize]), async ([page, pageSize]) => {
  try {
    const res = await request.get<any, {
      total: number,
      data: User[],
    }>(`/users?page=${page}&page_size=${pageSize}`)
    total.value = res.total
    users.value = res.data
  } catch {}
}, {immediate: true})

const currentRow = ref<User|null>(null)
async function handleCurrentChange(user: User|null) {
  if (!user) return
  try {
    const res = await request.get<any, User>(`/users/${user.id}`)
    currentRow.value = res
  } catch {}
}

const activeCollapses = ref([])
</script>

<template>
  <div class="h-full flex gap-2 p-2">

    <el-card shadow="hover" class="basis-0 grow flex flex-col"
      body-class="grow overflow-y-auto" header-class="flex justify-between"
    >

      <template #header>
        <div class="text-lg">
          用户列表
        </div>
        <el-button type="primary">
          添加用户
        </el-button>
      </template>

      <el-table :data="users" highlight-current-row @current-change="handleCurrentChange">
        <el-table-column label="用户名" prop="name" />
        <el-table-column label="昵称" prop="nickname" />
        <el-table-column label="邮箱" prop="email" />
        <el-table-column label="手机号" prop="phone" />
        <el-table-column label="注册时间" prop="createdAt" :formatter="formatDate" />
        <el-table-column label="操作" />
      </el-table>

      <template #footer>
        <el-pagination layout="sizes, prev, pager, next, total" :total="total"
          v-model:current-page="page" v-model:page-size="pageSize"
        />
      </template>

    </el-card>

    <el-card shadow="hover" class="basis-0 grow">

      <div v-if="currentRow">

        <el-descriptions title="用户信息">
          <el-descriptions-item label="用户名">
            {{currentRow.name}}
          </el-descriptions-item>
          <el-descriptions-item label="昵称">
            {{currentRow.nickname}}
          </el-descriptions-item>
          <el-descriptions-item label="邮箱">
            {{currentRow.email}}
          </el-descriptions-item>
          <el-descriptions-item label="手机号">
            {{currentRow.phone}}
          </el-descriptions-item>
          <el-descriptions-item label="注册时间">
            {{dayjs(currentRow.createdAt).format('MM月DD日 HH:mm')}}
          </el-descriptions-item>
          <el-descriptions-item label="是否管理员">
            <el-tag :type="currentRow.admin?'success':'danger'">
              {{currentRow.admin?'是':'否'}}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>


        <el-collapse v-model="activeCollapses">
          <el-collapse-item title="上传的专利" name="patent">
            <el-button v-for="patent in currentRow.patents" :key="patent.id" class="mt-2"
              @click="$router.push(`/patents/${patent.id}`)"
            >
              {{patent.name}}
            </el-button>
          </el-collapse-item>
          <el-collapse-item title="上传的年报" name="report">
            <el-button v-for="report in currentRow.reports" :key="report.id" class="mt-2"
              @click="$router.push(`/reports/${report.id}`)"
            >
              {{report.name}}
            </el-button>
          </el-collapse-item>
          <el-collapse-item title="生成的建议" name="suggestion">
            <el-button v-for="suggestion in currentRow.suggestions" :key="suggestion.id" class="mt-2"
              @click="$router.push(`/suggestions/${suggestion.id}`)"
            >
              {{suggestion.title}}
            </el-button>
          </el-collapse-item>
          <el-collapse-item title="添加的关键词" name="keyword">
            <el-button v-for="keyword in currentRow.keywords" :key="keyword.id" class="mt-2"
              @click="$router.push(`/keywords/${keyword.id}`)"
            >
              {{keyword.value}}
            </el-button>
          </el-collapse-item>
        </el-collapse>
      </div>

      <el-empty v-else description="请选择一个用户查看具体信息" />

    </el-card>

  </div>
</template>
