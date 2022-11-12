<template>
  <el-container>
    <!-- <el-header style="background: #3399ff;border-radius: 4px">
      <span style="font-size: 40px;color:white;">调试脚本</span>
    </el-header> -->

    <multipane class="vertical-panes" layout="vertical">
      <div class="pane" :style="{ width: '65%' }">
        <el-card class="box-card">
          <div v-if="!destroy" ref="report" />
        </el-card>
      </div>
      <multipane-resizer />
      <div class="pane" :style="{ width: '35%' }">
        <el-form
          ref="form"
          :model="form"
          label-width="80px"
          style="padding: 5px"
        >
          <el-form-item label="data">
            <el-input v-model="form.data" type="textarea" :rows="4" />
          </el-form-item>
          <el-form-item label="vue 模板">
            <code-mirror
              :read-only="false"
              :code="form.model"
              height="500px"
              mode="javascript"
              theme="material"
              @change="updateScriptContent"
            />
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              :loading="loading"
              @click="debug"
            >调试</el-button>
          </el-form-item>
        </el-form>
      </div>
    </multipane>
  </el-container>
</template>

<script>
import Vue from 'vue'
import { Message } from 'element-ui'
import { Multipane, MultipaneResizer } from 'vue-multipane'
import CodeMirror from '@/components/codemirror/code-mirror'
import { getNode, updateNodes } from '@/api/pipeline-api'
export default {
  name: 'ScriptDebug',
  components: {
    Multipane,
    MultipaneResizer,
    CodeMirror
  },

  data() {
    return {
      code: '',
      form: {
        data: '',
        model: ''
      },
      id: '',
      nodeId: '',
      result: '',
      loading: false,
      model: '',
      destroy: false
    }
  },
  mounted() {
    const id = this.$route.query.id
    this.pipelineId = this.$route.query.pipelineId
    this.id = id
    getNode({ id }).then((res) => {
      const config = JSON.parse(res.data.config)
      const { model, data } = config
      this.form.model =
        model === undefined || model === ''
          ? '<div>没有设置template</div>'
          : model
      this.form.data = data === undefined ? '{}' : data
      this.compile()
    })
  },
  methods: {
    updateScriptContent(content) {
      this.form.model = content
    },
    debug() {
      getNode({ id: this.id }).then((res) => {
        const { data } = res
        const config = JSON.parse(data.config)
        config.model = this.form.model
        config.data = this.form.data
        data.config = Buffer.from(JSON.stringify(config)).toString('base64')
        updateNodes([data])
      })
      this.compile()
    },
    compile() {
      const that = this
      Vue.config.warnHandler = function (err) {
        Message({
          message:
            `Custom vue error handler: ${err}, template: ${
              that.template
            }, data: ${JSON.stringify(that.options)}` || 'Error',
          type: 'error',
          duration: 5 * 1000
        })
      }
      const Component = Vue.extend({
        data() {
          return JSON.parse(that.form.data)
        },
        template: this.form.model,
        methods: {
          tableRowStyle: new Function(
            'object',
            'const { row, rowIndex } = object; if (rowIndex === 0) { return { background: "#f0f9eb" }; } return "";'
          )
        }
      })
      const report = new Component().$mount()
      this.$refs['report'].innerHTML = ''
      this.$refs['report'].appendChild(report.$el)
    }
  }
}
</script>

<style>
.vertical-panes {
  width: 100%;
  height: 100%;
  border: 1px solid #ffffff;
}
.vertical-panes > .pane {
  text-align: left;
  padding: 5px;
  overflow: hidden;
  background: #ffffff;
}
.vertical-panes > .pane ~ .pane {
  border-left: 1px solid #ccc;
}
.text {
  font-size: 14px;
}

.item {
  margin-bottom: 18px;
}

.clearfix:before,
.clearfix:after {
  display: table;
  content: "";
}
.clearfix:after {
  clear: both;
}

.box-card {
  width: 100%;
}
</style>
