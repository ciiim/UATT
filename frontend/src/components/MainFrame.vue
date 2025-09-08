<template>

	<a-layout style="height: 100vh; overflow: hidden;">
		<a-layout-sider :style="siderStyle" breakpoint="lg" collapsed-width="0" v-model:collapsed="leftCollapsed" :trigger="null" collapsible>Sider</a-layout-sider>
		<a-layout>
			<a-layout-header :style="headerStyle">
				<!-- 左菜单折叠按钮 -->
				<menu-unfold-outlined v-if="leftCollapsed" class="left-trigger" @click="() => (leftCollapsed = !leftCollapsed)" />
				<menu-fold-outlined v-else class="left-trigger" @click="() => (leftCollapsed = !leftCollapsed)" />

				<!-- 右菜单折叠按钮 -->
				<menu-fold-outlined v-if="rightCollapsed" class="right-trigger" @click="() => (rightCollapsed = !rightCollapsed)" />
				<menu-unfold-outlined v-else class="right-trigger" @click="() => (rightCollapsed = !rightCollapsed)" />

				

			</a-layout-header>

			<a-layout-content :style="contentStyle">
				<MainContent></MainContent>
			</a-layout-content>
		</a-layout>

		<a-layout-sider :style="siderStyle" breakpoint="lg" collapsed-width="0" :trigger="null" v-model:collapsed="rightCollapsed" collapsible
			reverseArrow>Sider</a-layout-sider>
	</a-layout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { CSSProperties } from 'vue'
import MainContent from './MainContent.vue';
import {
	MenuUnfoldOutlined,
	MenuFoldOutlined,
} from '@ant-design/icons-vue';

const leftCollapsed = ref<boolean>(false);
const rightCollapsed = ref<boolean>(false);

const siderStyle: CSSProperties = {
	textAlign: 'center',
	color: '#fff',
	backgroundColor: '#3ba0e9'
}


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
	backgroundColor: '#108ee9',
	overflow: 'auto',
	flex: 1
}

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
	flex: 1;
	text-align: right;
}
</style>