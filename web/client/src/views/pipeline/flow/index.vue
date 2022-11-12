<template>
  <div class="app-container">
    <flow :pipelineId="pipelineId" :readOnly="readOnly" :taskId="taskId" />
  </div>
</template>

<script setup name="flow">
import Flow from "@/components/Flow/index.vue";
const route = useRoute();
const pipelineId = ref(route.query.pipelineId);
const taskId = ref(route.query.taskId);
const readOnly = ref(false);
if (route.query.from === "Task" || route.query.from === "PipelineRun") {
  readOnly.value = true;
}
</script>
