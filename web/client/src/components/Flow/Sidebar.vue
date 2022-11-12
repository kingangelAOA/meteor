<script setup>
import { Calendar, Search } from "@element-plus/icons-vue";
import { getBindingPluginNodes } from "@/api/node.js";
import { reactive } from "vue-demi";
const search = ref("");
const onDragStart = (event, node) => {
  if (event.dataTransfer) {
    console.log(event.clientX, event.clientY);
    event.dataTransfer.setData("application/vueflow", node);
    event.dataTransfer.effectAllowed = "move";
  }
};
let nodeList = ref([]);
watch(search, (newValue) => {
  getBindingPluginNodes({ search: `name:${newValue}` }).then((res) => {
    nodeList.value = res.data;
  });
});
</script>

<template>
  <div>
    <el-input
      v-model="search"
      placeholder="Type something"
      :prefix-icon="Search"
      style="padding-bottom: 5px"
    />
    <aside style="background: #ffffff">
      <div class="nodes">
        <div
          v-for="(item, index) in nodeList"
          :key="index"
          class="ef-node-menu-li"
          :draggable="true"
          @dragstart="onDragStart($event, JSON.stringify(item))"
        >
          {{ item.name }}
        </div>
      </div>
    </aside>
  </div>
</template>

<style scoped lang="scss">
.box-card {
  width: auto;
  padding-bottom: 5px;
}
.ef-node-menu-li {
  color: #000000;
  width: 150px;
  border: 1px dashed #e0e3e7;
  padding: 10px;
  border-radius: 5px;
  margin: auto;
  text-align: center;
}
</style>
