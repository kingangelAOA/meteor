<script setup name="pipeline">
import {
  getPiplines,
  updatePipeline,
  createPipeline,
  deletePipeline,
  getTypes,
  run,
} from "@/api/pipeline";
// import waves from '@/directive/waves' // waves directive
import { parseTime } from "@/utils";
import { ElMessage, ElNotification } from "element-plus";
import { ref, reactive, toRefs } from "vue";
import { useRouter } from "vue-router";
import Pagination from "@/components/Pagination/index.vue"; // secondary package based on el-pagination

const list = reactive({
  data: [],
});
const types = ref([]);
const router = useRouter();
const total = ref(0);
const listLoading = ref(true);
const page = ref(1);
const listQuery = reactive({
  offset: 0,
  limit: 20,
  sortField: "id",
  sortType: -1,
  search: "",
});
let temp = reactive({
  id: undefined,
  describe: "",
  type: "",
  config: {},
});
const dialogFormVisible = ref(false);
const dialogStatus = ref("");
const textMap = reactive({
  update: "Edit",
  create: "Create",
});
const rules = reactive({
  type: [{ required: false, message: "type is required", trigger: "change" }],
  timestamp: [
    { type: "date", required: true, message: "timestamp is required", trigger: "change" },
  ],
  title: [{ required: false, message: "title is required", trigger: "blur" }],
});
const paginationHandleGetList = (data) => {
  const { page, limit } = data;
  listQuery.limit = limit;
  listQuery.offset = (page - 1) * limit;
  getList();
};

const getPipelineTypes = () => {
  getTypes().then((res) => {
    types.value = res.data;
  });
};

getPipelineTypes();

const getList = () => {
  listLoading.value = true;
  getPiplines(listQuery)
    .then((response) => {
      list.data = response.data.items;
      total.value = response.data.total;
    })
    .finally(() => {
      listLoading.value = false;
    });
};
getList();
const handleFilter = () => {
  listQuery.search = `describe:${listQuery.filter}`;
  page.value = 1;
  getList();
};

const handleModifyStatus = (row, status) => {
  ElMessage({
    message: "操作成功",
    type: "success",
  });
  row.status = status;
};

const sortChange = (data) => {
  const { prop, order } = data;
  console.log(prop, order);
  if (order === "ascending") {
    listQuery.sortType = 1;
  } else {
    listQuery.sortType = -1;
  }
  listQuery.sortField = prop;
  handleFilter();
};

const resetTemp = () => {
  temp.id = undefined;
  temp.describe = "";
  temp.timestamp = new Date();
};

const handleCreate = () => {
  resetTemp();
  dialogStatus.value = "create";
  dialogFormVisible.value = true;
  nextTick(() => {
    dataForm.value?.clearValidate();
  });
};
const dataForm = ref(null);

const handleUpdate = (row) => {
  console.log(row);
  temp.id = row.id;
  temp.describe = row.describe; // copy obj
  temp.type = row.type;
  temp.config = row.config;
  dialogStatus.value = "update";
  dialogFormVisible.value = true;
  nextTick(() => {
    dataForm.value?.clearValidate();
  });
};
const updateData = () => {
  dataForm.value.validate((valid) => {
    if (valid) {
      updatePipeline(temp)
        .then(() => {
          getList();
        })
        .finally(() => {
          dialogFormVisible.value = false;
        });
    }
  });
};
const createData = () => {
  dataForm.value.validate((valid) => {
    if (valid) {
      createPipeline(temp)
        .then(() => {
          getList();
        })
        .finally(() => {
          dialogFormVisible.value = false;
        });
    }
  });
};
const handleDelete = (row, index) => {
  deletePipeline({ id: row.id })
    .then(() => {
      ElNotification({
        title: "成功",
        message: "删除成功",
        type: "success",
        duration: 2000,
      });
    })
    .finally(() => {
      getList();
    });
};

const formatJson = (filterVal) => {
  return list.map((v) =>
    filterVal.map((j) => {
      if (j === "timestamp") {
        return parseTime(v[j]);
      } else {
        return v[j];
      }
    })
  );
};

const getSortClass = (key) => {
  const sort = listQuery.sortType;
  if (sort === 1) {
    return "ascending";
  } else {
    return "descending";
  }
};

const handleTask = (row) => {
  router.push({
    name: "Task",
    query: {
      pipelineId: row.id,
    },
  });
};

const handleRun = (row) => {
  run({ id: row.id }).then((res) => {
    if (row.type === "default") {
      router.push({
        name: "Flow",
        query: {
          from: "PipelineRun",
          pipelineId: row.id,
          taskId: res.data,
        },
      });
    } else {
      router.push({
        name: "Monitor",
        query: {
          from: "Pipeline",
          pipelineId: row.id,
          taskId: res.data,
        },
      });
    }
  });
};

watch(
  () => page.value,
  (val) => {
    if (val < 1) {
      val = 0;
    }
    listQuery.offset = (val - 1) * listQuery.limit;
  },
  { immediate: true }
);
</script>
<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input
        v-model="listQuery.filter"
        placeholder="描述"
        style="width: 200px"
        class="filter-item"
        @keyup.enter="handleFilter"
      />
      <el-button
        class="filter-item"
        style="margin-left: 10px"
        type="primary"
        @click="handleFilter"
      >
        {{ $t("table.search") }}
      </el-button>
      <el-button
        class="filter-item"
        style="margin-left: 10px"
        type="primary"
        @click="handleCreate"
      >
        {{ $t("table.add") }}
      </el-button>
    </div>

    <el-table
      v-loading="listLoading"
      :data="list.data"
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
      <el-table-column
        label="描述"
        prop="describe"
        width="410px"
        sortable="custom"
        :class-name="getSortClass('describe')"
      >
        <template #default="{ row }">
          <router-link
            :to="{ name: 'Flow', query: { pipelineId: row.id, from: 'Pipeline' } }"
            >{{ row.describe }}</router-link
          >
        </template>
      </el-table-column>
      <el-table-column
        label="类型"
        prop="type"
        align="center"
        sortable="custom"
        :class-name="getSortClass('type')"
      >
        <template #default="{ row }">
          <span>{{ row.type }}</span>
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
        <template #default="{ row, $index }">
          <el-button-group>
            <el-button type="primary" size="small" @click="handleUpdate(row)">
              Edit
            </el-button>
            <el-button type="success" size="small" @click="handleTask(row)">
              Task
            </el-button>
            <el-button type="success" size="small" @click="handleRun(row)">
              Run
            </el-button>
            <el-popconfirm
              title="这是一段内容确定删除吗？"
              @confirm="handleDelete(row, $index)"
            >
              <template #reference>
                <el-button v-if="row.status != 'deleted'" size="small" type="danger"
                  >Delete</el-button
                >
              </template>
            </el-popconfirm>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>

    <pagination
      v-show="total > 0"
      v-model:currentPage="page"
      v-model:pageSize="listQuery.limit"
      :total="total"
      @pagination="paginationHandleGetList"
    />

    <el-dialog v-model="dialogFormVisible" :title="textMap[dialogStatus]">
      <el-form
        ref="dataForm"
        :rules="rules"
        :model="temp"
        label-position="left"
        label-width="70px"
        style="width: 400px; margin-left: 50px"
      >
        <el-form-item label="描述">
          <el-input
            v-model="temp.describe"
            :autosize="{ minRows: 2, maxRows: 4 }"
            type="textarea"
            placeholder="Please input"
          />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="temp.type">
            <el-option
              v-for="item in types"
              :key="item.value"
              :label="item.name"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <div v-if="temp.type === 'performance'">
          <el-form-item label="qps">
            <el-input-number v-model="temp.config.qps" :min="0" />
          </el-form-item>
          <el-form-item label="time">
            <el-input-number v-model="temp.config.time" :min="1" />
          </el-form-item>
          <el-form-item label="users">
            <el-input-number v-model="temp.config.users" :min="1" />
          </el-form-item>
        </div>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogFormVisible = false">
            {{ $t("table.cancel") }}
          </el-button>
          <el-button
            type="primary"
            @click="dialogStatus === 'create' ? createData() : updateData()"
          >
            {{ $t("table.confirm") }}
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>
