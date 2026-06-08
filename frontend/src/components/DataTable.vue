<script setup lang="ts" generic="T extends Record<string, any>">
import { computed, ref, watch } from 'vue'
import {
  type ColumnDef,
  type SortingState,
  FlexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  useVueTable,
} from '@tanstack/vue-table'
import { Search, ArrowUpDown, ArrowUp, ArrowDown, ChevronLeft, ChevronRight } from '@lucide/vue'

export interface Column<R> {
  key: string
  label?: string
  sortable?: boolean
  align?: 'left' | 'right' | 'center'
  width?: string
  cellClass?: string
  headerClass?: string
  hideBelow?: 'sm' | 'md' | 'lg' | 'xl'
  accessor?: (row: R) => unknown
}

const props = withDefaults(
  defineProps<{
    data: T[]
    columns: Column<T>[]
    searchable?: boolean
    searchPlaceholder?: string
    pageSize?: number
    showPagination?: boolean
    rowKey?: keyof T | ((row: T) => string | number)
    rowClass?: (row: T) => string | undefined
  }>(),
  {
    searchable: false,
    searchPlaceholder: 'Search...',
    pageSize: 20,
    showPagination: true,
    rowKey: 'id' as any,
  }
)

const sorting = ref<SortingState>([])
const globalFilter = ref('')

const tanstackColumns = computed<ColumnDef<T>[]>(() =>
  props.columns.map((c) => ({
    id: c.key,
    accessorFn: c.accessor ? c.accessor : (row: T) => (row as any)[c.key],
    header: c.label ?? '',
    enableSorting: c.sortable === true,
  }))
)

const table = useVueTable({
  get data() { return props.data },
  get columns() { return tanstackColumns.value },
  state: {
    get sorting() { return sorting.value },
    get globalFilter() { return globalFilter.value },
  },
  initialState: {
    pagination: { pageSize: props.pageSize, pageIndex: 0 },
  },
  onSortingChange: (updater) => {
    sorting.value = typeof updater === 'function' ? updater(sorting.value) : updater
  },
  onGlobalFilterChange: (updater) => {
    globalFilter.value = typeof updater === 'function' ? updater(globalFilter.value) : updater
  },
  globalFilterFn: 'includesString',
  getCoreRowModel: getCoreRowModel(),
  getSortedRowModel: getSortedRowModel(),
  getFilteredRowModel: getFilteredRowModel(),
  getPaginationRowModel: getPaginationRowModel(),
})

watch(globalFilter, () => table.setPageIndex(0))

const pageCount = computed(() => table.getPageCount())
const pagination = computed(() => table.getState().pagination)

function rowKeyFor(row: T, fallback: number) {
  if (typeof props.rowKey === 'function') return props.rowKey(row)
  const k = props.rowKey as keyof T | undefined
  if (k && (row as any)[k] !== undefined) return (row as any)[k] as string | number
  return fallback
}
</script>

<template>
  <div>
    <div v-if="searchable || $slots.toolbar" class="d-flex align-center ga-2 flex-wrap mb-3">
      <div v-if="searchable" class="position-relative" style="max-width: 320px; width: 100%;">
        <Search
          :size="16"
          :stroke-width="1.5"
          class="position-absolute"
          style="color: #94A3B8; left: 12px; top: 50%; transform: translateY(-50%); z-index: 1;"
        />
        <input
          v-model="globalFilter"
          :placeholder="searchPlaceholder"
          class="w-100"
          style="
            background: rgba(0,0,0,0.5); border: 1px solid rgba(30, 41, 59, 0.8);
            border-radius: 8px; color: white; padding: 8px 12px 8px 36px;
            font-family: 'JetBrains Mono', monospace; font-size: 12px;
            outline: none; transition: border-color 0.2s;
          "
          @focus="$event.target.style.borderColor = '#F7931A'"
          @blur="$event.target.style.borderColor = 'rgba(30, 41, 59, 0.8)'"
        />
      </div>
      <slot name="toolbar" />
    </div>

    <div class="rounded-xl overflow-hidden" style="background: #0F1115; border: 1px solid rgba(30, 41, 59, 0.6);">
      <v-table density="compact">
        <thead>
          <tr>
            <th
              v-for="header in table.getHeaderGroups()[0].headers"
              :key="header.id"
              :style="{
                width: columns.find((c) => c.key === header.column.id)?.width,
                textAlign: columns.find((c) => c.key === header.column.id)?.align || 'left',
              }"
              class="eyebrow"
            >
              <button
                v-if="header.column.getCanSort()"
                class="d-inline-flex align-center ga-1"
                style="background: none; border: none; color: inherit; cursor: pointer; font-family: inherit; font-size: inherit; letter-spacing: inherit; text-transform: inherit;"
                @click="header.column.toggleSorting()"
              >
                <slot :name="`header-${header.column.id}`" :column="header.column">
                  {{ columns.find((c) => c.key === header.column.id)?.label }}
                </slot>
                <ArrowUp v-if="header.column.getIsSorted() === 'asc'" :size="12" :stroke-width="2" style="color: #F7931A;" />
                <ArrowDown v-else-if="header.column.getIsSorted() === 'desc'" :size="12" :stroke-width="2" style="color: #F7931A;" />
                <ArrowUpDown v-else :size="12" :stroke-width="1.5" style="opacity: 0.3;" />
              </button>
              <span v-else>
                <slot :name="`header-${header.column.id}`" :column="header.column">
                  {{ columns.find((c) => c.key === header.column.id)?.label }}
                </slot>
              </span>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="(row, i) in table.getRowModel().rows"
            :key="rowKeyFor(row.original, i)"
            :class="rowClass ? rowClass(row.original) : undefined"
            style="transition: background 0.2s;"
            @mouseenter="$event.currentTarget.style.background = 'rgba(247, 147, 26, 0.03)'"
            @mouseleave="$event.currentTarget.style.background = 'transparent'"
          >
            <td
              v-for="cell in row.getVisibleCells()"
              :key="cell.id"
              :style="{ textAlign: columns.find((c) => c.key === cell.column.id)?.align || 'left' }"
            >
              <slot
                :name="`cell-${cell.column.id}`"
                :row="row.original"
                :value="cell.getValue()"
                :index="i"
              >
                <FlexRender :render="cell.column.columnDef.cell" :props="cell.getContext()" />
              </slot>
            </td>
          </tr>
          <tr v-if="!table.getRowModel().rows.length">
            <td :colspan="columns.length" class="text-center py-8" style="color: #94A3B8;">
              <slot name="empty">
                {{ globalFilter ? 'No results match your search' : 'No data' }}
              </slot>
            </td>
          </tr>
        </tbody>
      </v-table>
    </div>

    <div
      v-if="showPagination && table.getFilteredRowModel().rows.length > pageSize"
      class="d-flex align-center justify-space-between ga-3 mt-3"
    >
      <span class="font-mono text-caption" style="color: #94A3B8;">
        {{ pagination.pageIndex * pagination.pageSize + 1 }}&ndash;{{ Math.min((pagination.pageIndex + 1) * pagination.pageSize, table.getFilteredRowModel().rows.length) }} of {{ table.getFilteredRowModel().rows.length }}
      </span>
      <div class="d-flex align-center ga-2">
        <button
          class="d-inline-flex align-center ga-1 rounded-pill px-3"
          style="height: 32px; background: rgba(247, 147, 26, 0.1); border: 1px solid rgba(247, 147, 26, 0.2); color: #F7931A; font-family: 'JetBrains Mono', monospace; font-size: 12px; transition: all 0.2s;"
          :disabled="!table.getCanPreviousPage()"
          :style="{ opacity: !table.getCanPreviousPage() ? 0.4 : 1, cursor: !table.getCanPreviousPage() ? 'not-allowed' : 'pointer' }"
          @click="table.previousPage()"
        >
          <ChevronLeft :size="14" :stroke-width="1.5" />
          Prev
        </button>
        <span class="font-mono text-caption" style="color: #94A3B8;">
          {{ pagination.pageIndex + 1 }} / {{ pageCount }}
        </span>
        <button
          class="d-inline-flex align-center ga-1 rounded-pill px-3"
          style="height: 32px; background: rgba(247, 147, 26, 0.1); border: 1px solid rgba(247, 147, 26, 0.2); color: #F7931A; font-family: 'JetBrains Mono', monospace; font-size: 12px; transition: all 0.2s;"
          :disabled="!table.getCanNextPage()"
          :style="{ opacity: !table.getCanNextPage() ? 0.4 : 1, cursor: !table.getCanNextPage() ? 'not-allowed' : 'pointer' }"
          @click="table.nextPage()"
        >
          Next
          <ChevronRight :size="14" :stroke-width="1.5" />
        </button>
      </div>
    </div>
  </div>
</template>
