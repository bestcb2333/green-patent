<script setup lang="ts">
import {request} from '@/axios';
import type {Report} from '@/tables';
import {formatDate} from '@/utils';
import {useUrlSearchParams} from '@vueuse/core';
import type {EChartsOption} from 'echarts';
import {LineChart} from 'echarts/charts';
import {GridComponent, TitleComponent, TooltipComponent} from 'echarts/components';
import {use} from 'echarts/core';
import {CanvasRenderer} from 'echarts/renderers';
import {computed, reactive, ref, watch} from 'vue';
import VChart from 'vue-echarts'
import {useRoute} from 'vue-router';

const route = useRoute()
const params = useUrlSearchParams('history')

use([
  TitleComponent,
  TooltipComponent,
  LineChart,
  GridComponent,
  CanvasRenderer,
])

const page = computed({
  get: () => parseInt(params.page as string)||1,
  set: (page: number) => params.page = page.toString(),
})

const pageSize = computed({
  get: () => parseInt(params.page_size as string)||10,
  set: (pageSize: number) => params.page_size = pageSize.toString(),
})

const total = ref(0)
const reports = ref<Report[]>([])
watch(([page, pageSize]), async ([page, pageSize]) => {
  try {
    request.get<any, {
      total: number,
      data: Report[],
    }>(`/reports?page=${page}&page_size=${pageSize}`).then(res => {
      total.value = res.total
      reports.value = res.data
    }).catch(() => {})
  } catch {}
}, {immediate: true})

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
</script>

<template>
  <div class="h-full p-2 grid grid-cols-2 grid-rows-[1fr_2fr] gap-2">

    <el-card shadow="hover" body-class="h-full">
      <v-chart :option="chartOption" autoresize />
    </el-card>

    <el-card class="row-span-2" shadow="hover"
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
          该年报尚未生成关键词
        </template>
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
        <el-button type="primary" @click="isDialogOpen=true">
          上传年报
        </el-button>
      </template>

      <el-table :data="reports" highlight-current-row
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

  </div>
  <el-dialog v-model="isDialogOpen" title="上传年报">
    <el-form :model="addReportForm">
      <el-form-item label="名称">
        <el-input v-model="addReportForm.name" />
      </el-form-item>
      <el-form-item label="年份">
        <el-input-number v-model="addReportForm.year" />
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
