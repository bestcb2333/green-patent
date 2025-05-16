<script setup lang="ts">
import {request} from '@/axios';
import usePersistedStore from '@/stores/persisted';
import type {Keyword, Report} from '@/tables';
import {formatDate} from '@/utils';
import {useRouteQuery} from '@vueuse/router';
import type {EChartsOption} from 'echarts';
import {LineChart} from 'echarts/charts';
import {GridComponent, TitleComponent, TooltipComponent} from 'echarts/components';
import {use} from 'echarts/core';
import {CanvasRenderer} from 'echarts/renderers';
import {ElNotification, type TableInstance} from 'element-plus';
import {computed, reactive, ref, watch} from 'vue';
import VChart from 'vue-echarts'
import {useRoute} from 'vue-router';

use([
  TitleComponent,
  TooltipComponent,
  LineChart,
  GridComponent,
  CanvasRenderer,
])

const route = useRoute()
const persisted = usePersistedStore()
const tableRef = ref<TableInstance>()

const page = useRouteQuery('page', 1, {transform: Number})
const pageSize = useRouteQuery('page_size', 100, {transform: Number})
const status = useRouteQuery<any, boolean>('status', false, {transform: Boolean})
const total = ref(0)
const reports = ref<Report[]>([])
watch(([page, pageSize, status]), loadTable, {immediate: true})

async function loadTable() {
  try {
    const res = await request.get<any, {
      total: number,
      data: Report[],
    }>('/reports', {params: {
      page: page.value,
      page_size: pageSize.value,
      status: status.value,
    }})
    total.value = res.total
    reports.value = res.data
  } catch {}
}

const trend = ref<number[]>([])
request.get<any, number[]>('/trend').then(res => {
  trend.value = res
}).catch(() => {})
const chartOption = computed<EChartsOption>(() => ({
  title: {
    text: '不同年的年报数量',
  },
  tooltip: {
    trigger: 'axis',
  },
  grid: {
    top: 30,
    bottom: 20,
  },
  xAxis: {
    type: 'category',
    data: ['2021', '2022', '2023', '2024', '2025']
  },
  yAxis: {
    type: 'value'
  },
  series: [
    {
      data: trend.value,
      type: 'line'
    }
  ]
}))

interface Stats {
  pending: number,
  solved: number,
  total: number,
}

const stats = ref<Stats|null>(null)
request.get<any, Stats>('/stats/reports').then(res => {
  stats.value = res
}).catch(() => {})

const isDialogOpen = ref(false)

const addReportForm = reactive({
  name: '',
  year: 2025,
})

const currentRow = ref<Report|null>(null)

watch(() => route.params.id as string, async id => {
  if (!id) return
  try {
    const res = await request.get<any, Report>(`/reports/${id}`)
    currentRow.value = res
  } catch {}
}, {immediate: true})

async function getKeywords() {
  if (!currentRow.value) return
  try {
    ElNotification({
      title: '已收到请求，请勿重复点击',
    })
    const res = await request.get<any, Keyword[]>(`/extract/report/${currentRow.value.id}`, {
      timeout: 60000,
    })
    const keywords = res.map(item => item.value)
    ElNotification({
      title: '关键词生成成功请刷新页面',
      message: keywords.toString(),
    })
  } catch {}
}

async function getSuggestion() {
   if (!currentRow.value) return
  try {
    ElNotification({
      title: '已收到请求，请勿重复点击'
    })
    const res = await request.get<any, string>(`/generate/report/${currentRow.value.id}`, {
      timeout: 60000,
    })
    ElNotification({
      title: '建议生成成功，请查看',
      message: res,
    })
  } catch {}
}
</script>

<template>
  <div class="h-full p-2 grid grid-cols-2 grid-rows-[1fr_2fr] gap-2">

    <el-card shadow="hover" body-class="h-full">
      <v-chart :option="chartOption" autoresize />
    </el-card>

    <el-card class="row-span-2 flex flex-col" shadow="hover"
      body-class="grow flex flex-col gap-2"
      header-class="flex justify-between"
    >

      <template v-if="currentRow" #header>
        <div>
          {{currentRow.name}}
        </div>
        <div class="flex gap-2">
          <el-button type="warning" @click="getKeywords">
            生成关键词
          </el-button>
          <a :href="`${persisted.setting.reportAddr}${currentRow.file}`" target="_blank">
            <el-button type="success">
              查看正文
            </el-button>
          </a>
          <el-button @click="$router.push(`/suggestions?report_id=${currentRow.id}`)" type="primary">
            查看创新建议
          </el-button>
          <el-button type="success" @click="getSuggestion">
            生成创新建议
          </el-button>
        </div>
      </template>

      <template v-if="currentRow">
        <div v-if="currentRow.keywords.length" class="flex flex-wrap gap-2">
          <el-tag v-for="keyword in currentRow.keywords" :key="keyword.id">
            {{keyword.value}}
          </el-tag>
        </div>
        <template v-else>
          该年报尚未生成关键词
        </template>
        <el-card shadow="hover" class="grow" body-class="h-full p-0">
          <iframe :src="`${persisted.setting.reportAddr}${currentRow.file}`" frameborder="0"
            width="100%" height="100%" allowfullscreen
          >
          </iframe>
        </el-card>
      </template>

      <el-empty v-else description="请选择一项年报以查看" class="mx-auto" />

    </el-card>

    <el-card shadow="hover" class="flex flex-col"
      header-class="flex justify-between"
      body-class="grow min-h-0 overflow-y-auto"
    >

      <template #header>
        <div class="space-x-2">
          <span>年报列表</span>
          <el-tag type="primary">总数：{{stats?.total}}</el-tag>
          <el-tag type="success">已处理：{{stats?.solved}}</el-tag>
          <el-tag type="warning">待处理：{{stats?.pending}}</el-tag>
        </div>
        <div>
          <el-switch v-model="status" active-text="仅未处理" inactive-text="全部年报" />
          <el-button class="ms-2" type="primary" @click="isDialogOpen=true">
            上传年报
          </el-button>
        </div>
      </template>

      <el-table :data="reports" highlight-current-row ref="tableRef"
        @current-change="val=>$router.push(`/reports/${val.id}`)"
      >
        <el-table-column label="创建时间" prop="createdAt" :formatter="formatDate" />
        <el-table-column label="名称" prop="name" />
        <el-table-column label="年份" prop="year" />
        <el-table-column label="状态">
          <template #default="{row}">
            <el-tag :type="row.status?'success':'danger'">
              {{row.status?'已处理':'处理中'}}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="上传者" prop="user.nickname" />
        <el-table-column label="操作">
          <template #default="{row}">
            <a :href="`${persisted.setting.reportAddr}${row.file}`" download>
              <el-button>
                下载
              </el-button>
            </a>
          </template>
        </el-table-column>
      </el-table>

      <template #footer>
        <el-pagination layout="sizes, prev, pager, next, total" :total="total"
          v-model:current-page="page" v-model:page-size="pageSize"
        />
      </template>

    </el-card>

  </div>
  <el-dialog v-model="isDialogOpen" title="上传年报">
    <el-form :model="addReportForm">
      <el-form-item label="名称">
        <el-input v-model="addReportForm.name" />
      </el-form-item>
      <el-form-item label="年份">
        <el-input-number v-model="addReportForm.year" />
      </el-form-item>
      <el-form-item>
        请先填写名称和年份，在上传专利，只支持PDF格式
      </el-form-item>
      <el-form-item label="文件">
        <el-upload
          :action="`${persisted.setting.apiAddr}reports`"
          :data="addReportForm"
        >
          <el-button type="primary">
            点击上传文件
          </el-button>
        </el-upload>
      </el-form-item>
      <el-form-item>
        <el-button @click="isDialogOpen=false">
          取消
        </el-button>
      </el-form-item>
    </el-form>
  </el-dialog>
</template>
