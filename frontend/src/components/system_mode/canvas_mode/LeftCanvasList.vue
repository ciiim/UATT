<template>
  <div class="sidebar">
    <!-- 顶部区域 -->
    <div class="sidebar-header">
      <span class="title">画布列表</span>
      <a-button type="primary" shape="circle" @click="addCanvas">
        <template #icon>
          <PlusOutlined></PlusOutlined>
        </template>
      </a-button>
    </div>

    <!-- 列表区域 -->
    <div class="canvas-list">
      <a-card
        v-for="canvas in canvasList"
        :key="canvas"
        class="app-card"
        :class="{ active: canvas === selectedCanvas }"
        @click="selectCanvas(canvas)"
      >
        <div class="card-content">
          <span class="canvas-name">{{ canvas }}</span>
          <a-popconfirm
            title="确认删除该画布？"
            ok-text="是"
            cancel-text="否"
            @confirm="deleteCanvas(canvas)"
            @click.stop
          >
            <a-button type="text" danger>
              <template #icon>
                <DeleteOutlined></DeleteOutlined>
              </template>
            </a-button>
          </a-popconfirm>
        </div>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { DeleteOutlined, PlusOutlined } from "@ant-design/icons-vue";
import { GetAllAppName, LoadApp, CreateApp, DeleteApp } from "../../../../wailsjs/go/bsd_testtool/Manager";
import { message } from "ant-design-vue";
import { useActionStore } from "../../../stores/action_store";
import { bsd_testtool } from "../../../../wailsjs/go/models";

const canvasList = ref<string[]>([]);
const selectedCanvas = ref<string>("");
const store = useActionStore();

// 加载应用列表
const loadCanvasList = async () => {
  try {
  } catch (err) {
    message.error("加载应用列表失败");
  }
};


const emit = defineEmits<{
  (e: "canvas-loaded"): void;
}>();

onMounted(() => {
  loadCanvasList();
});

const selectCanvas = async (app: string) => {
  if (app == selectedCanvas.value) {
    return;
  }
  selectedCanvas.value = app;
  store.nowApp = app;
  store.selectedAction = undefined;
  try {
    // message.loading('加载中')
    // await LoadApp(app);
    message.info("加载完毕", 1);

    emit("canvas-loaded");
  } catch (err) {
    message.error("加载 " + app + " 失败");
  }
};


const addCanvas = () => {

};

// 表单数据
const formData = ref(new bsd_testtool.AppConfigSettings({
  AppName: "",
  SerialConfig: {
    BaudRate: 9600,
    DataBits: 8,
    Parity: "None",
    StopBits: 1,
  },
  LogEnable: false,
  LogExportEnable: false,
  LogExportLoaction: "",
}));

// 提交表单
const handleAddCanvas = async () => {
  try {
    // await CreateApp(formData.value);
    message.success("新增画布成功");
    loadCanvasList();
  } catch (err) {
    console.error(err);
    message.error("新增画布失败");
  }
};
const deleteCanvas = async (app: string) => {
  message.success(`已删除：${app}`);
//   await DeleteApp(app)
  canvasList.value = canvasList.value.filter((a) => a !== app);
};
</script>

<style scoped>
.sidebar {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px;
  border-bottom: 1px solid #eee;
}

.canvas-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.app-card {
  margin-bottom: 8px;
  cursor: pointer;
  transition: all 0.2s;
}
.app-card:hover {
  border-color: #1890ff;
}
.app-card.active {
  border: 1px solid #1890ff;
  background-color: #e6f7ff;
}

.card-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.canvas-name {
  font-weight: 500;
}
</style>
