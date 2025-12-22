<template>
  <a-layout>
    <a-layout-header :style="headerStyle">

      <ViewCanvasHeaderToolBar></ViewCanvasHeaderToolBar>
    </a-layout-header>

    <a-layout-content :style="contentStyle">
      <MainViewCanvas
      ></MainViewCanvas>
    </a-layout-content>
  </a-layout>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import type { CSSProperties } from "vue";
import ViewCanvasHeaderToolBar from "./ViewCanvasHeaderToolBar.vue";
import MainViewCanvas from "./MainViewCanvas.vue";

import { useActionStore } from "../../../stores/action_store";

defineProps<{rightSiderWidth : number}>();

const store = useActionStore();

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


</script>