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

const emit = defineEmits<{
  'update:search': [value: string]
}>()

const sorting = ref<SortingState>([])
const globalFilter = ref('')
const searchInput = ref('')

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
  sm: 'd-sm-none d-md-table-cell',
  md: 'd-none d-md-table-cell',
  lg: 'd-none d-lg-table-cell',
  xl: 'd-none d-xl-table-cell',
}

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
      <v-text-field
        v-if="searchable"
        v-model="globalFilter"
        :placeholder="searchPlaceholder"
        prepend-inner-icon="mdi-magnify"
        clearable
        hide-details
        density="compact"
        variant="outlined"
        style="max-width: 320px;"
      />
      <slot name="toolbar" />
    </div>

    <v-card>
      <v-table density="compact" class="table-tight">
        <thead>
          <tr>
            <th
              v-for="header in table.getHeaderGroups()[0].headers"
              :key="header.id"
              :class="columns.find((c) => c.key === header.column.id)?.hideBelow ? hideBelowMap[columns.find((c) => c.key === header.column.id)!.hideBelow!] : ''"
              :style="{ width: columns.find((c) => c.key === header.column.id)?.width, textAlign: columns.find((c) => c.key === header.column.id)?.align || 'left' }"
              class="text-caption font-weight-bold text-uppercase text-medium-emphasis"
            >
              <v-btn
                v-if="header.column.getCanSort()"
                variant="text"
                density="compact"
                size="small"
                :class="columns.find((c) => c.key === header.column.id)?.align === 'right' ? 'flex-row-reverse' : ''"
                @click="header.column.toggleSorting()"
              >
                <slot :name="`header-${header.column.id}`" :column="header.column">
                  {{ columns.find((c) => c.key === header.column.id)?.label }}
                </slot>
                <v-icon
                  v-if="header.column.getIsSorted() === 'asc'"
                  size="x-small"
                  class="ml-1"
                >mdi-arrow-up</v-icon>
                <v-icon
                  v-else-if="header.column.getIsSorted() === 'desc'"
                  size="x-small"
                  class="ml-1"
                >mdi-arrow-down</v-icon>
                <v-icon
                  v-else
                  size="x-small"
                  class="ml-1 text-disabled"
                >mdi-arrow-up-down</v-icon>
              </v-btn>
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
              :class="columns.find((c) => c.key === cell.column.id)?.hideBelow ? hideBelowMap[columns.find((c) => c.key === cell.column.id)!.hideBelow!] : ''"
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
            <td :colspan="columns.length" class="text-center text-medium-emphasis py-8">
              <slot name="empty">
                {{ globalFilter ? 'No results match your search' : 'No data' }}
              </slot>
            </td>
          </tr>
        </tbody>
      </v-table>
    </v-card>

    <div
      v-if="showPagination && table.getFilteredRowModel().rows.length > pageSize"
      class="d-flex align-center justify-space-between ga-3 mt-3"
    >
      <span class="text-caption text-medium-emphasis tabular-nums">
        {{ pagination.pageIndex * pagination.pageSize + 1 }}&ndash;{{ Math.min((pagination.pageIndex + 1) * pagination.pageSize, table.getFilteredRowModel().rows.length) }} of {{ table.getFilteredRowModel().rows.length }}
      </span>
      <div class="d-flex align-center ga-2">
        <v-btn
          size="small"
          variant="tonal"
          :disabled="!table.getCanPreviousPage()"
          @click="table.previousPage()"
        >
          <v-icon size="small">mdi-chevron-left</v-icon>
          Previous
        </v-btn>
        <span class="text-caption text-medium-emphasis tabular-nums">
          Page {{ pagination.pageIndex + 1 }} / {{ pageCount }}
        </span>
        <v-btn
          size="small"
          variant="tonal"
          :disabled="!table.getCanNextPage()"
          @click="table.nextPage()"
        >
          Next
          <v-icon size="small">mdi-chevron-right</v-icon>
        </v-btn>
      </div>
    </div>
  </div>
</template>

<style scoped>
.table-tight th {
  height: 36px;
  white-space: nowrap;
}
.table-tight td {
  height: 36px;
}
.table-tight tbody tr:hover {
  background: #fbfaf6;
}
</style>
