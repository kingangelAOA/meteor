<template>
  <div class="yaml-editor">
    <textarea ref="textarea" />
  </div>
</template>

<script>
import CodeMirror from 'codemirror'
import util from 'util'
import 'codemirror/addon/lint/lint.css'
import 'codemirror/addon/search/search'
import 'codemirror/addon/search/searchcursor'
import 'codemirror/addon/selection/mark-selection'
import 'codemirror/addon/selection/active-line.js'
import 'codemirror/lib/codemirror.css'
import 'codemirror/theme/rubyblue.css'
import 'codemirror/theme/base16-dark.css'
import 'codemirror/theme/material.css'
import 'codemirror/theme/idea.css'
import 'codemirror/mode/javascript/javascript'
import 'codemirror/mode/yaml/yaml'
import 'codemirror/mode/yaml-frontmatter/yaml-frontmatter'
import 'codemirror/addon/lint/lint'
import 'codemirror/addon/lint/yaml-lint'
window.jsyaml = require('js-yaml')

export default {
  name: 'CodeMirror',
  /* eslint-disable vue/require-prop-types */
  props: ['value', 'mode', 'theme', 'searchBaseText', 'interfaceText'],
  data() {
    return {
      codeEditor: false,
      arrObject: {}
    }
  },
  watch: {
    value(value) {
      const editorValue = this.codeEditor.getValue()
      if (value !== editorValue) {
        this.codeEditor.setValue(this.value)
        this.searchLineByText()
      }
    },
    searchBaseText(val) {
      this.searchBase(val)
    },
    interfaceText(val) {
      this.searchInterface(val)
    }
  },
  mounted() {
    this.codeEditor = CodeMirror.fromTextArea(this.$refs.textarea, {
      lineNumbers: true,
      mode: 'yaml',
      gutters: ['CodeMirror-lint-markers'],
      theme: 'material',
      styleSelectedText: true,
      lint: true
    })

    this.codeEditor.setValue(JSON.stringify(this.value, null, 2))
    this.codeEditor.on('change', cm => {
      this.$emit('changed', cm.getValue())
      this.$emit('input', cm.getValue())
    })
  },
  methods: {
    getValue() {
      return this.codeEditor.getValue()
    },
    search(line) {
      this.codeEditor.focus()
      this.codeEditor.setCursor(parseInt(line, 10), 0, {
        bias: 10,
        scroll: true
      })
    },
    searchLineByText() {
      var lineTotal = this.codeEditor.lineCount()
      var arrObject = {}
      for (var i = 0; i < lineTotal; i++) { // 将每行数据对应行数
        arrObject[i] = this.codeEditor.lineInfo(i).text
      }
      this.arrObject = arrObject
    },
    searchBase(baseField) {
      let line = 0
      for (const [key, value] of Object.entries(this.arrObject)) {
        if (value === baseField) {
          line = key
        }
      }
      this.search(line)
    },
    searchInterface(text) {
      const methodPath = text.split('-')
      const method = methodPath[0]
      const path = methodPath[1]
      let flag = false
      for (const [key, value] of Object.entries(this.arrObject)) {
        const v = value.trim().replace(/\"/g, '')
        if (method === 'SCHEMA') {
          if (v === 'schemas:') {
            flag = true
          }
        } else {
          console.log(v)
          if (util.format('%s:', path.trim()) === v) {
            flag = true
          }
        }
        if (flag) {
          if (method === 'SCHEMA') {
            if (util.format('%s:', path) === v) {
              this.search(key)
              return
            }
          } else {
            if (util.format('%s:', method) === v) {
              this.search(key)
              return
            }
          }
        }
      }
    }
  }
}
</script>

<style scoped>
.yaml-editor {
  height: 100%;
  position: relative;
}

.CodeMirror-selected {

  background-color: blue !important;

}

.CodeMirror-selectedtext {

  color: white !important;

}

.cm-matchhighlight {

  background-color: #ae00ae;

}

.yaml-editor>>>.CodeMirror {
  height: auto;
  min-height: 300px;
}

.yaml-editor>>>.CodeMirror-scroll {
  min-height: 300px;
}

.yaml-editor>>>.cm-s-rubyblue span.cm-string {
  color: #F08047;
}
</style>
