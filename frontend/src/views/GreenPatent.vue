<script setup lang="ts">
import {request} from '@/axios';
import usePersistedStore from '@/stores/persisted';
import type {Keyword, Patent} from '@/tables';
import {formatDate} from '@/utils';
import {useRouteQuery} from '@vueuse/router';
import {ElNotification, type TableInstance} from 'element-plus';
import {reactive, ref, watch} from 'vue';
import {useRoute} from 'vue-router';

const tableRef = ref<TableInstance>()
const persisted = usePersistedStore()
const route = useRoute()
const page = useRouteQuery('page', 1, {transform: Number})
const pageSize = useRouteQuery('page_size', 100, {transform: Number})
const status = useRouteQuery<any, boolean>('status', false, {transform: Boolean})
const total = ref(0)
const patents = ref<Patent[]>([])
watch(([page, pageSize, status]), loadTable, {immediate: true})

async function loadTable() {
  tableRef.value?.setCurrentRow(undefined)
  try {
    const res = await request.get<any, {
      total: number,
      data: Patent[],
    }>('/patents', {params: {
      page: page.value,
      page_size: pageSize.value,
      status: status.value,
    }})
    total.value = res.total
    patents.value = res.data
  } catch {}
}

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
  number: '',
})

const currentRow = ref<Patent|null>(null)

watch(() => route.params.id as string, async id => {
  if (!id) return
  try {
    const res = await request.get<any, Patent>(`/patents/${id}`)
    currentRow.value = res
  } catch {}
}, {immediate: true})

async function getKeywords() {
  if (!currentRow.value) return
  try {
    ElNotification({
      title: '已收到请求，请勿重复点击',
    })
    const res = await request.get<any, Keyword[]>(`/extract/patent/${currentRow.value.id}`, {
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
    const res = await request.get<any, string>(`/generate/patent/${currentRow.value.id}`, {
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
        <div>
          <el-switch v-model="status" active-text="仅未处理" inactive-text="全部专利" />
          <el-button class="ms-2" type="primary" @click="isDialogOpen=true">
            上传专利
          </el-button>
        </div>
      </template>

      <el-table :data="patents" highlight-current-row
        @current-change="val=>$router.push({path: `/patents/${val.id}`, query:route.query})" ref="tableRef">
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
          <template #default="{row}">
            <a :href="`${persisted.setting.patentAddr}${row.file}`" download>
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

    <el-card class="basis-0 grow row-span-2 flex flex-col" shadow="hover"
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
          <a :href="`${persisted.setting.patentAddr}${currentRow.file}`" target="_blank">
            <el-button type="success">
              查看正文
            </el-button>
          </a>
          <el-button @click="$router.push(`/suggestions?patent_id=${currentRow.id}`)" type="primary">
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
          该专利尚未生成关键词
        </template>
        <el-card shadow="hover" body-class="h-full p-0" class="grow">
          <iframe :src="`${persisted.setting.patentAddr}${currentRow.file}`" frameborder="0"
            class="h-full w-full"
          >
          </iframe>
        </el-card>
      </template>

      <el-empty v-else description="请选择一项专利以查看" class="mx-auto" />

    </el-card>

  </div>
  <el-dialog v-model="isDialogOpen" title="上传专利">
    <el-form :model="addPatentForm">
      <el-form-item label="名称">
        <el-input v-model="addPatentForm.name" />
      </el-form-item>
      <el-form-item label="专利号">
        <el-input v-model="addPatentForm.number" />
      </el-form-item>
      <el-form-item>
        请先填写名称和专利号，再上传文件，只支持PDF格式
      </el-form-item>
      <el-form-item label="文件">
        <el-upload
          :action="`${persisted.setting.apiAddr}patents`"
          :data="addPatentForm"
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
