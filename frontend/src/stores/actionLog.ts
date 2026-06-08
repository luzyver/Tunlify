import { defineStore } from 'pinia'

export interface ActionEntry {
  id: number
  label: string
  projectName?: string
  status: 'running' | 'success' | 'error'
  lines: string[]
  createdAt: number
  endedAt?: number
}

export const useActionLogStore = defineStore('actionLog', {
  state: () => ({
    _nextId: 1,
    actions: [] as ActionEntry[],
  }),
  getters: {
    running: (state) => state.actions.filter((a) => a.status === 'running'),
    hasRunning: (state) => state.actions.some((a) => a.status === 'running'),
    finished: (state) => state.actions.filter((a) => a.status !== 'running'),
  },
  actions: {
    start(label: string, projectName?: string): number {
      const id = this._nextId++
      this.actions.unshift({ id, label, projectName, status: 'running', lines: [], createdAt: Date.now() })
      return id
    },
    append(id: number, line: string) {
      const a = this.actions.find((a) => a.id === id)
      if (a) a.lines.push(line)
    },
    appendLines(id: number, lines: string[]) {
      const a = this.actions.find((a) => a.id === id)
      if (a) a.lines.push(...lines)
    },
    end(id: number, status: 'success' | 'error' = 'success') {
      const a = this.actions.find((a) => a.id === id)
      if (a) { a.status = status; a.endedAt = Date.now() }
    },
    clearFinished() {
      this.actions = this.actions.filter((a) => a.status === 'running')
    },
    remove(id: number) {
      this.actions = this.actions.filter((a) => a.id !== id)
    },
  },
})
