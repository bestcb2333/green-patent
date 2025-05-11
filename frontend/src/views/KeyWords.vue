<script setup lang="ts">
import {request} from '@/axios';
import type {Keyword} from '@/tables';
import {formatDate} from '@/utils';
import {reactive, ref, watch} from 'vue';

const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const currentRow = ref<Keyword|null>(null)
const keywords = ref<Keyword[]>([])
watch(([page, pageSize]), async ([page, pageSize]) => {
  try {
    const res = await request.get<any, {
      total: number,
      data: Keyword[],
    }>('/keywords', {params: {
      page: page,
      page_size: pageSize,
    }})
    total.value = res.total
    keywords.value = res.data
  } catch {}
}, {immediate: true})

async function handleSelectionChange(row: Keyword|null) {
  if (!row) return
  try {
    const res = await request.get<any, Keyword>(`/keywords/${row.id}`)
    currentRow.value = res
  } catch {}
}

const isDialogOpen = ref(false)

const addKeywordForm = reactive({
  value: '',
})
</script>

<template>
  <div class="h-full p-2 flex gap-2">

    <el-card shadow="hover" class="basis-0 grow flex flex-col"
      header-class="flex justify-between" body-class="grow min-h-0 overflow-y-auto"
    >

      <template #header>
        <div>
          关键词列表
        </div>
        <el-button type="primary" @click="isDialogOpen=true">
          添加关键词
        </el-button>
      </template>

      <el-table :data="keywords" highlight-current-row @current-change="handleSelectionChange">
        <el-table-column label="创建时间" prop="createdAt" :formatter="formatDate" />
        <el-table-column label="创建者" prop="user.nickname" />
        <el-table-column label="关键词" prop="value" />
      </el-table>

      <template #footer>
        <el-pagination layout="sizes, prev, pager, next, total" :total="total"
          v-model:current-page="page" v-model:page-size="pageSize"
        />
      </template>

    </el-card>

    <el-card shadow="hover" class="basis-0 grow">
      <template v-if="currentRow">
        <div class="font-bold text-xl">
          {{currentRow.value}}
        </div>
        <el-divider />
        <div class="font-bold text-lg">
          相关的绿色专利
        </div>
        <el-button v-for="patent in currentRow.patents" :key="patent.id" type="success"
          @click="$router.push(`/patents?id=${patent.id}`)"
        >
          {{patent.name}}（{{patent.number}}）
        </el-button>
        <el-divider />
        <div class="font-bold text-lg">
          相关的年度报告
        </div>
        <el-button v-for="report in currentRow.reports" :key="report.id" type="primary"
          @click="$router.push(`/reports?id=${report.id}`)"
        >
          {{report.name}}（{{report.year}}）
        </el-button>
      </template>
    </el-card>

  </div>
  <el-dialog v-model="isDialogOpen" title="添加关键词">
    <el-form :model="addKeywordForm">
      <el-form-item label="名称">
        <el-input v-model="addKeywordForm.value" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" class="ms-auto">
          确认添加
        </el-button>
        <el-button @click="isDialogOpen=false">
          取消
        </el-button>
      </el-form-item>
    </el-form>
  </el-dialog>
</template>
