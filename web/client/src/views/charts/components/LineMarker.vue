<template>
  <div :id="id" :class="className" :style="{ height: height, width: width }" />
</template>

<script setup>
import * as echarts from "echarts";
import { reactive } from "vue";
import { v4 as uuidv4 } from "uuid";

const props = defineProps({
  className: {
    type: String,
    default: "chart",
  },
  width: {
    type: String,
    default: "200px",
  },
  height: {
    type: String,
    default: "200px",
  },
  data: {
    type: Array,
    required: true,
    default: [{ name: "", points: [] }],
  },
});
const state = reactive({
  chart: null,
});
const id = uuidv4();
const opt = reactive({
  title: {
    // text: "Stacked Line",
  },
  tooltip: {
    trigger: "axis",
  },
  legend: {
    // data: data[0],
  },
  grid: {
    left: "3%",
    right: "4%",
    bottom: "3%",
    containLabel: true,
  },
  toolbox: {
    feature: {
      // saveAsImage: {},
    },
  },
  xAxis: {
    type: "time",
  },
  yAxis: {
    type: "value",
  },
  series: [],
});
watch(
  () => [...props.data],
  (val) => {
    updateData(val);
  }
);

const updateData = (val) => {
  let ss = [];
  console.log(val);
  for (const item of val) {
    ss.push({
      symbol: "none",
      name: item.name,
      type: "line",
      // stack: "Total",
      data: item.data,
    });
  }
  opt.series = ss;
  state.chart.setOption(opt);
};

onMounted(() => {
  initChart();
  updateData(props.data);
});
onBeforeUnmount(() => {
  if (!state.chart) {
    return;
  }
  state.chart.dispose(document.getElementById(id));
  state.chart = null;
});
const initChart = () => {
  state.chart = markRaw(echarts.init(document.getElementById(id)));
  state.chart.setOption(opt);
};
// let helloFunc = () => {
//   console.log("helloFunc");
// };
//导出给refs使用
// defineExpose({ helloFunc });
//导出属性到页面中使用
// let {levelList} = toRefs(state);
</script>

<style scoped lang="scss"></style>
