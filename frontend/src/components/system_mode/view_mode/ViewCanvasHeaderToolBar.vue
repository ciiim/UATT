<template>
  <div class="header-toolbar">
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
      <a-button class="btn" type="primary"@click="exportButton">导出生成Viewer</a-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { message } from "ant-design-vue";
import { useActionStore } from "../../../stores/action_store";
import {
  GetAllSerial,
  SelectSerialCom,
  OpenSerialPort,
  CloseSerialPort,
  ExportViewer,
} from "../../../../wailsjs/go/bsd_testtool/Manager";
import { bsd_testtool } from "../../../../wailsjs/go/models";

const store = useActionStore();

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

  
  await GetSerialList()
  if (serialList.value.length > 0) {
      selectedSerial.value = serialList.value[0];
      // 默认选中第一个串口
      await SelectSerialCom(selectedSerial.value);
  }

  await closeSerialPort()

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
    await SelectSerialCom(port)
    //message.success(`已选择串口: ${port}`);
  } catch (err) {
    console.error(err);
    message.error("选择串口失败");
  }
};


const exportButton = async () => {
  try {
    let appNum = await ExportViewer()
    message.info(['生成成功, 收集了', appNum, '个App'])
  } catch(err) {
    console.log('Close serial port err:', err);
  }
}

const closeSerialPort = async () => {
  nowRunningStatus.value = 2
  try {
    await CloseSerialPort()
  } catch(err) {
    console.log('Close serial port err:', err);
  }
  nowRunningStatus.value = 0
}

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
