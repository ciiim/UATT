<template>
  <div class="module-card">
    <!-- 第一行：模块名 + 标签 -->
    <div class="module-header">
      <component :is="getIcon(data.modUid)" style="color: black; font-size: large;"/>
      <span class="module-title">{{ data.name }}</span>
      <div class="module-tags">
        <a-tag
          v-for="(tag, index) in data.tags"
          :key="index"
          color="blue"
          class="module-tag"
        >
          {{ tag.label }}({{ tag.len }})
        </a-tag>
      </div>
    </div>

    <!-- 分割线 -->
    <div class="module-divider"></div>

    <!-- 第二行：展示项 + 按钮 -->
    <div class="module-footer">
      <span class="module-label">超时</span>
      <span class="module-value">{{ data.timeout }}ms</span>
      <span class="module-label">状态</span>
      <span class="module-value">{{ data.status }}</span>
      <div style="display: flex;margin-left: auto;">
      <a-button type="primary" size="small" style="margin-right: 10px;"><PlayCircleOutlined /></a-button>
      <a-button type="default" size="small"><FlagOutlined /></a-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps } from 'vue';
import { PlayCircleOutlined, FlagOutlined, SendOutlined, DownloadOutlined, QuestionOutlined } from '@ant-design/icons-vue';

interface Tag {
  label: string;
  len: number;
}
interface ModuleCardData {
  name: string;
  modUid: number;
  tags: Tag[];
  timeout: number;
  status: string;
}

defineProps<{
  data: ModuleCardData;
}>();

const getIcon = (uid : number) => {
    let res = iconList.findIndex((e) => {return e.modUid == uid})
    //console.log('uid' + uid + " res:" + res)
    if(res != -1) {
        return iconList[res].icon;
    }
    return QuestionOutlined
} 

const iconList = [
    {modUid: 1, icon: SendOutlined},
    {modUid: 2, icon: DownloadOutlined}
];

</script>

<style scoped>
.module-card {
  border: 1px solid #d9d9d9;
  border-radius: 8px;
  padding: 8px;
  background-color: #fff;
  min-width: 300px;
}

.module-header {
  display: flex;
  align-items: center;
  gap: 6px;
}

.module-title {
  font-weight: bold;
  font-size: medium;
  color: black;
}

.module-tags {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.module-tag {
  padding: 0 6px;
  height: 20px;
  line-height: 18px;
  font-size: 12px;
}

.module-divider {
  border-top: 1px solid #f0f0f0;
  margin: 6px 0;
}

.module-footer {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.module-label {
  font-size: 12px;
  color: rgba(0, 0, 0, 0.45);
}

.module-value {
  font-size: 12px;
  color: #000;
}
</style>