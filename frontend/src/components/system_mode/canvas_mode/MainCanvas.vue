<template>
  <div class="container" v-if="store.nowCanvas != ''">
    <!-- 左侧组件库 -->
    <div class="sidebar">
      <div
        v-for="item in widgetLibrary"
        :key="item.type"
        class="library-item"
        draggable="true"
        @dragstart="onDragStart(item)"
      >
        {{ item.label }}
      </div>
    </div>

    <!-- 画布 -->
    <div
      class="canvas"
      ref="canvasRef"
      @dragover.prevent
      @drop="onDrop"
      @click="hideContextMenu"
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
          'widget-selected': selectedComponentId === comp.ID,
          'widget-button': comp.Type === 'button',
          'widget-text': comp.Type === 'text'
        }"
        :style="{ left: comp.Position.X + 'px', top: comp.Position.Y + 'px' }"
        @mousedown="startDrag(comp, $event)"
        @contextmenu="showContextMenu($event, comp)"
        @click.stop="selectComponent(comp)"
        @dblclick="handleDoubleClick(comp)"
      >
        {{ comp.Label }}
        <div v-if="comp.Type === 'text' && comp.Value" class="text-output">
          {{ comp.Value }}
        </div>
      </div>

      <!-- 右键菜单 -->
      <div
        v-if="contextMenu.visible"
        class="context-menu"
        :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }"
      >
        <a-menu>
          <a-menu-item key="edit" @click="editComponent">
            <EditOutlined />
            编辑
          </a-menu-item>
          
          <!-- 按钮特有的关联选项 -->
<template v-if="contextMenu.targetComponent?.Type === 'button'">
  <a-menu-item key="connect-action" @click="showActionDialog">
    <SettingOutlined />
    关联到Action
  </a-menu-item>
  <a-menu-item 
    v-if="getButtonConnections(contextMenu.targetComponent.ID).length > 0"
    key="disconnect" 
    @click="showDisconnectDialog"
  >
    <DisconnectOutlined />
    解除关联
  </a-menu-item>
</template>

<!-- 文字框特有的关联选项 -->
<template v-if="contextMenu.targetComponent?.Type === 'text'">
  <a-menu-item key="connect" @click="showConnectDialog">
    <LinkOutlined />
    关联到按钮
  </a-menu-item>
  <a-menu-item 
    v-if="getTextConnections(contextMenu.targetComponent.ID).length > 0"
    key="disconnect" 
    @click="showDisconnectDialog"
  >
    <DisconnectOutlined />
    解除关联
  </a-menu-item>
</template>

          <a-menu-item key="delete" @click="deleteComponent">
            <DeleteOutlined />
            删除
          </a-menu-item>
        </a-menu>
      </div>
    </div>

    <!-- 编辑弹窗 -->
    <a-modal
      v-model:open="editModalVisible"
      title="编辑组件"
      @ok="saveEdit"
      @cancel="cancelEdit"
    >
      <a-form :model="editForm" layout="vertical">
        <a-form-item label="标签文本">
          <a-input v-model:value="editForm.label" placeholder="请输入标签文本" />
        </a-form-item>
        <a-form-item label="类型">
          <a-select v-model:value="editForm.type" placeholder="请选择类型">
            <a-select-option value="button">按钮</a-select-option>
            <a-select-option value="text">文字框</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- Action关联弹窗 -->
<a-modal
  v-model:open="appModalVisible"
  title="关联到Action"
  @ok="saveAction"
  @cancel="cancelAction"
>
  <a-form layout="vertical">
    <a-form-item label="选择Action">
      <a-select v-model:value="appForm.selectedApp" placeholder="请选择Action">
        <a-select-option 
          v-for="action in availableActions"
          :key="action"
          :value="action"
        >
          {{ action }}
        </a-select-option>
      </a-select>
    </a-form-item>
  </a-form>
</a-modal>

    <!-- 关联弹窗 -->
    <a-modal
      v-model:open="connectModalVisible"
      title="关联到按钮"
      @ok="saveConnection"
      @cancel="cancelConnection"
    >
      <a-form layout="vertical">
        <a-form-item label="选择要关联的按钮">
          <a-select v-model:value="connectForm.targetId" placeholder="请选择按钮">
            <a-select-option 
              v-for="buttonComp in buttonComponents"
              :key="buttonComp.ID"
              :value="buttonComp.ID"
            >
              {{ buttonComp.Label }}
            </a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 解绑弹窗 -->
    <a-modal
      v-model:open="disconnectModalVisible"
      title="解除关联"
      @ok="saveDisconnection"
      @cancel="cancelDisconnection"
    >
      <a-form layout="vertical">
        <a-form-item label="选择要解除的关联">
          <a-select v-model:value="disconnectForm.connectionId" placeholder="请选择要解除的关联">
            <a-select-option 
              v-for="conn in currentConnections"
              :key="`${conn.FromID}-${conn.ToID}`"
              :value="`${conn.FromID}-${conn.ToID}`"
            >
              {{ getConnectionLabel(conn) }}
            </a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
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

const prop = defineProps<{
  needGetCanvasList: boolean;
}>();


watch(
  () => prop.needGetCanvasList,
  () => {
    GetCanvasData()
      .then((res) => {
        console.log(res);
        store.canvasComponents = res.Data.ComponentList || []
        store.canvasConnections = res.Data.Connections || []
      })
      .catch((err) => {
        console.log("错误 " + err);
      });
  }
);


// 组件库
const widgetLibrary = ref<LibraryItem[]>([
  { type: 'button', label: '按钮' },
  { type: 'text', label: '文字框' },
]);

// 画布尺寸
const canvasSize = reactive({ width: 800, height: 600 });

// 选中的组件ID
const selectedComponentId = ref<string | null>(null);

// 运行中的组件ID
const runningComponentID = ref<string | null>(null);

// 拖动数据缓存
const dragData = ref<LibraryItem | null>(null);

// 右键菜单状态
const contextMenu = reactive({
  visible: false,
  x: 0,
  y: 0,
  targetComponent: null as CanvasComponent | null
});

// Action弹窗状态
const appModalVisible = ref(false);
const appForm = reactive({
  buttonId: '',
  selectedApp: ''
});

// Action列表
const availableActions = ref<string[]>([]);

// 编辑弹窗状态
const editModalVisible = ref(false);
const editForm = reactive({
  label: '',
  type: '',
  id: '',
});

// 关联弹窗状态
const connectModalVisible = ref(false);
const connectForm = reactive({
  myselfId: '',
  targetId: ''
});

// 解绑弹窗状态
const disconnectModalVisible = ref(false);
const disconnectForm = reactive({
  connectionId: ''
});

// 计算属性：获取所有按钮组件
const buttonComponents = computed(() => 
  store.canvasComponents.filter(c => c.Type === 'button')
);

// 计算属性：当前组件的关联关系
const currentConnections = computed(() => {
  if (!contextMenu.targetComponent) return [];
  
  const compId = contextMenu.targetComponent.ID;
  if (contextMenu.targetComponent.Type === 'button') {
    return store.canvasConnections.filter(c => c.FromID === compId);
  } else if (contextMenu.targetComponent.Type === 'text') {
    return store.canvasConnections.filter(c => c.ToID === compId);
  }
  return [];
});

// 获取组件中心点坐标
const getComponentCenter = (componentId: string) => {
  const comp =  store.canvasComponents.find(c => c.ID === componentId);
  if (!comp) return { x: 0, y: 0 };
  
  return {
    x: comp.Position.X + 50, // 组件宽度的一半
    y: comp.Position.Y + 20  // 组件高度的一半
  };
};

// 获取按钮的关联关系
const getButtonConnections = (buttonId: string) => {
  return store.canvasConnections.filter(c => c.FromID === buttonId);
};

// 获取文字框的关联关系
const getTextConnections = (textId: string) => {
  return store.canvasConnections.filter(c => c.ToID === textId);
};

// 获取关联关系的显示标签
const getConnectionLabel = (conn: Connection) => {
  const fromComp =  store.canvasComponents.find(c => c.ID === conn.FromID);
  const toComp =  store.canvasComponents.find(c => c.ID === conn.ToID);
  return `${fromComp?.Label || '未知'} → ${toComp?.Label || '未知'}`;
};

// 更新画布尺寸
const updateCanvasSize = () => {
  if (canvasRef.value) {
    canvasSize.width = canvasRef.value.clientWidth;
    canvasSize.height = canvasRef.value.clientHeight;
  }
};

// 原生拖入事件开始
const onDragStart = (item: LibraryItem) => {
  dragData.value = item;
};

// 放到画布内
const onDrop = (e: DragEvent) => {
  if (!dragData.value) return;
  const canvasBounds = (canvasRef.value as HTMLDivElement).getBoundingClientRect();

  const newComp: CanvasComponent = {
    ID: `${dragData.value.type}_${Date.now()}`,
    Type: dragData.value.type,
    Label: dragData.value.label,
    Position: {
      X: e.clientX - canvasBounds.left - 40,
      Y: e.clientY - canvasBounds.top - 20
    }
  };

  console.log('id:' +newComp.ID);
  

  // 为文字框初始化 value，为按钮初始化 action
  if (newComp.Type === 'text') {
    newComp.Value = '';
  } else if (newComp.Type === 'button') {
    newComp.AttachApp = '';
  }

   store.canvasComponents.push(newComp);
  dragData.value = null;
};

// 自由拖动相关
const draggingComp = ref<CanvasComponent | null>(null);
let offsetX = 0;
let offsetY = 0;

const startDrag = (comp: CanvasComponent, e: MouseEvent) => {
  if (e.button === 2) return;
  
  draggingComp.value = comp;
  offsetX = e.clientX - comp.Position.X;
  offsetY = e.clientY - comp.Position.Y;
  window.addEventListener('mousemove', onMouseMove);
  window.addEventListener('mouseup', stopDrag);
};

const onMouseMove = (e: MouseEvent) => {
  if (draggingComp.value) {
    draggingComp.value.Position.X = e.clientX - offsetX;
    draggingComp.value.Position.Y = e.clientY - offsetY;
  }
};

const stopDrag = () => {
  draggingComp.value = null;
  window.removeEventListener('mousemove', onMouseMove);
  window.removeEventListener('mouseup', stopDrag);
};

// 选中组件
const selectComponent = (comp: CanvasComponent) => {
  selectedComponentId.value = comp.ID;
};

// 双击处理（按钮执行操作）
const handleDoubleClick = (comp: CanvasComponent) => {
  if (comp.Type === 'button' && comp.AttachApp != '') {
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



// 显示Action选择对话框
const showActionDialog = async () => {
  if (!contextMenu.targetComponent) return;
  
  try {
    availableActions.value = await GetAllAppName();
    appForm.buttonId = contextMenu.targetComponent.ID;
    appForm.selectedApp = contextMenu.targetComponent.AttachApp || '';
    appModalVisible.value = true;
    hideContextMenu();
  } catch (error) {
    message.error('获取Action列表失败');
  }
};

// 保存Action关联
const saveAction = () => {
  if (!appForm.buttonId || !appForm.selectedApp) {
    message.warning('请选择一个Action');
    return;
  }

  const targetButton = store.canvasComponents.find(c => c.ID === appForm.buttonId);
  if (targetButton) {
    targetButton.AttachApp = appForm.selectedApp;
    message.success('App关联成功');
  }
  
  appModalVisible.value = false;
};

// 取消Action关联
const cancelAction = () => {
  appForm.buttonId = '';
  appForm.selectedApp = '';
  appModalVisible.value = false;
};

// ===== App接口结束 =====

// 显示右键菜单
const showContextMenu = (e: MouseEvent, comp: CanvasComponent) => {
  e.preventDefault();
  contextMenu.visible = true;
  contextMenu.x = e.pageX;
  contextMenu.y = e.pageY;
  contextMenu.targetComponent = comp;
  selectedComponentId.value = comp.ID;
  console.log(comp);
  
};

// 隐藏右键菜单
const hideContextMenu = () => {
  contextMenu.visible = false;
  // contextMenu.targetComponent = null;
};

const showConnectDialog = () => {
  if (buttonComponents.value.length === 0) {
    message.warning('画布上没有按钮可以关联');
    return;
  }
  connectForm.targetId = '';
  connectForm.myselfId = contextMenu.targetComponent!.ID;
  connectModalVisible.value = true;
  hideContextMenu();
};

const saveConnection = () => {
  if (!connectForm.myselfId || !connectForm.targetId) {
    message.warning('请选择要关联的按钮');
    return;
  }

  const newConnection: Connection = {
    FromID: connectForm.targetId,  // 现在是按钮ID
    ToID: connectForm.myselfId  // 现在是文字框ID
  };

  // 检查是否已经存在相同的关联
  const exists = store.canvasConnections.some(
    c => c.FromID === newConnection.FromID && c.ToID === newConnection.ToID
  );

  if (exists) {
    message.warning('该关联已存在');
    return;
  }

  store.canvasConnections.push(newConnection);
  message.success('关联创建成功');
  connectModalVisible.value = false;
};

// 取消关联
const cancelConnection = () => {
  connectModalVisible.value = false;
};

// 显示解绑对话框
const showDisconnectDialog = () => {
  disconnectForm.connectionId = '';
  disconnectModalVisible.value = true;
  hideContextMenu();
};

// 保存解绑
const saveDisconnection = () => {
  if (!disconnectForm.connectionId) {
    message.warning('请选择要解除的关联');
    return;
  }

  console.log(disconnectForm.connectionId);
  

  const [fromId, toId] = disconnectForm.connectionId.split('-');
  const index = store.canvasConnections.findIndex(
    c => c.FromID === fromId && c.ToID === toId
  );

  if (index > -1) {
    store.canvasConnections.splice(index, 1);
    message.success('关联解除成功');
  }

  disconnectModalVisible.value = false;
};

// 取消解绑
const cancelDisconnection = () => {
  disconnectModalVisible.value = false;
};

// 删除组件
const deleteComponent = () => {
  if (!contextMenu.targetComponent) return;
  const targetComponent = contextMenu.targetComponent;
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除这个组件吗？删除后相关的连接关系也会被清除。',
    okText: '确认',
    cancelText: '取消',
    onOk() {

      const compId = targetComponent!.ID;
      
      // 删除组件
      const index = store.canvasComponents.findIndex(c => c.ID === compId);
      if (index > -1) {
        store.canvasComponents.splice(index, 1);
      }

      // 删除相关的连接关系
      store.canvasConnections = store.canvasConnections.filter(
        c => c.FromID !== compId && c.ToID !== compId
      );

      if (selectedComponentId.value === compId) {
        selectedComponentId.value = null;
      }

      message.success('组件删除成功');
      hideContextMenu();
    }
  });
};

// 编辑组件
const editComponent = () => {
  if (!contextMenu.targetComponent) return;
  
  editForm.label = contextMenu.targetComponent.Label;
  editForm.type = contextMenu.targetComponent.Type;
  editForm.id = contextMenu.targetComponent.ID;
  editModalVisible.value = true;
  hideContextMenu();
};

// 保存编辑
const saveEdit = () => {
  if (!editForm.id) return;
  
  const targetComp = store.canvasComponents.find(
    c => c.ID === editForm.id
  );
  if (targetComp) {
    targetComp.Label = editForm.label;
    targetComp.Type = editForm.type;
    targetComp.ID = editForm.id;
    message.success('组件编辑成功');
  }
  editModalVisible.value = false;
};

// 取消编辑
const cancelEdit = () => {
  editModalVisible.value = false;
};

const canvasRef = ref<HTMLDivElement | null>(null);

onMounted(() => {
  updateCanvasSize();
  window.addEventListener('resize', updateCanvasSize);

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
  cursor: move;
  user-select: none;
  padding: 8px 12px;
  border: 1px solid #999;
  border-radius: 4px;
  transition: border-color 0.2s;
  z-index: 2;
  min-width: 80px;
}

.widget:hover {
  border-color: #1890ff;
}

.widget-selected {
  border-color: #1890ff;
  border-width: 2px;
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