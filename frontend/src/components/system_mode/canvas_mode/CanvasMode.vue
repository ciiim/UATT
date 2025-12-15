<template><a-layout-sider
    :style="siderStyle"
    breakpoint="lg"
    collapsed-width="0"
    v-model:collapsed="leftCollapsed"
    :trigger="null"
    collapsible
  >
    <LeftCanvasList @canvas-loaded="notifyGetCanvasContent"></LeftCanvasList>
  </a-layout-sider>
  <a-layout>
    <a-layout-header :style="headerStyle">
      <!-- 左菜单折叠按钮 -->
      <menu-unfold-outlined
        v-if="leftCollapsed"
        class="left-trigger"
        @click="() => (leftCollapsed = !leftCollapsed)"
      />
      <menu-fold-outlined
        v-else
        class="left-trigger"
        @click="() => (leftCollapsed = !leftCollapsed)"
      />

      <CanvasHeaderToolBar></CanvasHeaderToolBar>

      <!-- 右菜单折叠按钮 -->
      <menu-fold-outlined
        v-if="rightCollapsed"
        class="right-trigger"
        @click="() => (rightCollapsed = !rightCollapsed)"
      />
      <menu-unfold-outlined
        v-else
        class="right-trigger"
        @click="() => (rightCollapsed = !rightCollapsed)"
      />
    </a-layout-header>

    <a-layout-content :style="contentStyle">
      <MainCanvas
        :need-get-canvas-list="needGetCanvas"
        :canvas-tool-Library="canvasToolLibrary"
      ></MainCanvas>
    </a-layout-content>
  </a-layout>

  <a-layout-sider
    :width="rightSiderWidth"
    :style="siderStyle"
    breakpoint="lg"
    collapsed-width="0"
    :trigger="null"
    v-model:collapsed="rightCollapsed"
    collapsible
    reverseArrow
  >
    <RightToolBox :canvas-tool-library="canvasToolLibrary"></RightToolBox>
  </a-layout-sider>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import type { CSSProperties } from "vue";
import LeftCanvasList from "./LeftCanvasList.vue";
import CanvasHeaderToolBar from "./CanvasHeaderToolBar.vue";
import MainCanvas from "./MainCanvas.vue";
import {
  MenuUnfoldOutlined,
  MenuFoldOutlined,
  FolderViewOutlined,
  EditOutlined,
  FundProjectionScreenOutlined,
} from "@ant-design/icons-vue";
import RightToolBox from "./RightToolBox.vue";

defineProps<{rightSiderWidth : number}>();

const headerStyle: CSSProperties = {
  display: "flex",
  alignItems: "center",
  color: "#fff",
  height: "64px",
  paddingInline: "10px",
  lineHeight: "64px",
  backgroundColor: "#7dbcea",
};

const siderStyle: CSSProperties = {
  textAlign: "center",
  color: "#fff",
  backgroundColor: "#3ba0e9",
};

const contentStyle: CSSProperties = {
  textAlign: "center",
  color: "#fff",
  backgroundColor: "whitegray",
  overflow: "auto",
  flex: 1,
};

const leftCollapsed = ref<boolean>(false);
const rightCollapsed = ref<boolean>(false);

const needGetCanvas = ref<boolean>(false);

const notifyGetCanvasContent = () => {
  needGetCanvas.value = !needGetCanvas.value
}

const canvasToolLibrary = ref([
])
</script>