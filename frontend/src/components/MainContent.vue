<template>
  <div style="position: relative; width: 100%; height: 100%">
    <draggable
      style="width: 100%; height: 100%"
      v-model="contentActions"
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
          :class="{ active: nowSelectedAction?.ActionUID === element.ActionUID }"
          :style="{ marginLeft: (element.indent ?? 0) * 20 + 'px' }"
          @click="selectAction(element)"
        >
          <TestModuleCard :data="element" class="action-card" />
        </div>
      </template>
    </draggable>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";

import { GetActionList } from "../../wailsjs/go/bsd_testtool/Manager";

import { parseActionTags, computeIndents } from "../utils/action_utils";

import type { ConfigActionBase } from "../types/Action";

import draggable from "vuedraggable";

import { defineProps } from "vue";
import { useActionStore } from "../stores/action_store";

const prop = defineProps<{
  actionLibrary: any[];
  needGetActionList: boolean;
}>();

const store = useActionStore()

const contentActions = ref<ConfigActionBase[]>([]);

const nowSelectedAction = ref<ConfigActionBase | undefined>();

watch(
  () => prop.needGetActionList,
  () => {
    GetActionList()
      .then((res) => {
        const withTags = res.map((action) => ({
          ...action,
          tags: parseActionTags(action), // 之前的标签解析
        }));
        contentActions.value = computeIndents(withTags); // 增加缩进属性
        store.actions = contentActions.value
      })
      .catch((err) => {
        console.log("错误 " + err);
      });
  }
);

const onDragChange = () => {
  contentActions.value = computeIndents(contentActions.value)
  store.actions = contentActions.value
}

const drag = ref(false);



const selectAction = (action: any) => {
  nowSelectedAction.value = action;
  store.selectedAction = action;
  
};
</script>

<style>
.action-card {
  margin: 5px;
  box-shadow: 0px 3px 5px #aaaaaa;
}

.action-item {
  border: 2px solid transparent;
  border-radius: 4px;
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
