<script setup lang="ts">
import {request} from '@/axios';
import type {Patent} from '@/tables';
import {formatDate} from '@/utils';
import {useUrlSearchParams} from '@vueuse/core';
import {computed, reactive, ref, watch} from 'vue';
import {useRoute} from 'vue-router';

const params = useUrlSearchParams('history')
const route = useRoute()

const page = computed({
  get: () => parseInt(params.page as string)||1,
  set: (page: number) => params.page = page.toString(),
})

const pageSize = computed({
  get: () => parseInt(params.page_size as string)||10,
  set: (pageSize: number) => params.page_size = pageSize.toString(),
})

const total = ref(0)
const patents = ref<Patent[]>([])
watch(([page, pageSize]), async ([page, pageSize]) => {
  try {
    request.get<any, {
      total: number,
      data: Patent[],
    }>('/patents', {params: {
      page: page,
      page_size: pageSize,
    }}).then(res => {
      total.value = res.total
      patents.value = res.data
    }).catch(() => {})
  } catch {}
}, {immediate: true})

interface Stats {
  pending: number,
  solved: number,
  total: number,
}

const stats = ref<Stats|null>(null)
request.get<any, Stats>('/stats/patents').then(res => {
  stats.value = res
}).catch(() => {})

const isDialogOpen = ref(false)

const addPatentForm = reactive({
  name: '',
  num: '',
})

const currentRow = ref<Patent|null>(null)

watch(() => route.params.id as string, async id => {
  if (!id) return
  try {
    const res = await request.get<any, Patent>(`/patents/${id}`)
    currentRow.value = res
  } catch {}
}, {immediate: true})
</script>

<template>
  <div class="h-full p-2 flex gap-2">

    <el-card shadow="hover" class="basis-0 grow flex flex-col"
      body-class="grow min-h-0 overflow-y-auto" header-class="flex justify-between"
    >

      <template #header>
        <div class="space-x-2">
          <span>绿色专利列表</span>
          <el-tag type="primary">总数：{{stats?.total}}</el-tag>
          <el-tag type="success">已处理：{{stats?.solved}}</el-tag>
          <el-tag type="warning">待处理：{{stats?.pending}}</el-tag>
        </div>
        <el-button type="primary" @click="isDialogOpen=true">
          上传年报
        </el-button>
      </template>

      <el-table :data="patents" highlight-current-row @current-change="val=>$router.push(`/patents/${val.id}`)">
        <el-table-column label="创建时间" prop="createdAt" :formatter="formatDate" />
        <el-table-column label="名称" prop="name" />
        <el-table-column label="专利号" prop="number" />
        <el-table-column label="状态">
          <template #default="{row}">
            <el-tag :type="row.status?'success':'danger'">
              {{row.status?'已处理':'处理中'}}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="上传者" prop="user.nickname" />
        <el-table-column label="操作">
          <el-button>
            下载
          </el-button>
        </el-table-column>
      </el-table>

      <template #footer>
        <el-pagination layout="sizes, prev, pager, next, total" :total="total"
          v-model:current-page="page" v-model:page-size="pageSize"
        />
      </template>

    </el-card>

    <el-card class="basis-0 grow row-span-2" shadow="hover"
      body-class="flex flex-wrap gap-2"
      header-class="flex justify-between"
    >

      <template v-if="currentRow" #header>
        <div>
          {{currentRow.name}}
        </div>
        <div>
          <el-button type="warning">
            生成关键词
          </el-button>
          <el-button type="success">
            查看正文
          </el-button>
          <el-button @click="$router.push(`/suggestions/${currentRow.id}`)" type="primary">
            查看创新建议
          </el-button>
        </div>
      </template>

      <template v-if="currentRow">
        <template v-if="currentRow.keywords.length">
          <el-tag v-for="keyword in currentRow.keywords" :key="keyword.id">
            {{keyword.value}}
          </el-tag>
        </template>
        <template v-else>
          该专利尚未生成关键词
        </template>
      </template>

      <el-empty v-else description="请选择一项专利以查看" class="mx-auto" />

    </el-card>

  </div>
  <el-dialog v-model="isDialogOpen" title="上传年报">
    <el-form :model="addPatentForm">
      <el-form-item label="名称">
        <el-input v-model="addPatentForm.name" />
      </el-form-item>
      <el-form-item label="专利号">
        <el-input v-model="addPatentForm.num" />
      </el-form-item>
      <el-form-item label="文件">

      </el-form-item>
      <el-form-item>
        <el-button type="primary" class="ms-auto">
          确认上传
        </el-button>
        <el-button @click="isDialogOpen=false">
          取消
        </el-button>
      </el-form-item>
    </el-form>
  </el-dialog>
</template>
