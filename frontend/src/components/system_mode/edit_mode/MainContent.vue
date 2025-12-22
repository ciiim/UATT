<template>
  <div v-if="store.nowApp != ''" style="position: relative; width: 100%; height: 100%">
    <draggable
      style="width: 100%; height: 100%"
      v-model="store.actions"
      :group="{ name: 'mods', pull: true, put: true }"
      :animation="200"
      :ghostClass="'ghost'"
      :component-data="{
        type: 'transition-group',
        name: !drag ? 'flip-list' : null,
      }"
      itemKey="ActionUID"
      @start="drag = true"
      @end="drag = false"
      @change="onDragChange"
      
    >
      <template #item="{ element }">
        <div
          class="action-item"
          :class="{ active: store.selectedAction?.ActionUID === element.ActionUID }"
          :style="{ marginLeft: (element.indent ?? 0) * 20 + 'px' }"
          @click="selectAction(element)"
        >
          <TestModuleCard :data="element" class="action-card" />
        </div>
      </template>
    </draggable>
  </div>
  <div v-else>
    <h2 style="color: black;">请选择应用</h2>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from "vue";

import { GetActionList } from "../../../../wailsjs/go/bsd_testtool/Manager";

import { EventsOff, EventsOn } from "../../../../wailsjs/runtime/runtime"

import { parseActionTags, computeIndents } from "../../../utils/action_utils";

import type { ConfigActionBase, ActionReport } from "../../../types/Action";

import draggable from "vuedraggable";


import { useActionStore } from "../../../stores/action_store";
import { message } from "ant-design-vue";
import { bsd_testtool } from "../../../../wailsjs/go/models";

const prop = defineProps<{
  actionLibrary: any[];
  needGetActionList: boolean;
}>();

const store = useActionStore()


watch(
  () => prop.needGetActionList,
  () => {
    GetActionList()
      .then((res) => {
        const withTags = res.map((action) => ({
          ...action,
          Tags: parseActionTags(action), // 之前的标签解析
          Status: "无"
        }));
        store.actions = computeIndents(withTags); // 增加缩进属性
      })
      .catch((err) => {
        console.log("错误 " + err);
      });
  }
);

const onDragChange = () => {
  store.actions = computeIndents(store.actions)
}

const drag = ref(false);

onMounted(() => {
  console.log("mounted");

  EventsOn("begin-action", (data) => {
    store.actions.forEach(a => {
      a.Status = "未开始"
    })
  })
  
  EventsOn("now-action", (data : ActionReport) => {
    
    const actionIdx = store.actions.findIndex(a => a.ActionUID == data.ActionUID)
    
    store.actions[actionIdx].Status = "运行中"
  })
  EventsOn("done-action", (data : ActionReport) => {
    const actionIdx = store.actions.findIndex(a => a.ActionUID == data.ActionUID)
    if (data.Result != "success") {
      store.actions[actionIdx].Status = "失败"
      message.error(data.Result)
    } else {
      store.actions[actionIdx].Status = "完成"
    }
  })
})

onUnmounted(() => {
  EventsOff("begin-action")
  EventsOff("now-action")
  EventsOff("done-action")
})


const selectAction = (action: any) => {
  store.selectedAction = action;
  store.nowRightSiderTabIndex = 1;
  
  // message.info('选择:'+ nowSelectedAction.value?.ActionUID)
};
</script>

<style>
.action-card {
  margin: 1px;
  box-shadow: 0px 2px 2px #aaaaaa;
}

.action-item {
  border: 2px solid transparent;
  border-radius: 2px;
  cursor: pointer;
  transition: border-color 0.2s;
}
.action-item.active {
  background-color: rgba(24, 144, 255, 0.5);
}
.action-item:hover {
  transition: all 0.2s;
  background-color: rgba(57, 159, 255, 0.5); /* hover 时浅蓝色 */
}

.flip-list-move {
  transition: transform 0.2s;
}

.no-move {
  transition: transform 0s;
}

.ghost {
  opacity: 0.5;
  background: #c8ebfb;
}

.list-group {
  min-height: 20px;
}

.list-group-item {
  cursor: move;
}

.list-group-item i {
  cursor: pointer;
}
</style>
