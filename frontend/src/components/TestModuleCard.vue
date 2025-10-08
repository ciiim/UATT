<template>
  <div class="action-card" :key="data.ActionUID">
    <!-- 第一行：图标 + 名称 + tag 列 -->
    <div class="action-header">
      <component :is="getIcon(data.ActionTypeID)" style="color: black; font-size: large;" />
      <span class="action-title">{{ data.Name }}</span>

      <div class="action-tags">
        <a-tag
          v-for="(tag, idx) in data.tags"
          :key="idx"
          color="blue"
          class="action-tag"
        >
          {{ tag.label }}<span v-if="tag.len != 0">({{ tag.len }})</span> 
        </a-tag>
      </div>
    </div>

    <!-- 分割线 -->
    <div class="action-divider"></div>

    <!-- 第二行：额外信息 + 操作按钮 -->
    <div class="action-footer">
      <span class="action-label">超时</span>
      <span class="action-value">
        {{ getTimeout(data) }}ms
      </span>
      <span class="action-label">状态</span>
      <span class="action-value">{{ data.status ?? '待定' }}</span>

      <div style="display: flex; margin-left: auto;">
        <a-button type="primary" size="small" style="margin-right: 10px;" @click="">
          <PlayCircleOutlined />
        </a-button>
        <a-button type="default" size="small" style="margin-right: 10px;" @click="">
          <FlagOutlined />
        </a-button>
        <a-button type="primary" danger size="small" @click="doDelete">
          <DeleteOutlined />
        </a-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, ref } from 'vue';
import { DeleteOutlined, ClockCircleOutlined, StopOutlined, LineOutlined, MoreOutlined, BranchesOutlined,FontColorsOutlined, EnterOutlined, PrinterOutlined,PlayCircleOutlined, FlagOutlined, SendOutlined, DownloadOutlined, QuestionOutlined } from '@ant-design/icons-vue';

import type { ConfigActionBase } from "../types/Action";
import { useActionStore } from '../stores/action_store';

const store = useActionStore();

const prop = defineProps<{
  data: any;
}>();

const getIcon = (actionId : number) => {
    let res = iconList.findIndex((e) => {return e.actionId == actionId})
    if(res != -1) {
        return iconList[res].icon;
    }
    return QuestionOutlined
}

const getTimeout = (action: ConfigActionBase) => {
  if (action.ActionTypeID === 1 || action.ActionTypeID === 2) {
    return action.TypeFeatureField?.TimeoutMs ?? '--'
  }
  return '--'
}

const doDelete = () => {
  
  store.actions.map((item,index) => {
    
    if (item.ActionUID == prop.data.ActionUID) {
      store.actions.splice(index,1)
    }
  })
}


const iconList = [
    {actionId: 1, icon: SendOutlined}, // send
    {actionId: 2, icon: DownloadOutlined}, // receive
    {actionId: 90, icon: PrinterOutlined}, // print
    {actionId: 91, icon: ClockCircleOutlined}, // delay
    {actionId: 20, icon: EnterOutlined}, // goto
    {actionId: 23, icon: FontColorsOutlined}, // declare
    {actionId: 24, icon: BranchesOutlined}, // if
    {actionId: 25, icon: MoreOutlined}, // else
    {actionId: 27, icon: LineOutlined}, // endblock
    {actionId: 31, icon: StopOutlined}, // stop


];

</script>

<style scoped>
.action-card {
  border: 1px solid #d9d9d9;
  border-radius: 8px;
  padding: 8px;
  background-color: #fff;
  min-width: 300px;
}

.action-header {
  display: flex;
  align-items: center;
  gap: 6px;
}

.action-title {
  font-weight: bold;
  font-size: medium;
  color: black;
}

.action-tags {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.action-tag {
  padding: 0 6px;
  height: 20px;
  line-height: 18px;
  font-size: 12px;
}

.action-divider {
  border-top: 1px solid #f0f0f0;
  margin: 6px 0;
}

.action-footer {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.action-label {
  font-size: 12px;
  color: rgba(0, 0, 0, 0.45);
}

.action-value {
  font-size: 12px;
  color: #000;
}
</style>