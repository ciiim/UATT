<template>
	<a-float-button-group shape="square" style="right: 24px">
    	<a-float-button class="float-btn-half" description="Can vas" @click="() => {changeSystemMode(SystemMode.CanvasMode)}" />
		<a-float-button class="float-btn-half" description="View" @click="() => {changeSystemMode(SystemMode.CanvasViewMode)}" />
    	<a-float-button class="float-btn-half" description="Edit" @click="() => {changeSystemMode(SystemMode.EditMode)}" />
  	</a-float-button-group>
	<a-layout style="height: 100vh; overflow: hidden;">
		<component :is="nowSystemModeFrame" :rightSiderWidth="rightSiderWidth"></component>
	</a-layout>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { SystemMode } from '../stores/global';
import EditMode from './system_mode/edit_mode/EditMode.vue';
import CanvasMode from './system_mode/canvas_mode/CanvasMode.vue';


const rightSiderWidth = ref(300)



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

const nowSystemModeFrame = ref(EditMode)

const changeSystemMode = (mode : SystemMode) => {
	console.log('switch to:', mode);
	
	switch(mode) {
		case SystemMode.EditMode:
			nowSystemModeFrame.value = EditMode
		case SystemMode.CanvasMode:
			nowSystemModeFrame.value = CanvasMode
		case SystemMode.CanvasViewMode:
	}
}

onMounted(() => {
	updateRightSiderWidth()
	window.addEventListener('resize', updateRightSiderWidth)
})

onUnmounted(() => {
	window.removeEventListener('resize', updateRightSiderWidth)
})

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

.float-btn-half {
  opacity: 0.3;
  transition: opacity 0.2s;
}

/* 悬停时完全不透明 */
.float-btn-half:hover {
  opacity: 1;
}
</style>