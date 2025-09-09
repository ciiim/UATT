<template>
  <div style="position: relative; width: 100%; height: 100%;">
      <draggable
        style="width: 100%; height: 100%;"
        v-model="contentModules"
        :group="{ name: 'mods', pull: true, put: true }"
        :animation="200"
        :ghostClass="'ghost'"
        :component-data="{
          type: 'transition-group',
          name: !drag ? 'flip-list' : null
        }"
        itemKey="id"
        @start="drag = true"
        @end="drag = false"
      >
        <template #item="{ element }">
          <div class="module-item" :class="{active: nowSelectedModule?.id === element.id}"  @click="selectModule(element)">
          <TestModuleCard :data="element" class="module-card"/>
          </div>
        </template>
      </draggable>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

import draggable from 'vuedraggable'

import { defineProps } from 'vue';

defineProps<{
  moduleLibrary: any[];
}>();

const drag = ref(false);

const contentModules = ref<any[]>([
]);

const nowSelectedModule = ref<any | null>(null);

const selectModule = (module : any) => {
  nowSelectedModule.value = module;
}

</script>

<style>

.module-card {
  margin: 5px;
  box-shadow: 0px 3px 5px #AAAAAA;
}

.module-item {
  border: 2px solid transparent;
  border-radius: 4px;
  cursor: pointer;
  transition: border-color 0.2s;
}
.module-item.active {
  background-color: rgba(24, 144, 255, 0.5);
}
.module-item:hover {
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