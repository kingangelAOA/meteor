<script setup name="task">
import { ElMessage, ElNotification } from "element-plus";
import { getTasks } from "@/api/pipeline";
import { parseTime } from "@/utils";
import Pagination from "@/components/Pagination/index.vue";
import { reactive, ref } from "vue";
import { useRoute, useRouter } from "vue-router";

const route = new useRoute();
const router = new useRouter();
const tableKey = ref(0);
const list = ref([]);
const total = ref(0);
const listLoading = ref(false);
const page = ref(0);
const listQuery = reactive({
  offset: 0,
  limit: 20,
  sortField: "id",
  sortType: -1,
  describe: "",
});
watch(
  () => page,
  (val) => {
    if (val < 1) {
      val = 0;
    }
    listQuery.offset = (val - 1) * listQuery.limit;
  }
);
const handlePaginationTasks = (data) => {
  const { page, limit } = data;
  listQuery.limit = limit;
  listQuery.offset = (page - 1) * limit;
  getList();
};
const getList = () => {
  listQuery.pipelineId = route.query.pipelineId;
  getTasks(listQuery)
    .then((response) => {
      list.value = response.data.items;
      total.value = response.data.total;
    })
    .finally(() => {
      listLoading.value = false;
    });
};
getList();
const sortChange = (data) => {
  const { prop, order } = data;
  if (order === "ascending") {
    listQuery.sortType = 1;
  } else {
    listQuery.sortType = -1;
  }
  listQuery.sortField = prop;
  handleFilter();
};
const getSortClass = (key) => {
  const sort = listQuery.sortType;
  if (sort === 1) {
    return "ascending";
  } else {
    return "descending";
  }
};

const goToDetail = (row) => {
  const { id: taskId, type } = row;
  if (type === "default") {
    router.push({
      name: "Flow",
      query: {
        from: "Task",
        pipelineId: route.query.pipelineId,
        taskId,
      },
    });
  } else {
    router.push({
      name: "Monitor",
      query: {
        pipelineId: route.query.pipelineId,
        taskId,
      },
    });
  }
};
</script>
<template>
  <div class="app-container">
    <el-table
      :key="tableKey"
      v-loading="listLoading"
      :data="list"
      size="small"
      border
      fit
      highlight-current-row
      style="width: 100%"
      @sort-change="sortChange"
    >
      <el-table-column
        label="index"
        prop="index"
        sortable="custom"
        align="center"
        width="120px"
      >
        <template #default="{ row }">
          <span>{{ row.index }}</span>
        </template>
      </el-table-column>
      <el-table-column label="任务耗时(毫秒)" align="center" prop="taskConsume">
      </el-table-column>
      <el-table-column label="状态" align="center" prop="status">
        <template #default="{ row }">
          <el-tag v-if="row.status" class="ml-2" type="success">成功</el-tag>
          <el-tag v-else class="ml-2" type="danger">失败</el-tag>
        </template>
      </el-table-column>
      <el-table-column
        label="创建时间"
        prop="createTime"
        align="center"
        sortable="custom"
        :class-name="getSortClass('createTime')"
      >
        <template #default="{ row }">
          <span>{{ row.createTime }}</span>
        </template>
      </el-table-column>
      <el-table-column
        label="Actions"
        align="center"
        width="200px"
        class-name="small-padding fixed-width"
      >
        <template #default="{ row }">
          <el-button-group>
            <el-button type="primary" size="small" @click="goToDetail(row)">
              详情
            </el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>
    <pagination
      v-show="total > 0"
      v-model:page="listQuery.page"
      v-model:limit="listQuery.limit"
      :total="total"
      @pagination="handlePaginationTasks"
    />
  </div>
</template>
