<template>
  <a-card title="组件" :style="{overflow: 'auto', borderRadius: '0'}">
  <draggable
    :list="actionLibrary"
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
import type { ConfigActionBase } from '../../types/Action';
import draggable from 'vuedraggable';

defineProps<{ actionLibrary: any[] }>();

const onClone = (evt: any) => {

  let newAction : ConfigActionBase = {
    ActionUID: Date.now(),
    ActionType: evt.actionType,
    ActionTypeID: evt.actionID,
    Name: evt.name,
    TypeFeatureField: evt.feat,
    BreakPoint: false,
  } 
  
  return newAction
}

</script>