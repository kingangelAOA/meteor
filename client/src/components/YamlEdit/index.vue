<template>
  <div class="editor-container">
    <code-mirror ref="yamlEditor" v-model="value" mode="yaml" theme="material" />
  </div>
</template>

<script>
import CodeMirror from '@/components/CodeMirror'

const jsonData = '[{"items":[{"market_type":"forexdata","symbol":"XAUUSD"},{"market_type":"forexdata","symbol":"UKOIL"},{"market_type":"forexdata","symbol":"CORN"}],"name":""},{"items":[{"market_type":"forexdata","symbol":"XAUUSD"},{"market_type":"forexdata","symbol":"XAGUSD"},{"market_type":"forexdata","symbol":"AUTD"},{"market_type":"forexdata","symbol":"AGTD"}],"name":"贵金属"},{"items":[{"market_type":"forexdata","symbol":"CORN"},{"market_type":"forexdata","symbol":"WHEAT"},{"market_type":"forexdata","symbol":"SOYBEAN"},{"market_type":"forexdata","symbol":"SUGAR"}],"name":"农产品"},{"items":[{"market_type":"forexdata","symbol":"UKOIL"},{"market_type":"forexdata","symbol":"USOIL"},{"market_type":"forexdata","symbol":"NGAS"}],"name":"能源化工"}]'

export default {
  name: 'YamlEdit',
  components: { CodeMirror },
  data() {
    return {
      value: JSON.parse(jsonData)
    }
  }
}
</script>

<style scoped>
  .editor-container{
    height: 100%;
    overflow: scroll;
    width: 100%;
    position: relative;
  }
</style>

