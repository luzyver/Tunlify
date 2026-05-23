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
import { ArrowDown, ArrowUp, ArrowUpDown, ChevronLeft, ChevronRight, Search } from 'lucide-vue-next'

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
    searchPlaceholder: 'Search…',
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

const hideBelowMap: Record<string, string> = {
  sm: 'hidden sm:table-cell',
  md: 'hidden md:table-cell',
  lg: 'hidden lg:table-cell',
  xl: 'hidden xl:table-cell',
}
const alignMap: Record<string, string> = {
  left: 'text-left',
  right: 'text-right',
  center: 'text-center',
}

function colClass(c: Column<T>) {
  const parts: string[] = []
  if (c.hideBelow) parts.push(hideBelowMap[c.hideBelow])
  if (c.align && c.align !== 'left') parts.push(alignMap[c.align])
  if (c.cellClass) parts.push(c.cellClass)
  return parts.join(' ')
}

function headerClass(c: Column<T>) {
  const parts: string[] = []
  if (c.hideBelow) parts.push(hideBelowMap[c.hideBelow])
  if (c.align && c.align !== 'left') parts.push(alignMap[c.align])
  if (c.headerClass) parts.push(c.headerClass)
  return parts.join(' ')
}

function rowKeyFor(row: T, fallback: number) {
  if (typeof props.rowKey === 'function') return props.rowKey(row)
  const k = props.rowKey as keyof T | undefined
  if (k && (row as any)[k] !== undefined) return (row as any)[k] as string | number
  return fallback
}

const rangeLabel = computed(() => {
  const filtered = table.getFilteredRowModel().rows.length
  if (!filtered) return '0'
  const { pageIndex, pageSize } = table.getState().pagination
  const from = pageIndex * pageSize + 1
  const to = Math.min(from + pageSize - 1, filtered)
  return `${from}–${to} of ${filtered}`
})
</script>

<template>
  <div class="space-y-3">
    <div v-if="searchable || $slots.toolbar" class="flex items-center gap-2 flex-wrap">
      <div v-if="searchable" class="relative">
        <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-text-dim" :stroke-width="1.75" />
        <input
          v-model="globalFilter"
          type="text"
          :placeholder="searchPlaceholder"
          class="input !pl-9 !py-1.5 w-[280px]"
        />
      </div>
      <slot name="toolbar" />
    </div>

    <div class="card overflow-hidden">
      <table class="table-tight">
        <thead>
          <tr>
            <th
              v-for="header in table.getHeaderGroups()[0].headers"
              :key="header.id"
              :class="[colClass(columns.find((c) => c.key === header.column.id)!), columns.find((c) => c.key === header.column.id)?.width && '']"
              :style="{ width: columns.find((c) => c.key === header.column.id)?.width }"
            >
              <button
                v-if="header.column.getCanSort()"
                type="button"
                class="inline-flex items-center gap-1.5 hover:text-text transition-colors"
                @click="header.column.toggleSorting()"
              >
                <slot :name="`header-${header.column.id}`" :column="header.column">
                  {{ columns.find((c) => c.key === header.column.id)?.label }}
                </slot>
                <ArrowUp v-if="header.column.getIsSorted() === 'asc'" class="w-3 h-3" :stroke-width="2" />
                <ArrowDown v-else-if="header.column.getIsSorted() === 'desc'" class="w-3 h-3" :stroke-width="2" />
                <ArrowUpDown v-else class="w-3 h-3 opacity-40" :stroke-width="2" />
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
          >
            <td
              v-for="cell in row.getVisibleCells()"
              :key="cell.id"
              :class="colClass(columns.find((c) => c.key === cell.column.id)!)"
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
            <td :colspan="columns.length" class="text-center text-text-muted py-8 text-sm">
              <slot name="empty">
                {{ globalFilter ? 'No results match your search' : 'No data' }}
              </slot>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div
      v-if="showPagination && table.getFilteredRowModel().rows.length > pageSize"
      class="flex items-center justify-between gap-3"
    >
      <span class="text-2xs text-text-dim tabular-nums">{{ rangeLabel }}</span>
      <div class="flex items-center gap-2">
        <button
          type="button"
          class="btn-secondary !py-1 !px-2.5 !text-xs"
          :disabled="!table.getCanPreviousPage()"
          @click="table.previousPage()"
        >
          <ChevronLeft class="w-3.5 h-3.5" :stroke-width="1.75" />
          Previous
        </button>
        <span class="text-2xs text-text-dim tabular-nums">
          Page {{ table.getState().pagination.pageIndex + 1 }} / {{ table.getPageCount() }}
        </span>
        <button
          type="button"
          class="btn-secondary !py-1 !px-2.5 !text-xs"
          :disabled="!table.getCanNextPage()"
          @click="table.nextPage()"
        >
          Next
          <ChevronRight class="w-3.5 h-3.5" :stroke-width="1.75" />
        </button>
      </div>
    </div>
  </div>
</template>
