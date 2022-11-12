<script setup name="Node">
import { Handle, Position, useVueFlow } from "@braks/vue-flow";
import { reactive } from "vue-demi";
import { QuestionFilled, Edit } from "@element-plus/icons-vue";
import { Delete } from "@element-plus/icons-vue";
import { v4 as uuidv4 } from "uuid";
const props = defineProps({
  data: {
    type: Object,
    required: true,
  },
  id: {
    type: String,
    required: true,
  },
  events: {
    type: Object,
    required: true,
  },
});
// console.log(props);
const { applyNodeChanges } = useVueFlow();
const onDelete = (evt) => {
  props.events.deleteNode(props.id);
  evt.stopPropagation();
};

// const onConnect = (params) => console.log('handle onConnect', params)
const data = ref(props.data);
// const readOnly = ref(data.readOnly);
const inputs = ref(data.value.defaultInputs);
const runningInfo = ref({});
const top = uuidv4();
const bottom = uuidv4();
watch(
  data,
  (val) => {
    const { runningInfo: newRunningInfo, readOnly: newReadOnly } = val;
    if (newRunningInfo !== undefined) {
      runningInfo.value = newRunningInfo;
      let result = [];
      for (const [key, value] of Object.entries(newRunningInfo.inputs)) {
        result.push({ name: key, value });
      }
      inputs.value = result;
      // readOnly.value = newReadOnly;
    }
  },
  { deep: true }
);
</script>

<template>
  <div>
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>{{ props.data.name }}</span>
          <el-button
            v-if="!props.data.readOnly"
            type="danger"
            :icon="Delete"
            @click="(event) => onDelete(event)"
            circle
          />
        </div>
      </template>
      <el-collapse accordion>
        <el-collapse-item title="输入" name="1">
          <el-form label-width="40px" style="max-width: 460px">
            <el-form-item v-for="(item, index) in inputs" :key="index" :label="item.name">
              <el-row :gutter="10">
                <el-col :span="22">
                  <el-input v-if="!props.data.readOnly" v-model="item.value" />
                  <span v-else>{{ item.value }}</span>
                </el-col>
              </el-row>
            </el-form-item>
          </el-form>
        </el-collapse-item>
        <el-collapse-item title="输出" name="2">
          <el-form label-width="40px" style="max-width: 460px">
            <el-form-item
              v-for="(value, key) in runningInfo.output"
              :key="key"
              :label="key"
            >
              <el-row :gutter="10">
                <el-col :span="22">
                  <span>{{ value }}</span>
                </el-col>
              </el-row>
            </el-form-item>
          </el-form>
        </el-collapse-item>
        <el-collapse-item v-if="runningInfo.status == 'failed'" name="3">
          <template #title>
            错误信息
            <el-icon class="header-icon">
              <info-filled />
            </el-icon>
          </template>
          {{ runningInfo.err }}
        </el-collapse-item>
      </el-collapse>
    </el-card>
    <Handle id="a" type="target" :position="Position.Top" />
    <Handle id="b" type="source" :position="Position.Bottom" />
  </div>
</template>

<style>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 5px;
}

.text {
  font-size: 14px;
}

.item {
  margin-bottom: 18px;
}

.box-card {
  width: 250px;
}
</style>
