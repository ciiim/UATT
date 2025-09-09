<template>
    <div class="container">
    <Transition name="fade" mode="out-in">
      <component :is="tabs[tabsIndex]" :key="tabsIndex" :moduleLibrary="moduleLibrary"></component>
    </Transition>
    <a-segmented v-model:value="value" block :options="data" @change="switchTab" :style="{marginTop: 'auto', borderRadius:'0%'}" size="middle" />
    </div>
</template>

<script setup lang="ts">
import { reactive, ref, Transition, defineProps } from 'vue';
import ModuleList from './right_sider/ModuleList.vue';
import PropPage from './right_sider/PropPage.vue';
const data = reactive(["组件库","属性"]);
const tabs = [ModuleList, PropPage];
const value = ref(data[0]);
const tabsIndex = ref(0);

defineProps<{ moduleLibrary: any[] }>();

const switchTab = (v : string) => {
  let res = data.findIndex((e) => {return e === v});
  
  if (res != -1) {
    tabsIndex.value = res;
  }
};

</script>

<style lang="css">

.container {
    display: flex;
    width: 100%;
    height: 100%;
    flex-direction: column;
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