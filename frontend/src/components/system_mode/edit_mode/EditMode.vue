<template>
  <a-layout-sider
    :style="siderStyle"
    breakpoint="lg"
    collapsed-width="0"
    v-model:collapsed="leftCollapsed"
    :trigger="null"
    collapsible
  >
    <LeftSiderMenu @app-loaded="notifyGetActionList"></LeftSiderMenu>
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

      <HeaderToolBar></HeaderToolBar>

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
      <MainContent
        :need-get-action-list="needGetActionList"
        :actionLibrary="actionLibrary"
      ></MainContent>
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
    <RightSider :action-library="actionLibrary"></RightSider>
  </a-layout-sider>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import type { CSSProperties } from "vue";
import { message, FloatButton } from "ant-design-vue";
import { SaveApp, SyncActions } from "../../../../wailsjs/go/bsd_testtool/Manager";
import MainContent from "./MainContent.vue";
import {
  MenuUnfoldOutlined,
  MenuFoldOutlined,
  FolderViewOutlined,
  EditOutlined,
  FundProjectionScreenOutlined,
} from "@ant-design/icons-vue";
import HeaderToolBar from "./HeaderToolBar.vue";
import RightSider from "../../RightSider.vue";
import LeftSiderMenu from "./LeftSiderMenu.vue";
import { EventsOn, EventsOff } from "../../../../wailsjs/runtime/runtime";
import { useActionStore } from "../../../stores/action_store";
import { SystemMode } from "../../../stores/global";

const leftCollapsed = ref<boolean>(false);
const rightCollapsed = ref<boolean>(false);

const store = useActionStore();

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

const actionLibrary = ref([
  {
    name: "发送",
    actionID: 1,
    actionType: "IO",
    feat: { TimeoutMs: 1000, Modules: [] },
  },
  {
    name: "接收",
    actionID: 2,
    actionType: "IO",
    feat: { TimeoutMs: 1000, Modules: [] },
  },
  { name: "Print", actionID: 90, actionType: "Debug", feat: { PrintFmt: "" } },
  { name: "Delay", actionID: 91, actionType: "Debug", feat: { DelayMs: 0 } },
  { name: "Show", actionID: 92, actionType: "Debug", feat: { OutputIdx: 0, FmtStr: '' } },
  {
    name: "Label",
    actionID: 28,
    actionType: "Control",
    feat: { LabelName: "" },
  },
  { name: "Goto", actionID: 29, actionType: "Control", feat: { Label: "" } },
  {
    name: "Declare",
    actionID: 23,
    actionType: "Control",
    feat: {
      VarName: "",
      VarType: "",
      VarNumberValue: 0,
      VarStringValue: "",
      VarByteArrayValue: [],
    },
  },
  {
    name: "Assign",
    actionID: 32,
    actionType: "Control",
    feat: { AssignTargetVar: "", Expression: "" },
  },
  { name: "If", actionID: 24, actionType: "Control", feat: { Condition: "" } },
  { name: "Else", actionID: 25, actionType: "Control", feat: {} },
  {
    name: "For",
    actionID: 26,
    actionType: "Control",
    feat: { UseVar: "", EnterCondition: "", VarOp: "" },
  },
  { name: "EndBlock", actionID: 27, actionType: "Control", feat: {} },
]);

const needGetActionList = ref<boolean>(false);
const notifyGetActionList = () => {
  needGetActionList.value = !needGetActionList.value;
};

const handleKeyDown = async (e: KeyboardEvent) => {
  if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === "s") {
    e.preventDefault(); // 阻止浏览器保存网页
    try {
      await SyncActions(store.actions);
      await SaveApp();
      message.success("保存成功", 1);
    } catch (err) {
      if (err != "could not found app") {
        message.error("保存失败," + err);
      }
    }
  }
};

onMounted(() => {
  window.addEventListener("keydown", handleKeyDown);
  EventsOn("runtime-log", (msg: string) => {
    store.logs.push(msg);
  });
});

onUnmounted(() => {
  window.removeEventListener("keydown", handleKeyDown);
  EventsOff("runtime-log");
});
</script>
