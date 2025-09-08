<template>
  <div style="position: relative; width: 100%; height: 100%;">
      <a-float-button
      tooltip="添加组件"
      :style="{
        position: 'absolute',
        right: '24px',
      }"
    >
      <template #icon>
        <PlusOutlined></PlusOutlined>
      </template>
  </a-float-button>
  </div>
</template>

<script setup lang="ts">
import { PlusOutlined } from '@ant-design/icons-vue';



</script>