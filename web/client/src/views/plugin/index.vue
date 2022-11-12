<script setup name="plugin">
import { ElMessage, ElNotification } from "element-plus";
import {
  getPlugins,
  updatePlugin,
  createPlugin,
  deletePlugin,
  publishPlugin,
} from "@/api/plugin";
import { parseTime } from "@/utils";
import Pagination from "@/components/Pagination/index.vue"; // secondary package based on el-pagination
import { reactive, ref } from "vue";
import { useRouter } from "vue-router";

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
const showReviewer = ref(false);
const temp = ref({
  id: undefined,
  describe: "",
  type: "",
});
const dialogFormVisible = ref(false);
const dialogStatus = ref("");
const textMap = ref({
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
const downloadLoading = ref(false);
watch(
  () => page,
  (val) => {
    if (val < 1) {
      val = 0;
    }
    listQuery.offset = (val - 1) * listQuery.limit;
  }
);
const getList = () => {
  listLoading.value = true;
  getPlugins(listQuery)
    .then((response) => {
      list.value = response.data.items;
      total.value = response.data.total;
    })
    .finally(() => {
      listLoading.value = false;
    });
};
getList();
const handleFilter = () => {
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
  if (order === "ascending") {
    listQuery.sortType = 1;
  } else {
    listQuery.sortType = -1;
  }
  listQuery.sortField = prop;
  handleFilter();
};
const resetTemp = () => {
  temp.value = {
    id: undefined,
    importance: 1,
    remark: "",
    timestamp: new Date(),
    title: "",
    status: "published",
    type: "",
  };
};
const handleCreate = () => {
  router.push({
    path: "/plugin/edit",
    query: { id: "" },
  });
};
const dataForm = ref(null);
const updateData = () => {
  dataForm.value.validate((valid) => {
    if (valid) {
      dialogFormVisible.value = true;
      updatePlugin(temp.value)
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
      dialogFormVisible.value = true;
      createPlugin(temp.value)
        .then(() => {
          getList();
        })
        .finally(() => {
          dialogFormVisible.value = false;
        });
    }
  });
};
const handleDelete = (id) => {
  deletePlugin({ id })
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
const releaseLoading = ref(false);
const handlePublish = (id) => {
  releaseLoading.value = true;
  publishPlugin({ id }).finally(() => {
    releaseLoading.value = false;
    getList();
  });
};
const formatJson = (filterVal) => {
  return list.value.map((v) =>
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
        size="small"
        placeholder="描述"
        style="width: 400px"
        class="filter-item"
        @keyup.enter="handleFilter"
      />
      <el-button
        v-waves
        class="filter-item"
        type="primary"
        size="small"
        icon="el-icon-search"
        @click="handleFilter"
      >
        {{ $t("table.search") }}
      </el-button>
      <el-button
        class="filter-item"
        style="margin-left: 10px"
        size="small"
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
      <el-table-column
        label="名称"
        prop="name"
        sortable="custom"
        :class-name="getSortClass('describe')"
      >
        <template #default="{ row }">
          <router-link :to="{ name: 'PluginEdit', query: { id: row.id } }">{{
            row.name
          }}</router-link>
        </template>
      </el-table-column>
      <el-table-column
        label="描述"
        prop="desc"
        sortable="custom"
        :class-name="getSortClass('describe')"
      >
        <template #default="{ row }">
          <router-link :to="{ name: 'PluginEdit', query: { id: row.id } }">{{
            row.desc
          }}</router-link>
        </template>
      </el-table-column>
      <el-table-column
        label="语言"
        prop="language"
        sortable="custom"
        align="center"
        :class-name="getSortClass('language')"
      >
        <template #default="{ row }">
          <router-link :to="{ name: 'PluginEdit', query: { id: row.id } }">{{
            row.language
          }}</router-link>
        </template>
      </el-table-column>
      <el-table-column label="状态" prop="status" align="center">
        <template #default="{ row }">
          <el-tag v-if="row.status === 'offline'" class="ml-2" type="danger">离线</el-tag>
          <el-tag v-else-if="row.status === 'health'" class="ml-2" type="success"
            >正常</el-tag
          >
          <el-tag v-else-if="row.status === 'exception'" class="ml-2" type="danger"
            >错误</el-tag
          >
          <el-tag v-if="row.status === 'exception'" class="ml-2" type="warning"
            >未发布</el-tag
          >
        </template>
      </el-table-column>
      <el-table-column label="版本数量" prop="status" align="center">
        <template #default="{ row }">
          <el-tag class="ml-2" type="success">{{ row.versionNum }}</el-tag>
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
          <el-button-group class="ml-4">
            <el-popconfirm
              title="这是一段内容确定删除吗？"
              @confirm="handleDelete(row.id)"
            >
              <template #reference>
                <el-button v-if="row.status != 'deleted'" size="small" type="danger"
                  >Delete</el-button
                >
              </template>
            </el-popconfirm>
            <el-button
              :loading="releaseLoading"
              v-if="row.status != 'deleted'"
              size="small"
              type="success"
              @click="handlePublish(row.id)"
              >Publish</el-button
            >
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>
    <pagination
      v-show="total > 0"
      v-model:page="page"
      v-model:limit="listQuery.limit"
      :total="total"
      @pagination="getList"
    />
  </div>
</template>
