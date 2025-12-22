<template>
  <div class="container" v-if="store.nowCanvas != ''">

    <!-- 画布 -->
    <div
      class="canvas"
      ref="canvasRef"
    >
      <!-- SVG 连线层 -->
      <svg class="connection-layer" :width="canvasSize.width" :height="canvasSize.height">
        <line
          v-for="conn in store.canvasConnections"
          :key="`${conn.FromID}-${conn.ToID}`"
          :x1="getComponentCenter(conn.FromID).x"
          :y1="getComponentCenter(conn.FromID).y"
          :x2="getComponentCenter(conn.ToID).x"
          :y2="getComponentCenter(conn.ToID).y"
          stroke="#1890ff"
          stroke-width="2"
          marker-end="url(#arrowhead)"
        />
        <!-- 箭头标记定义 -->
        <defs>
          <marker id="arrowhead" markerWidth="10" markerHeight="7" 
            refX="9" refY="3.5" orient="auto">
            <polygon points="0 0, 10 3.5, 0 7" fill="#1890ff" />
          </marker>
        </defs>
      </svg>

      <!-- 组件层 -->
      <div
        v-for="comp in store.canvasComponents"
        :key="comp.ID"
        class="widget"
        :class="{ 
          'widget-button': comp.Type === 'button',
          'widget-text': comp.Type === 'text'
        }"
        :style="{ left: comp.Position.X + 'px', top: comp.Position.Y + 'px' }"
        @click="() => {handleClick(comp)}"
      >
        {{ comp.Label }}
        <div v-if="comp.Type === 'text' && comp.Value" class="text-output">
          {{ comp.Value }}
        </div>
      </div>
    </div>
  </div>
  <div v-else>
    <h2 style="color: black;">请选择画布</h2>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, watch } from 'vue';
import { Modal, message } from 'ant-design-vue';
import { EditOutlined, DeleteOutlined, LinkOutlined, DisconnectOutlined } from '@ant-design/icons-vue';
import { GetAllAppName, GetCanvasData, StartCanvasApp } from '../../../../wailsjs/go/bsd_testtool/Manager';
import { CanvasComponent, Connection } from '../../../types/Canvas';
import { useActionStore } from "../../../stores/action_store";
import { EventsOn, EventsOff } from '../../../../wailsjs/runtime/runtime';

const store = useActionStore();

interface LibraryItem {
  type: string;
  label: string;
}

// 画布尺寸
const canvasSize = reactive({ width: 800, height: 600 });


// 运行中的组件ID
const runningComponentID = ref<string | null>(null);

// 获取组件中心点坐标
const getComponentCenter = (componentId: string) => {
  const comp =  store.canvasComponents.find(c => c.ID === componentId);
  if (!comp) return { x: 0, y: 0 };
  
  return {
    x: comp.Position.X + 50, // 组件宽度的一半
    y: comp.Position.Y + 20  // 组件高度的一半
  };
};


// 更新画布尺寸
const updateCanvasSize = () => {
  if (canvasRef.value) {
    canvasSize.width = canvasRef.value.clientWidth;
    canvasSize.height = canvasRef.value.clientHeight;
  }
};

const sleep = (ms: number): Promise<void> => {
  return new Promise(resolve => setTimeout(resolve, ms));
};

const handleClick = async (comp: CanvasComponent) => {
  if (comp.Type === 'button' && comp.AttachApp != '') {
    
    updateConnectedTextBoxes(comp.ID, 0, "")
    updateConnectedTextBoxes(comp.ID, 1, "")
    executeButtonApp(comp);
  } else {
    console.log('executeButtonApp failed:', comp);
  }
};

const executeButtonApp = async (buttonComp: CanvasComponent) => {

  runningComponentID.value = buttonComp.ID;

  
  console.log('按钮执行:', buttonComp.Label, '操作App:', buttonComp.AttachApp);

  StartCanvasApp(buttonComp.AttachApp!).catch((err)=> {
    message.error(["执行失败", err], 1)
  })
};

// 更新关联的文字框
const updateConnectedTextBoxes = (buttonId: string, index : number, result: string) => {
  const relatedConnections = store.canvasConnections.filter(c => c.FromID === buttonId);
  // 验证索引是否有效
  if (isNaN(index) || index >= relatedConnections.length) {
    return;
  }
  
  const connection = relatedConnections[index]; // 转换为0基础索引
  const textComp = store.canvasComponents.find(c => c.ID === connection.ToID);
  
  if (textComp) {
    textComp.Value = result;
  } else {
    console.warn(`Text component not found for connection: ${connection.ToID}`);
  }
};


const canvasRef = ref<HTMLDivElement | null>(null);

onMounted(() => {
  updateCanvasSize();
  window.addEventListener('resize', updateCanvasSize);

  GetCanvasData()
    .then((res) => {
      console.log(res);
      store.canvasComponents = res.Data.ComponentList
      store.canvasConnections = res.Data.Connections
    })
    .catch((err) => {
      console.log("错误 " + err);
    });

  EventsOn("output_text", (idx, str) => {
    updateConnectedTextBoxes(runningComponentID.value!, idx, str)
  })

});

onUnmounted(() => {
  window.removeEventListener('resize', updateCanvasSize);

  EventsOff("output_text")
});
</script>

<style scoped>
.container {
  display: flex;
  height: 100vh;
}

/* 左侧组件库 */
.sidebar {
  width: 200px;
  background: #f4f4f4;
  padding: 10px;
  box-sizing: border-box;
}
.library-item {
  padding: 8px;
  margin-bottom: 8px;
  background: #ddd;
  border: 1px solid #ccc;
  cursor: grab;
  border-radius: 4px;
}
.library-item:active {
  cursor: grabbing;
}

/* 画布区域 */
.canvas {
  flex: 1;
  position: relative;
  background: #eee;
  border: 1px solid #ccc;
}

/* SVG 连线层 */
.connection-layer {
  position: absolute;
  top: 0;
  left: 0;
  pointer-events: none;
  z-index: 1;
}

/* 画布组件 */
.widget {
  position: absolute;
  cursor: pointer;
  user-select: none;
  padding: 8px 12px;
  border: 1px solid #999;
  border-radius: 4px;
  transition: border-color 0.2s;
  transition: transform 0.1s ease;
  z-index: 2;
  min-width: 80px;
}

.widget:hover {
  border-color: #3bff18;
  border-width: 1px;
}


.widget-button:active {
  transform: scale(0.95);
}

.widget-button {
  background: #5a6266;
}

.widget-text {
  background: #7d9763;
  min-height: 40px;
}

.text-output {
  margin-top: 4px;
  padding: 4px;
  background: #fff;
  border: 1px solid #d9d9d9;
  border-radius: 2px;
  font-size: 12px;
  color: #666;
}

/* 右键菜单 */
.context-menu {
  position: fixed;
  z-index: 1000;
  background: white;
  border-radius: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.context-menu :deep(.ant-menu) {
  border: none;
  box-shadow: none;
}

.context-menu :deep(.ant-menu-item) {
  padding: 8px 16px;
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>