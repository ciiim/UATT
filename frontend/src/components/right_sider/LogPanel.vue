<template>
  <div class="log-panel">
    <div class="toolbar">
      <a-button type="primary" size="small" @click="saveLogs">保存日志</a-button>
      <a-button type="default" size="small" @click="clearLogs">清空</a-button>
    </div>

    <a-list
      size="small"
      bordered
      :data-source="store.logs"
      style="height: calc(100% - 40px); overflow-y: auto;"
    >
    <template #renderItem="{ item }">
        <a-list-item>
          <code>{{ item }}</code> 
        </a-list-item>
      </template>
</a-list>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { Button, List, message } from "ant-design-vue";
import { EventsOn } from "../../../wailsjs/runtime/runtime"; // 假设你用 Wails 的事件机制
import { useActionStore } from "../../stores/action_store";

const store = useActionStore();

onMounted(() => {
  EventsOn("runtime-log", (msg: string) => {
    store.logs.push(msg);
  });
});

const clearLogs = () => {
  store.logs = [];
};

const saveLogs = () => {
  const blob = new Blob([store.logs.join("\n")], { type: "text/plain;charset=utf-8" });

  // 利用浏览器 API 保存
  const url = URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.download = `logs_${new Date().toISOString()}.txt`;
  link.click();
  URL.revokeObjectURL(url);

  message.success("日志已保存");
};
</script>

<style scoped>
.log-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.toolbar {
  display: flex;
  gap: 8px;
  padding: 4px;
  border-bottom: 1px solid #eee;
  background-color: #fff;
}
</style>