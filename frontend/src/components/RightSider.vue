<template>
    <div class="container">
    <component :is="tabs[tabsIndex]"></component>
    <a-segmented v-model:value="value" block :options="data" @change="switchTab" :style="{marginTop: 'auto', borderRadius:'0%'}" size="middle" />
    </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import ModuleList from './right_sider/ModuleList.vue';
import PropPage from './right_sider/PropPage.vue';
const data = reactive(["组件库","属性"]);
const tabs = [ModuleList, PropPage];
const value = ref(data[0]);
const tabsIndex = ref(0);

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


</style>