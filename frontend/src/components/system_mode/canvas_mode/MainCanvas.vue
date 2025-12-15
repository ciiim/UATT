<template>
  <div class="canvas" ref="canvasRef">
    <draggable
      v-model="components"
      :group="{ name: 'canvas', pull: false, put: true }"
      item-key="id"
      :sort="false"
      class="canvas-draggable"
    >
      <template #item="{ element }">
        <WidgetWrapper
          :component="element"
          :onPositionChange="updatePosition"
        />
      </template>
    </draggable>
  </div>
</template>

<script setup lang="ts">
import draggable from 'vuedraggable';
import { ref } from 'vue';
import WidgetWrapper from './WidgetWrapper.vue';

interface CanvasComponent {
  id: string;
  type: string;
  label: string;
  position: { x: number; y: number };
}

const components = ref<CanvasComponent[]>([]);

const updatePosition = (id: string, pos: { x: number; y: number }) => {
  const item = components.value.find(c => c.id === id);
  if (item) {
    item.position = pos;
  }
};
</script>

<style scoped>
.canvas {
  flex: 1;
  position: relative;
  background: #fff;
}
.canvas-draggable {
  position: relative;
  width: 100%;
  height: 100%;
}
</style>