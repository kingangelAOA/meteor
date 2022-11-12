<script setup>
import {
  ConnectionMode,
  VueFlow,
  useVueFlow,
  MiniMap,
  Controls,
  BackgroundVariant,
  Background,
} from "@braks/vue-flow";
import { markRaw, onUnmounted } from "vue";
import PluginNode from "./Node/plugin.vue";
import DefaultEdge from "./edge/default.vue";
import Sidebar from "./Sidebar.vue";
import { v4 as uuidv4 } from "uuid";
import { updateFlow, getFlow, run, disabled } from "@/api/pipeline.js";
import { disabled } from "@/api/task.js";
const props = defineProps({
  pipelineId: {
    type: String,
    required: true,
  },
  readOnly: {
    type: Boolean,
    required: false,
    default: false,
  },
  taskId: {
    type: String,
    required: false,
  },
});
const span = ref(20);
if (props.readOnly) {
  span.value = 24;
}
const nodeTypes = {
  plugin: markRaw(PluginNode),
};
const edgeTypes = {
  default: markRaw(DefaultEdge),
};
const nodes = ref([]);
const edges = ref([]);
const wrapper = ref();
const runningInfo = ref([]);
const { getNodes, getEdges, onConnect, addEdges, addNodes, project } = useVueFlow({
  fitViewOnInit: true,
  connectionMode: ConnectionMode.Loose,
  elevateEdgesOnSelect: true,
  nodes: [],
});
const onDragOver = (event) => {
  event.preventDefault();
  if (event.dataTransfer) {
    event.dataTransfer.dropEffect = "move";
  }
};

const deleteNode = (id) => {
  let ns = [];
  let es = [];
  for (const node of getNodes.value) {
    if (node.id !== id) {
      ns.push(node);
    }
  }
  for (const edge of getEdges.value) {
    if (edge.source === id || edge.target === id) {
      continue;
    }
    es.push(edge);
  }
  nodes.value = ns;
  edges.value = es;
};

onConnect((params) => {
  params["animated"] = true;
  params["type"] = "default";
  params["id"] = uuidv4();
  addEdges([params]);
});

watch(runningInfo, (infos) => {
  for (const info of infos) {
    const { id } = info;
    for (let n of nodes.value) {
      if (n.id === id) {
        n.data.runningInfo = info;
      }
    }
  }
});
const onDrop = (event) => {
  const flowbounds = wrapper.value.$el.getBoundingClientRect();
  const data = JSON.parse(event.dataTransfer?.getData("application/vueflow"));
  let { type, id } = data;
  const position = project({
    x: event.clientX - flowbounds.left,
    y: event.clientY - flowbounds.top,
  });

  const uuid = uuidv4();
  const newNode = {
    id: uuid,
    type,
    data,
    position,
    pipelineId: props.pipelineId,
    events: {
      deleteNode: deleteNode,
    },
  };
  addNodes([newNode]);
};
const loading = ref(false);
const saveFlow = () => {
  let flow = { pipelineId: props.pipelineId, nodes: [], edges: [] };
  for (const node of getNodes.value) {
    const { defaultInputs, describe, name, pluginID, code, language } = node.data;
    flow.nodes.push({
      id: node.id,
      nodeId: node.data.id,
      position: node.position,
      type: node.type,
      data: { defaultInputs, describe, name, pluginID, code, language },
    });
  }
  for (const edge of getEdges.value) {
    flow.edges.push({
      id: edge.id,
      type: edge.type,
      source: edge.source,
      target: edge.target,
      animated: edge.animated,
    });
  }
  loading.value = true;
  updateFlow(flow).finally(() => {
    initFlow();
    loading.value = false;
  });
};

const initFlow = () => {
  getFlow({ pipelineId: props.pipelineId }).then((res) => {
    if (res.data == null) {
      nodes.value = [];
      edges.value = [];
    } else {
      const { nodes: ns, edges: es } = res.data;
      if (ns !== undefined) {
        for (let n of ns) {
          n.data.readOnly = props.readOnly;
          n.events = {
            deleteNode: deleteNode,
          };
        }
      }
      if (es !== undefined) {
        for (let e of es) {
          e.data = {};
          e.data.readOnly = props.readOnly;
        }
      }
      nodes.value = ns === undefined ? [] : ns;
      edges.value = es === undefined ? [] : es;
    }
  });
};
initFlow();
let taskId = ref(props.taskId);
const runPipeline = () => {
  loading.value = true;
  run({ id: props.pipelineId })
    .then((res) => {
      taskId.value = res.data;
    })
    .finally(() => (loading.value = false));
};
const initRunningInfo = () => {
  if (taskId.value !== undefined && taskId.value !== "") {
    disabled({ id: taskId.value }).then((res) => {
      runningInfo.value = res.data;
    });
  }
};
initRunningInfo();
const timer = ref(null);
onMounted(() => {
  timer.value = setInterval(async () => {
    initRunningInfo();
  }, 5000);
});
onUnmounted(() => clearInterval(timer.value));
</script>

<template>
  <el-row :gutter="24" class="dndflow" @drop="onDrop">
    <el-col :span="span">
      <el-button-group v-if="!props.readOnly" class="ml-4" style="float: right">
        <el-button :loading="loading" type="primary" @click="saveFlow">保存</el-button>
        <el-button :loading="loading" type="primary" @click="runPipeline">运行</el-button>
      </el-button-group>
      <VueFlow
        ref="wrapper"
        :nodes="nodes"
        :edges="edges"
        :node-types="nodeTypes"
        :edge-types="edgeTypes"
        @dragover="onDragOver"
      >
        <MiniMap />
        <Controls />
        <Background pattern-color="#aaa" gap="8" />
      </VueFlow>
    </el-col>
    <el-col v-if="!props.readOnly" :span="4">
      <Sidebar />
    </el-col>
  </el-row>
</template>

<style lang="scss">
.el-row {
  margin-bottom: 20px;
}
.el-row:last-child {
  margin-bottom: 0;
}
.el-col {
  border-radius: 4px;
}

.grid-content {
  border-radius: 4px;
  min-height: 36px;
}
.vue-flow__minimap {
  transform: scale(75%);
  transform-origin: bottom right;
}

.dndflow {
  flex-direction: column;
  display: flex;
  height: 100%;
}
.dndflow aside {
  color: #fff;
  font-weight: 700;
  border-right: 1px solid #eee;
  padding: 15px 10px;
  font-size: 12px;
  background: rgba(16, 185, 129, 0.75);
  -webkit-box-shadow: 0px 5px 10px 0px rgba(0, 0, 0, 0.3);
  box-shadow: 0 5px 10px #0000004d;
}
.dndflow aside .nodes > * {
  margin-bottom: 10px;
  cursor: grab;
  font-weight: 500;
  -webkit-box-shadow: 5px 5px 10px 2px rgba(0, 0, 0, 0.25);
  box-shadow: 5px 5px 10px 2px #00000040;
}
.dndflow aside .description {
  margin-bottom: 10px;
}
.dndflow .vue-flow-wrapper {
  flex-grow: 1;
  height: 100%;
}
@media screen and (min-width: 640px) {
  .dndflow {
    flex-direction: row;
  }
  .dndflow aside {
    min-width: 25%;
  }
}
@media screen and (max-width: 639px) {
  .dndflow aside .nodes {
    display: flex;
    flex-direction: row;
    gap: 5px;
  }
}
</style>
