<script setup name="node">
import { getNodes, updateNode, createNode, deleteNode } from "@/api/node";
import { getSelectList } from "@/api/plugin";
// import waves from '@/directive/waves' // waves directive
import { parseTime } from "@/utils";
import Pagination from "@/components/Pagination/index.vue"; // secondary package based on el-pagination
import { reactive, ref, nextTick } from "vue";
import { ElMessage, ElNotification } from "element-plus";

const dataForm = ref(null);
const tableKey = ref(0);
let list = reactive({
  data: [],
});
let total = ref(0);
let listLoading = ref(false);
let page = ref(1);
const listQuery = reactive({
  offset: 0,
  limit: 20,
  sortField: "id",
  sortType: -1,
  describe: "",
});
let temp = ref({
  id: undefined,
  describe: "",
  name: "",
  pluginID: "",
});
let dialogUpdateVisible = ref(false);
let dialogStatus = ref("");
const textMap = reactive({
  update: "Edit",
  create: "Create",
});
const rules = reactive({
  type: [{ required: true, message: "type is required", trigger: "change" }],
  timestamp: [
    { type: "date", required: true, message: "timestamp is required", trigger: "change" },
  ],
  title: [{ required: true, message: "title is required", trigger: "blur" }],
});
let plugins = ref([]);
let pluginMap = reactive({});
watch(
  () => page,
  (val) => {
    if (val < 1) {
      val = 0;
    }
    listQuery.offset = (val - 1) * listQuery.limit;
  }
);
const paginationHandleGetList = (data) => {
  const { page, limit } = data;
  listQuery.limit = limit;
  listQuery.offset = (page - 1) * limit;
  getList();
};
const getList = () => {
  listLoading.value = true;
  getNodes(listQuery)
    .then((response) => {
      list.data = response.data.items;
      total.value = response.data.total;
    })
    .finally(() => {
      listLoading.value = false;
    });
};
getList();
const getPluginList = () => {
  getSelectList().then((res) => {
    plugins.value = res.data;
    for (const plugin of plugins.value) {
      const { id, name } = plugin;
      pluginMap[id] = name;
    }
  });
};
getPluginList();
const handleFilter = () => {
  page = 1;
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
  if (order === "ascending") {
    listQuery.sortType = 1;
  } else {
    listQuery.sortType = -1;
  }
  listQuery.sortField = prop;
  handleFilter();
};
const resetTemp = () => {
  temp.value = Object.assign(
    {},
    {
      id: undefined,
      name: "",
      describe: "",
      pluginID: "",
      createTime: "",
    }
  );
};
const handleCreate = () => {
  resetTemp();
  dialogStatus.value = "create";
  dialogUpdateVisible.value = true;
  nextTick(() => {
    dataForm.value?.clearValidate();
  });
};
const handleUpdate = (row) => {
  temp.value = Object.assign({}, row);
  dialogStatus.value = "update";
  dialogUpdateVisible.value = true;
  nextTick(() => {
    dataForm.value?.clearValidate();
  });
};

const updateData = () => {
  dialogUpdateVisible.value = true;
  dataForm.value.validate((valid) => {
    if (valid) {
      updateNode(temp.value)
        .then(() => {
          getList();
        })
        .finally(() => {
          dialogUpdateVisible.value = false;
        });
    }
  });
};
const createData = () => {
  dialogUpdateVisible.value = true;
  dataForm.value.validate((valid) => {
    if (valid) {
      createNode(temp.value)
        .then(() => {
          getList();
        })
        .finally(() => {
          dialogUpdateVisible.value = false;
        });
    }
  });
};
const handleDelete = (row, index) => {
  deleteNode({ id: row.id })
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
  return list.data.map((v) =>
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
</script>
<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input
        v-model="listQuery.describe"
        placeholder="描述"
        style="width: 400px"
        class="filter-item"
        @keyup.enter="handleFilter"
      />
      <el-button
        v-waves
        class="filter-item"
        type="primary"
        icon="el-icon-search"
        @click="handleFilter"
      >
        {{ $t("table.search") }}
      </el-button>
      <el-button
        class="filter-item"
        style="margin-left: 10px"
        type="primary"
        icon="el-icon-edit"
        @click="handleCreate"
      >
        {{ $t("table.add") }}
      </el-button>
    </div>

    <el-table
      :key="tableKey"
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
        label="名称"
        prop="name"
        sortable="custom"
        :class-name="getSortClass('describe')"
      >
        <template #default="{ row }">
          <router-link :to="{ name: 'Flow', params: { nodeId: row.id } }">{{
            row.name
          }}</router-link>
        </template>
      </el-table-column>
      <el-table-column
        label="描述"
        prop="describe"
        sortable="custom"
        :class-name="getSortClass('describe')"
      >
        <template #default="{ row }">
          <router-link :to="{ name: 'Flow', params: { nodeId: row.id } }">{{
            row.describe
          }}</router-link>
        </template>
      </el-table-column>
      <el-table-column
        label="插件"
        prop="pluginID"
        align="center"
        sortable="custom"
        :class-name="getSortClass('type')"
      >
        <template #default="{ row }">
          <span>{{ pluginMap[row.pluginID] }}</span>
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
        class-name="small-padding fixed-width"
      >
        <template #default="{ row, $index }">
          <el-button type="primary" size="small" @click="handleUpdate(row)">
            {{ $t("table.edit") }}
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

    <el-dialog v-model="dialogUpdateVisible" :title="textMap[dialogStatus]">
      <el-form
        ref="dataForm"
        :rules="rules"
        :model="temp"
        label-position="left"
        label-width="70px"
        style="width: 400px; margin-left: 50px"
      >
        <el-form-item label="名称">
          <el-input v-model="temp.name" placeholder="Please input" />
        </el-form-item>
        <el-form-item label="plugin" prop="pluginID">
          <el-select
            v-model="temp.pluginID"
            class="filter-item"
            filterable
            placeholder="Please select"
          >
            <el-option
              v-for="item in plugins"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            v-model="temp.describe"
            :autosize="{ minRows: 2, maxRows: 4 }"
            type="textarea"
            placeholder="Please input"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogUpdateVisible = false">
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
