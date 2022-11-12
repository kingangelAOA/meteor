<script setup>
import LineMarker from "@/views/charts/components/LineMarker.vue";
import {
  getRunningMonitorInfo,
  stopTask,
  resetQPS,
  resetUsers,
  getTaskPipelineUsers,
  getTaskStatus,
  getTaskWorks,
} from "@/api/task";
import { reactive } from "vue";
import { useRoute } from "vue-router";
import { v4 as uuidv4 } from "uuid";

const route = new useRoute();
const taskId = route.query.taskId;
const lastTimeList = [1, 5, 10, 20, 30, 60];
const search = reactive({
  interval: [0, 0],
  lastTime: 5,
});

const getKey = () => {
  return uuidv4();
};

watch(
  () => [...search.interval],
  (val) => {
    search.lastTime = 0;
    getMonitorData();
  }
);

const LineData = ref([]);
const status = ref("");
const runningWorks = ref(0);
const getMonitorData = () => {
  getTaskStatus({ id: taskId }).then((res) => {
    status.value = res.data;
    if (status.value === "running") {
      getTaskWorks({ id: taskId }).then((res) => {
        runningWorks.value = res.data;
      });
    }
  });
  getRunningMonitorInfo({
    id: taskId,
    startTime: search.interval[0] / 1000,
    endTime: search.interval[1] / 1000,
    lastTime: search.lastTime * 60,
  }).then((res) => {
    const result = [];
    for (const item of res.data) {
      if (Object.keys(item.nodePerformance) === 0) {
        return;
      }
      let nodeData = {};
      const { RTAverage, ErrRate, RT90, RT95, RT99, QPS } = item.nodePerformance;
      const rt99D = RT99;
      const rt95D = RT95;
      const rt90D = RT90;
      const averageRTD = RTAverage;
      const qpsD = QPS;
      const errRateD = ErrRate;
      nodeData.name = item.name;
      nodeData.points = [
        [
          { name: "rt99", data: rt99D },
          { name: "rt95", data: rt95D },
          { name: "rt90", data: rt90D },
          { name: "averageRT", data: averageRTD },
        ],
        [{ name: "qps", data: qpsD }],
        [{ name: "err", data: errRateD }],
      ];
      result.push(nodeData);
    }
    LineData.value = result;
  });
};
getMonitorData();
let timer = null;
onMounted(() => {
  timer = setInterval(async () => {
    getMonitorData();
  }, 3000);
});
onUnmounted(() => clearInterval(timer));

watch(status, (val) => {
  if (val !== "running") {
    clearInterval(timer);
  } else {
    getUsers();
  }
});

const clickStatus = ref(false);

const handleStop = () => {
  clickStatus.value = true;
  stopTask({ id: taskId }).finally(() => {
    clickStatus.value = false;
  });
};

const users = ref(0);
const getUsers = () => {
  getTaskPipelineUsers({ id: taskId }).then((res) => {
    users.value = res.data;
  });
};
const handleResetUsers = () => {
  clickStatus.value = true;
  resetUsers({ id: taskId, users: users.value }).finally(() => {
    clickStatus.value = false;
  });
};

const qps = ref(0);
const handleResetQPS = () => {
  clickStatus.value = true;
  resetQPS({ id: taskId, qps: qps.value }).finally(() => {
    clickStatus.value = false;
  });
};
</script>
<template>
  <div>
    <el-form :inline="true" :model="search" class="demo-form-inline">
      <el-form-item label="时间段:">
        <el-date-picker
          v-model="search.interval"
          type="datetimerange"
          start-placeholder="Start Date"
          end-placeholder="End Date"
          format="YYYY/MM/DD HH:mm:ss"
        />
      </el-form-item>
      <el-form-item label="最近(分钟):">
        <el-select v-model="search.lastTime" placeholder="Select" size="small">
          <el-option
            v-for="item in lastTimeList"
            :key="item"
            :label="item + '分钟'"
            :value="item"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="协程数量:">
        <el-input-number
          v-model="users"
          placeholder="Please input"
          class="input-with-select"
        >
        </el-input-number>
        <el-button type="primary" :loading="clickStatus" @click="handleResetUsers"
          >设置</el-button
        >
      </el-form-item>
      <el-form-item label="qps:">
        <el-input-number
          v-model="qps"
          placeholder="Please input"
          class="input-with-select"
        >
        </el-input-number>
        <el-button type="primary" :loading="clickStatus" @click="handleResetQPS"
          >设置</el-button
        >
      </el-form-item>
      <el-form-item label="停止任务:">
        <el-button type="primary" :loading="clickStatus" @click="handleStop"
          >停止</el-button
        >
      </el-form-item>
      <el-form-item label="占用协程数量:">
        <span>{{ runningWorks }}</span>
      </el-form-item>
    </el-form>
    <el-row v-for="(item, index) in LineData" :key="index" :gutter="32">
      <el-divider content-position="center">{{ item.name }}</el-divider>
      <el-col v-for="(data, index) in item.points" :key="index" :xs="24" :sm="24" :lg="8">
        <div class="dashboard-editor-container">
          <div class="chart-wrapper">
            <LineMarker :data="data" height="200px" width="100%" />
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>
<style lang="scss" scoped>
.dashboard-editor-container {
  padding: 5px;
  background-color: rgb(240, 242, 245);
  position: relative;

  .github-corner {
    position: absolute;
    top: 0;
    border: 0;
    right: 0;
  }

  .chart-wrapper {
    background: #fff;
    // padding: 16px 16px 0;
    // margin-bottom: 32px;
  }
}

@media (max-width: 1024px) {
  .chart-wrapper {
    padding: 8px;
  }
}
</style>
