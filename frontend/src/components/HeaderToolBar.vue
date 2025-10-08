<template>
  <div class="header-toolbar">
    <div class="left">
      <!-- 左边内容 -->
    </div>
    <div class="right">
      <a-button class="btn" type="primary" @click="startAction">运行</a-button>
      <a-button class="btn" type="primary" @click="onSave">保存</a-button>
      <a-button class="btn" type="primary" @click="openSetting">设置</a-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { message } from 'ant-design-vue'
import { useActionStore } from '../stores/action_store'
import { SaveApp, SyncActions } from '../../wailsjs/go/bsd_testtool/Manager'


const actionStore = useActionStore()

const startAction = async () => {

}

const openSetting = async () => {

}

const onSave = async () => {
  try {
    // 调用 wails 后端保存
    await SyncActions(actionStore.actions)
    await SaveApp()

    message.success('保存成功')
  } catch (err) {
    console.error(err)
    message.error('保存失败')
  }
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