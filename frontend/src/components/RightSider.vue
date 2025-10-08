<template>
  <div class="right-sider">
    <div class="content-area">
      <Transition name="fade" mode="out-in">
        <component
          :is="tabs[tabsIndex]"
          :key="tabsIndex"
          :actionLibrary="actionLibrary"
        />
      </Transition>
    </div>
    <div class="tab-area">
      <a-segmented
        v-model:value="value"
        block
        :options="data"
        @change="switchTab"
        size="middle"
      />
    </div>
  </div>
</template>



<script setup lang="ts">
import { reactive, ref, Transition, defineProps } from "vue";
import type { ConfigActionBase } from "../types/Action";
import ModuleList from "./right_sider/ModuleList.vue";
import PropPage from "./right_sider/PropPage.vue";
const data = reactive(["组件库", "属性"]);
const tabs = [ModuleList, PropPage];
const value = ref(data[0]);
const tabsIndex = ref(0);

defineProps<{ 
  actionLibrary: any[],
 }>();

const switchTab = (v: string) => {
  let res = data.findIndex((e) => {
    return e === v;
  });

  if (res != -1) {
    tabsIndex.value = res;
  }
};
</script>

<style lang="css">
.right-sider {
  display: flex;
  flex-direction: column;
  height: 100%; /* 撑满右侧栏 */
}

/* 中间内容区域：自动滚动 */
.content-area {
  flex: 1;               /* 占满剩余空间 */
  overflow-y: auto;      /* 垂直滚动 */
  padding: 8px;
  background-color: whitesmoke;
}

/* 底部 Tab 区域：固定高度 */
.tab-area {
  flex-shrink: 0;
  border-top: 1px solid #eee;
  padding: 4px;
  background: #fff;
}

.fade-enter-active,
.fade-leave-active {
  transition: all 0.1s ease-in-out;
}

.fade-enter-from,
.fade-leave-to {
  transform: translateX(300px);
  opacity: 0.3;
}
</style>
