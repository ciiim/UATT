import { reactive, ref } from 'vue'
import type { ConfigActionBase } from '../types/Action'
import { CanvasComponent, Connection } from '../types/Canvas'

const state = reactive({
  nowApp: '',
  nowCanvas: '',
  canvasComponents: [] as CanvasComponent[],
  canvasConnections: [] as Connection[],
  actions: [] as ConfigActionBase[],
  selectedAction: undefined as ConfigActionBase | undefined,
  actionListChanges: [],
  nowRightSiderTabIndex: 0,
  logs: [] as string[]
})


export function useActionStore() {
  return state
}