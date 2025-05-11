<script setup lang="ts">
import {request} from '@/axios';
import type {Suggestion} from '@/tables';
import {formatDate} from '@/utils';
import {useUrlSearchParams} from '@vueuse/core';
import {ref, watch} from 'vue';
import {useRoute} from 'vue-router';

const route = useRoute()
const params = useUrlSearchParams('history')

const currentRow = ref<Suggestion|null>(null)
const suggestions = ref<Suggestion[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
watch(([page, pageSize, params]), ([page, pageSize, params]) => {
  try {
    request.get<any, {
      total: number,
      data: Suggestion[],
    }>('/suggestions', {params: {
      page: page,
      page_size: pageSize,
      ...params,
    }}).then(res => {
      total.value = res.total
      suggestions.value = res.data
    }).catch(() => {})
  } catch {}
}, {immediate: true})

watch(() => route.params.id as string, async id => {
  if (!id) return
  try {
    const res = await request.get<any, Suggestion>(`/suggestions/${id}`)
    currentRow.value = res
  } catch {}
}, {immediate: true})
</script>

<template>
  <div class="h-full p-2 flex gap-2">

    <el-card shadow="hover" class="basis-0 grow flex flex-col"
      body-class="grow min-h-0 overflow-y-auto"
    >

      <template #header>
        创新建议列表
      </template>

      <el-table :data="suggestions" highlight-current-row
        @current-change="val=>$router.push(`/suggestions/${val.id}`)"
      >
        <el-table-column label="创建时间" prop="createdAt" :formatter="formatDate" />
        <el-table-column label="来源(专利/年报)" width="200">
          <template #default="{row}">
            <el-button :type="row.patent?'success':'primary'"
              @click="$router.push(row.patent?`/patents/${row.patent.id}`:`/reports/${row.report.id}`)"
            >
              {{row.patent?row.patent.name:row.report.name}}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column label="创建者" prop="user.nickname" />
        <el-table-column label="标题" prop="title" />
      </el-table>

      <template #footer>
        <el-pagination layout="sizes, prev, pager, next, total" :total="total"
          v-model:current-page="page" v-model:page-size="pageSize"
        />
      </template>

    </el-card>

    <el-card shadow="hover" class="basis-0 grow">

      <template v-if="currentRow">
        <div class="font-bold text-lg">
          {{currentRow.title}}
        </div>
        <div>
          创建者：{{currentRow.user?.nickname}}
        </div>
        <div>
          {{currentRow.content}}
        </div>
      </template>

      <el-empty v-else class="mx-auto" description="请选择一项查看" />

    </el-card>

  </div>
</template>
