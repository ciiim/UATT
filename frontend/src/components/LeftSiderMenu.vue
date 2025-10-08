<template>
  <div class="sidebar">
    <!-- 顶部区域 -->
    <div class="sidebar-header">
      <span class="title">应用列表</span>
      <a-button type="primary" shape="circle" @click="showAddModal">
        <template #icon>
          <PlusOutlined></PlusOutlined>
        </template>
      </a-button>
    </div>

    <a-modal
    v-model:open="addModalVisible"
    title="新增应用"
    @ok="handleAddApp"
    width="600px"
  >
    <a-form
      :model="formData"
      :label-col="{ span: 6 }"
      :wrapper-col="{ span: 16 }"
    >
      <a-form-item label="应用名称" required>
        <a-input v-model:value="formData.AppName" />
      </a-form-item>

      <a-form-item label="串口波特率" required>
        <a-input-number v-model:value="formData.SerialConfig.BaudRate" :min="1" />
      </a-form-item>

      <a-form-item label="数据位" required>
        <a-input-number v-model:value="formData.SerialConfig.DataBits" :min="1" />
      </a-form-item>

      <a-form-item label="校验位" required>
        <a-select v-model:value="formData.SerialConfig.Parity">
          <a-select-option value="None">None</a-select-option>
          <a-select-option value="Odd">Odd</a-select-option>
          <a-select-option value="Even">Even</a-select-option>
        </a-select>
      </a-form-item>

      <a-form-item label="停止位" required>
        <a-input-number v-model:value="formData.SerialConfig.StopBits" :min="1" />
      </a-form-item>

      <a-form-item label="启用日志">
        <a-switch v-model:checked="formData.LogEnable" />
      </a-form-item>

      <a-form-item label="启用日志导出">
        <a-switch v-model:checked="formData.LogExportEnable" />
      </a-form-item>

      <a-form-item label="日志导出位置">
        <a-input v-model:value="formData.LogExportLoaction" />
      </a-form-item>
    </a-form>
  </a-modal>

    <!-- 列表区域 -->
    <div class="app-list">
      <a-card
        v-for="app in apps"
        :key="app"
        class="app-card"
        :class="{ active: app === selectedApp }"
        @click="selectApp(app)"
      >
        <div class="card-content">
          <span class="app-name">{{ app }}</span>
          <a-popconfirm
            title="确认删除该应用？"
            ok-text="是"
            cancel-text="否"
            @confirm="deleteApp(app)"
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
import { GetAllAppName, LoadApp, CreateApp, DeleteApp } from "../../wailsjs/go/bsd_testtool/Manager";
import { message } from "ant-design-vue";
import { useActionStore } from "../stores/action_store";

const apps = ref<string[]>([]);
const selectedApp = ref<string>("");
const store = useActionStore();

// 加载应用列表
const loadApps = async () => {
  try {
    apps.value = await GetAllAppName();
    console.log(apps.value);
  } catch (err) {
    message.error("加载应用列表失败");
  }
};


const emit = defineEmits<{
  (e: "app-loaded"): void;
}>();

onMounted(() => {
  loadApps();
});

const selectApp = async (app: string) => {
  if (app == selectedApp.value) {
    return;
  }
  selectedApp.value = app;
  store.nowApp = app;
  store.selectedAction = undefined;
  try {
    // message.loading('加载中')
    await LoadApp(app);
    message.info("加载完毕");

    emit("app-loaded");
  } catch (err) {
    message.error("加载 " + app + " 失败");
  }
};

// 弹窗控制
const addModalVisible = ref(false);

const showAddModal = () => {
  addModalVisible.value = true;
};

// 表单数据
const formData = ref({
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
});

// 提交表单
const handleAddApp = async () => {
  try {
    await CreateApp(formData.value);
    message.success("新增应用成功");
    addModalVisible.value = false;
    loadApps();
  } catch (err) {
    console.error(err);
    message.error("新增应用失败");
  }
};
const deleteApp = async (app: string) => {
  message.success(`已删除：${app}`);
  await DeleteApp(app)
  apps.value = apps.value.filter((a) => a !== app);
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

.app-list {
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
.app-name {
  font-weight: 500;
}
</style>
