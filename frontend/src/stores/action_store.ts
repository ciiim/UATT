import { reactive, ref } from 'vue'
import type { ConfigActionBase } from '../types/Action'

const state = reactive({
  nowApp: '',
  actions: [] as ConfigActionBase[],
  selectedAction: undefined as ConfigActionBase | undefined,
  actionListChanges: [],
  nowRightSiderTabIndex: 0,
  logs: [] as string[]
})


export function useActionStore() {
  return state
}