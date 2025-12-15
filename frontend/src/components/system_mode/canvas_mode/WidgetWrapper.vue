<template>
  <div
    class="widget-wrapper"
    :style="wrapperStyle"
    @mousedown="onMouseDown"
  >
    <template v-if="component.type === 'button'">
      <button>{{ component.label }}</button>
    </template>
    <template v-else-if="component.type === 'input'">
      <input placeholder="请输入文字" />
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { CSSProperties } from 'vue';

interface CanvasComponent {
  id: string;
  type: string;
  label: string;
  position: { x: number; y: number };
}

const props = defineProps<{
  component: CanvasComponent;
  onPositionChange: (id: string, pos: { x: number; y: number }) => void;
}>();

// 组件定位样式
const wrapperStyle = computed<CSSProperties>(() => ({
  position: 'absolute',
  left: `${props.component.position.x}px`,
  top: `${props.component.position.y}px`,
  cursor: 'move',
}));

const onMouseDown = (e: MouseEvent) => {
  const startX = e.clientX;
  const startY = e.clientY;
  const { x, y } = props.component.position;

  const onMouseMove = (moveEvent: MouseEvent) => {
    const dx = moveEvent.clientX - startX;
    const dy = moveEvent.clientY - startY;
    props.onPositionChange(props.component.id, { x: x + dx, y: y + dy });
  };

  const onMouseUp = () => {
    document.removeEventListener('mousemove', onMouseMove);
    document.removeEventListener('mouseup', onMouseUp);
  };

  document.addEventListener('mousemove', onMouseMove);
  document.addEventListener('mouseup', onMouseUp);
};
</script>

<style scoped>
.widget-wrapper {
  user-select: none;
}
</style>