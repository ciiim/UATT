<template>

	<a-layout style="height: 100vh; overflow: hidden;">
		<a-layout-sider :style="siderStyle" breakpoint="lg" collapsed-width="0" v-model:collapsed="leftCollapsed" :trigger="null" collapsible>
			<LeftSiderMenu
			@app-loaded="notifyGetActionList"
			></LeftSiderMenu>
		</a-layout-sider>
		<a-layout>
			<a-layout-header :style="headerStyle">
				<!-- 左菜单折叠按钮 -->
				<menu-unfold-outlined v-if="leftCollapsed" class="left-trigger" @click="() => (leftCollapsed = !leftCollapsed)" />
				<menu-fold-outlined v-else class="left-trigger" @click="() => (leftCollapsed = !leftCollapsed)" />


				<HeaderToolBar></HeaderToolBar>

				<!-- 右菜单折叠按钮 -->
				<menu-fold-outlined v-if="rightCollapsed" class="right-trigger" @click="() => (rightCollapsed = !rightCollapsed)" />
				<menu-unfold-outlined v-else class="right-trigger" @click="() => (rightCollapsed = !rightCollapsed)" />

			</a-layout-header>

			<a-layout-content :style="contentStyle">
				<MainContent 
				:need-get-action-list="needGetActionList"
      			:actionLibrary="actionLibrary"></MainContent>
			</a-layout-content>
		</a-layout>

		<a-layout-sider :width="rightSiderWidth" :style="siderStyle" breakpoint="lg" collapsed-width="0" :trigger="null" v-model:collapsed="rightCollapsed" collapsible
			reverseArrow>
			<RightSider 
			:action-library="actionLibrary"
			></RightSider>
		</a-layout-sider>
	</a-layout>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import type { CSSProperties } from 'vue'
import { message } from 'ant-design-vue'
import { SaveApp, SyncActions } from '../../wailsjs/go/bsd_testtool/Manager'
import MainContent from './MainContent.vue';
import {
	MenuUnfoldOutlined,
	MenuFoldOutlined,
} from '@ant-design/icons-vue';
import HeaderToolBar from './HeaderToolBar.vue';
import RightSider from './RightSider.vue';
import LeftSiderMenu from './LeftSiderMenu.vue';
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';
import { useActionStore } from '../stores/action_store';

const leftCollapsed = ref<boolean>(false);
const rightCollapsed = ref<boolean>(false);

const rightSiderWidth = ref(300)

const store = useActionStore();

const updateRightSiderWidth = () => {
  const winWidth = window.innerWidth
  if (winWidth >= 1600) {
    rightSiderWidth.value = 400
  } else if (winWidth >= 1200) {
    rightSiderWidth.value = 320
  } else if (winWidth >= 800) {
    rightSiderWidth.value = 270
  } else {
    rightSiderWidth.value = 200
  }
}

const needGetActionList = ref<boolean>(false);
const notifyGetActionList = () => {
	needGetActionList.value = !needGetActionList.value
}

const siderStyle: CSSProperties = {
	textAlign: 'center',
	color: '#fff',
	backgroundColor: '#3ba0e9',
}

const handleKeyDown = async (e: KeyboardEvent) => {
  if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 's') {
    e.preventDefault() // 阻止浏览器保存网页
    try {
      await SyncActions(store.actions)
      await SaveApp()
      message.success('保存成功', 1)
    } catch (err) {
		if (err != 'could not found app') {
			message.error('保存失败,' + err)
		}
    }
  }
}

onMounted(() => {
	updateRightSiderWidth()
  	window.addEventListener('resize', updateRightSiderWidth)
  	window.addEventListener('keydown', handleKeyDown)
	EventsOn("runtime-log", (msg: string) => {
    	store.logs.push(msg);
  	});
})

onUnmounted(() => {
	window.removeEventListener('resize', updateRightSiderWidth)
  	window.removeEventListener('keydown', handleKeyDown)
	EventsOff("runtime-log");
})

const headerStyle: CSSProperties = {
	display: "flex",
	alignItems: 'center',
	color: '#fff',
	height: '64px',
	paddingInline: '10px',
	lineHeight: '64px',
	backgroundColor: '#7dbcea',
}



const contentStyle: CSSProperties = {
	textAlign: 'center',
	color: '#fff',
	backgroundColor: 'whitegray',
	overflow: 'auto',
	flex: 1
}


const actionLibrary = ref([
  {name: '发送', actionID: 1, actionType: 'IO', feat : {TimeoutMs: 1000, Modules: []}},
  {name: '接收', actionID: 2, actionType: 'IO', feat : {TimeoutMs: 1000, Modules: []}},
  {name: 'PRINT', actionID: 90, actionType: 'Debug', feat : {PrintFmt: ""}},
  {name: 'Delay', actionID: 91, actionType: 'Debug', feat : {DelayMs: 0}},
  {name: 'Label', actionID: 28, actionType: 'Control', feat : {LabelName: ''}},
  {name: 'Goto', actionID: 29, actionType: 'Control', feat : {Label: ''}},
  {name: 'Declare', actionID: 23, actionType: 'Control', feat : {VarName: '', VarType: '', VarNumberValue: 0, VarStringValue: '', VarByteArrayValue: []}},
  {name: 'IF', actionID: 24, actionType: 'Control', feat : {Condition: ""}},
  {name: 'ELSE', actionID: 25, actionType: 'Control', feat : {}},
  {name: 'EndBlock', actionID: 27, actionType: 'Control', feat : {}},
])

const nowApp = ref('')

</script>

<style>
html,
body,
#app {
	height: 100%;
	margin: 0;
}

.left-trigger, .right-trigger {
	font-size: 20px;
	padding: 0 5px;
	cursor: pointer;
	transition: color 0.3s;
}

.left-trigger :hover, .right-trigger :hover {
	color: #1890ff;
}

.right-trigger {
	margin-left: auto;
}
</style>