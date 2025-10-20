<template>
  <div class="header-toolbar">
    <div class="left">
      <!-- 左边内容 -->
    </div>
    <div class="right">
      <!-- 串口下拉选择 -->
      <a-select
        v-model:value="selectedSerial"
        placeholder="选择串口"
        style="width: 160px"
        @change="onSelectSerial"
        @dropdownVisibleChange="GetSerialList"
      >
        <a-select-option v-for="port in serialList" :key="port" :value="port">
          {{ port }}
        </a-select-option>
      </a-select>
      <a-button class="btn" type="primary" v-if="nowRunningStatus == 0" @click="startAction">运行</a-button>
      <a-button class="btn" type="primary" v-if="nowRunningStatus == 2" >等待结束</a-button>
      <a-button class="btn" type="primary" v-if="nowRunningStatus == 1"  @click="stopAction">停止</a-button>
      <a-button class="btn" type="primary" @click="onSave">保存</a-button>
      <a-button class="btn" type="primary" @click="openSetting">设置</a-button>
    </div>
  </div>
  <a-modal
    v-model:open="addModalVisible"
    title="新增应用"
    @ok="handleAddApp"
    width="600px"
  >
    <a-form
      :model="configSettings"
      :label-col="{ span: 6 }"
      :wrapper-col="{ span: 16 }"
    >
      <a-form-item label="应用名称" required>
        <a-input v-model:value="configSettings.AppName" />
      </a-form-item>

      <a-form-item label="串口波特率" required>
        <a-input-number
          v-model:value="configSettings.SerialConfig.BaudRate"
          :min="1"
        />
      </a-form-item>

      <a-form-item label="数据位" required>
        <a-input-number
          v-model:value="configSettings.SerialConfig.DataBits"
          :min="1"
        />
      </a-form-item>

      <a-form-item label="校验位" required>
        <a-select v-model:value="configSettings.SerialConfig.Parity">
          <a-select-option value="None">None</a-select-option>
          <a-select-option value="Odd">Odd</a-select-option>
          <a-select-option value="Even">Even</a-select-option>
        </a-select>
      </a-form-item>

      <a-form-item label="停止位" required>
        <a-input-number
          v-model:value="configSettings.SerialConfig.StopBits"
          :min="1"
        />
      </a-form-item>

      <a-form-item label="启用日志">
        <a-switch v-model:checked="configSettings.LogEnable" />
      </a-form-item>

      <a-form-item label="启用日志导出">
        <a-switch v-model:checked="configSettings.LogExportEnable" />
      </a-form-item>

      <a-form-item label="日志导出位置">
        <a-input v-model:value="configSettings.LogExportLoaction" />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { message } from "ant-design-vue";
import { useActionStore } from "../stores/action_store";
import {
  GetAllSerial,
  SelectSerialCom,
  Start,
  Stop,
  GetAppSettings,
  SyncAppSettings,
  SaveApp,
  SyncActions,
} from "../../wailsjs/go/bsd_testtool/Manager";
import { bsd_testtool } from "../../wailsjs/go/models";
import { EventsOn } from "../../wailsjs/runtime/runtime"

const actionStore = useActionStore();

// 弹窗控制
const addModalVisible = ref(false);
// 串口列表和选中值
const serialList = ref<string[]>([]);
const selectedSerial = ref<string>("");
const configSettings = ref<bsd_testtool.AppConfigSettings>(
  new bsd_testtool.AppConfigSettings()
);

// 页面加载时获取串口列表
onMounted(async () => {
  
  EventsOn('stopped', (data) => {
    nowRunningStatus.value = 0
  })
  
  await GetSerialList()
  if (serialList.value.length > 0) {
      selectedSerial.value = serialList.value[0];
      // 默认选中第一个串口
      await SelectSerialCom(selectedSerial.value);
    }

});

const nowRunningStatus = ref<number>(0);

const GetSerialList = async ()=> {
  try {
    
    serialList.value = await GetAllSerial();
  } catch (err) {
    console.error(err);
    message.error("获取串口列表失败");
  }
} 

// 选择串口
const onSelectSerial = async (port: string) => {
  try {
    console.log('select');
    
    await SelectSerialCom(port);
    //message.success(`已选择串口: ${port}`);
  } catch (err) {
    console.error(err);
    message.error("选择串口失败");
  }
};

const startAction = async () => {
  try {
    await onSave()
    actionStore.nowRightSiderTabIndex = 2;
    await Start()
    nowRunningStatus.value = 1
  } catch(err) {
    message.error('运行失败:' + err)
  }

};

const stopAction = async () => {
  nowRunningStatus.value = 2
  await Stop()
  nowRunningStatus.value = 0
}

const openSetting = () => {
  GetAppSettings()
    .then((res) => {
      configSettings.value = res;
      addModalVisible.value = true;
    })
    .catch((err) => {
      console.error(err);
      message.error("获取应用设置失败");
    });
};

const handleAddApp = async () => {
  try {
    await SyncAppSettings(configSettings.value);
    message.success("修改应用成功");
    addModalVisible.value = false;
  } catch (err) {
    console.error(err);
    message.error("修改应用失败");
  }
};

const onSave = async () => {
  try {
    console.log("actions:", actionStore.actions);
    // 调用 wails 后端保存
    await SyncActions(actionStore.actions);
    console.log("step 1");
    
    
    await SaveApp();

    message.success("保存成功", 1);
  } catch (err) {
    console.error(err);
    message.error("保存失败");
  }
};
</script>

<style scoped>
.header-toolbar {
  display: flex;
  align-items: center;
}
.right {
  margin-left: auto;
}
.btn {
  margin: 10px;
}
</style>
