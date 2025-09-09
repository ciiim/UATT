<template>
  <a-card title="组件" :style="{overflow: 'auto', borderRadius: '0'}">
  <draggable
    :list="moduleLibrary"
    :group="{name:'mods', pull: 'clone', put: false}"
    :sort="false"
    :clone="onClone"
    itemKey="name"
  >
    <template #item="{ element, index }">
      <a-card-grid style="width: 100%; text-align: center">{{ element.name }}</a-card-grid>
    </template>
  </draggable>
  </a-card>
</template>

<script setup lang="ts">
import { ref, defineProps } from 'vue';
import draggable from 'vuedraggable';

defineProps<{ moduleLibrary: any[] }>();

const onClone = (evt: any) => {
  console.log('modUID' + evt.modUid);
  
  return {
    name: evt.name,
    modUid: evt.modUid,
    tags: [],
    timeout: 200,
    status: '无',
    id: Date.now()
  }
}

</script>