<script setup name="PluginEdit">
import {
  getPluginByID,
  createPlugin,
  updatePlugin,
  uploadPluginFile,
  debugPlugin,
  checkPluginFileStatus,
} from "@/api/plugin";
import MonacoEditor from "@/components/monaco/index.vue";
import { Multipane, MultipaneResizer } from "vue-multipane";
import { reactive, ref } from "vue";
import { useRoute } from "vue-router";
import { ElMessage } from "element-plus";

const route = useRoute();
const formInline = reactive({
  id: "",
  name: "",
  desc: "",
  inputs: [],
  code: "",
  language: "go",
});
const log = ref("");
const loading = ref(false);
let fileList = ref([]);
formInline.id = route.query.id;
const addInput = () => {
  formInline.inputs.push({ name: "", desc: "", value: "", required: false });
};
const getPlugin = () => {
  loading.value = true;

  getPluginByID({ id: formInline.id })
    .then((res) => {
      let { name, desc, inputs, code, language, id } = res.data;
      if (name === undefined) {
        name = "";
      }
      if (id === undefined) {
        id = "";
      }
      if (desc === undefined) {
        desc = "";
      }
      if (inputs === undefined) {
        inputs = [];
      }
      if (code === undefined) {
        code = "";
      }
      if (language === undefined) {
        language = "go";
      }
      for (const input of inputs) {
        input.originalName = input.name;
        input.originalDesc = input.desc;
        input.originalValue = input.value;
      }
      formInline.id = id;
      formInline.name = name;
      formInline.desc = desc;
      formInline.inputs = inputs;
      formInline.code = code;
      formInline.language = language;
    })
    .finally(() => {
      loading.value = false;
      checkPlugin();
    });
};
const checkPlugin = () => {
  checkPluginFileStatus({ id: formInline.id, language: formInline.language });
};
if (formInline.id !== "") {
  getPlugin();
}
const cancelEdit = (row, index) => {
  row.name = row.originalName;
  row.value = row.originalValue;
  row.desc = row.originalDesc;
  row.edit = false;
  formInline.inputs[index] = row;
  ElMessage({
    message: "The title has been restored to the original value",
    type: "warning",
  });
};
const edit = (row, index) => {
  row.edit = !row.edit;
  formInline.inputs[index] = row;
};
const confirmEdit = (row, index) => {
  row.edit = false;
  row.originalName = row.name;
  row.originalDesc = row.desc;
  row.originalValue = row.value;
  formInline.inputs[index] = row;
  ElMessage({
    message: "The title has been edited",
    type: "success",
  });
};
const handleCode = (value) => {
  formInline.code = value;
};
const save = () => {
  loading.value = true;
  if (formInline.id === "") {
    createPlugin(formInline)
      .then((res) => {
        formInline.id = res.data;
      })
      .finally(() => {
        getPlugin();
      });
  } else {
    updatePlugin(formInline).finally(() => {
      getPlugin();
    });
  }
};
const debug = () => {
  loading.value = true;
  debugPlugin(formInline)
    .then((res) => {
      log.value = res.data;
    })
    .finally(() => {
      getPlugin();
      loading.value = false;
    });
};
const handleUpload = (params) => {
  // progressPercent.value = 0
  const file = params.file;
  // this.ruleForm.packageSize = (file.size / (1024 * 1024)).toFixed(2) + 'M' // 文件大小，转化成M
  const forms = new FormData(); // 实例化一个formData，用来做文件上传
  forms.append("file", file);
  if (this.formInline.id === "") {
    ElMessage.error("上传插件文件前, 先保存插件配置");
  } else {
    forms.append("id", this.formInline.id);
    uploadPluginFile(forms).then((res) => {
      console.log(res);
    });
  }
};
const handleChange = (file, fl) => {
  fileList.value = fl.slice(-3);
};
</script>
<template>
  <div class="app-container">
    <el-form :inline="true" :model="formInline" class="demo-form-inline">
      <el-form-item label="名称">
        <el-input
          v-model="formInline.name"
          size="small"
          placeholder="名称"
          style="width: 300px"
        />
      </el-form-item>
      <el-form-item label="描述">
        <el-input
          v-model="formInline.desc"
          type="text"
          size="small"
          placeholder="描述"
          style="width: 300px"
        />
      </el-form-item>
    </el-form>
    <el-divider>入参</el-divider>
    <div style="padding-bottom: 5px">
      <el-button
        class="filter-item"
        type="primary"
        size="small"
        icon="el-icon-edit"
        @click="addInput"
      >
        {{ $t("table.add") }}
      </el-button>
    </div>
    <el-table
      v-loading="loading"
      :data="formInline.inputs"
      size="small"
      max-height="150"
      border
      fit
      highlight-current-row
      style="width: 100%"
    >
      <el-table-column label="参数名">
        <template #default="{ row }">
          <template v-if="row.edit">
            <el-input v-model="row.name" class="edit-input" size="small" />
          </template>
          <span v-else>{{ row.name }}</span>
        </template>
      </el-table-column>
      <el-table-column label="描述">
        <template #default="{ row }">
          <el-input v-if="row.edit" v-model="row.desc" class="edit-input" size="small" />
          <span v-else>{{ row.desc }}</span>
        </template>
      </el-table-column>
      <el-table-column label="必须">
        <template #default="{ row }">
          <el-select
            v-if="row.edit"
            v-model="row.required"
            class="m-2"
            placeholder="Select"
            size="small"
          >
            <el-option label="是" :value="true" />
            <el-option label="否" :value="false" />
          </el-select>
          <span v-else>{{ row.required }}</span>
        </template>
      </el-table-column>
      <el-table-column label="参数值">
        <template #default="{ row }">
          <el-input v-if="row.edit" v-model="row.value" class="edit-input" size="small" />
          <span v-else>{{ row.value }}</span>
        </template>
      </el-table-column>

      <el-table-column align="center" label="Actions" width="240">
        <template #default="{ row, $index }">
          <el-button-group>
            <el-button
              v-if="row.edit"
              size="small"
              icon="el-icon-refresh"
              type="warning"
              @click="cancelEdit(row, $index)"
            >
              cancel
            </el-button>
            <el-button
              v-if="row.edit"
              type="success"
              size="small"
              icon="el-icon-circle-check-outline"
              @click="confirmEdit(row, $index)"
            >
              Ok
            </el-button>
            <el-button
              v-else
              type="primary"
              size="small"
              icon="el-icon-edit"
              @click="edit(row, $index)"
            >
              Edit
            </el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>
    <el-divider>脚本</el-divider>
    <el-row :gutter="28">
      <el-col :span="12">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <el-select
                v-model="formInline.language"
                size="small"
                placeholder=""
                style="padding-bottom: 10px; padding-right: 5px"
              >
                <el-option label="golang" value="go" />
                <el-option label="python" value="python" />
              </el-select>
            </div>
          </template>
          <MonacoEditor
            :value="formInline.code"
            :language="formInline.language"
            :key="1"
            height="500"
            @change="handleCode"
          />
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <el-button-group style="padding-bottom: 10px">
                <el-button type="primary" :loading="loading" size="small" @click="save"
                  >保存</el-button
                >
                <el-button
                  v-if="formInline.id !== ''"
                  :loading="loading"
                  type="primary"
                  size="small"
                  @click="debug"
                  >调试</el-button
                >
              </el-button-group>
            </div>
          </template>
          <MonacoEditor
            :key="2"
            :value="log"
            language="text"
            :height="500"
            @change="handleCode"
            :options="{ readOnly: true }"
          />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>
<style lang="scss">
.edit-input {
  padding-right: 100px;
}
.cancel-btn {
  position: absolute;
  right: 15px;
  top: 10px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.text {
  font-size: 14px;
}

.item {
  margin-bottom: 18px;
}

.box-card {
  height: 100%;
  overflow: hidden;
}
</style>
