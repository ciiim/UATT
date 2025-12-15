<template>
  <div class="toolbox">
    <draggable
      :list="library"
      :group="{ name: 'canvas', pull: 'clone', put: false }"
      :sort="false"
      :clone="onClone"
      item-key="type"
    >
      <template #item="{ element }">
        <div class="toolbox-item">
          {{ element.label }}
        </div>
      </template>
    </draggable>
  </div>
</template>

<script setup lang="ts">
import draggable from 'vuedraggable';

interface LibraryItem {
  type: string;
  label: string;
}

const library: LibraryItem[] = [
  { type: 'button', label: '按钮' },
  { type: 'input', label: '输入框' },
];

// 克隆函数：生成新对象
const onClone = (item: LibraryItem) => {
  return {
    id: Date.now(),
    type: item.type,
    label: item.label,
    position: { x: 50, y: 50 }, // 初始位置
  };
};
</script>

<style scoped>
.toolbox {
  width: 150px;
  border-right: 1px solid #ddd;
  padding: 8px;
}
.toolbox-item {
  padding: 8px;
  margin-bottom: 6px;
  background: #f9f9f9;
  border: 1px solid #ccc;
  cursor: grab;
  text-align: center;
}
</style>