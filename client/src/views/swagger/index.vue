<template>
  <div>
    <!--    <div class="navbar">-->
    <!--      <div class="left-menu">-->
    <!--        <el-select v-model="version" placeholder="version" style="min-width: 30%">-->
    <!--          <el-option v-for="item in versions" :key="item" :label="item" :value="item" />-->
    <!--        </el-select>-->
    <!--        <el-dropdown split-button type="primary" style="min-width: 10%" @command="handleDropdown" @click="saveOrCreate">-->
    <!--          {{ $t('swagger.'+buttonText) }}-->
    <!--          <el-dropdown-menu slot="dropdown">-->
    <!--            <el-dropdown-item command="1">{{ $t('swagger.revertToLastSaved') }}</el-dropdown-item>-->
    <!--          </el-dropdown-menu>-->
    <!--        </el-dropdown>-->
    <!--      </div>-->
    <!--    </div>-->
    <multipane class="vertical-panes" layout="vertical">
      <el-menu default-active="1" style="height: 100%; overflow: scroll; width: 25%" @select="handleSideSelect">
        <el-form :inline="true" :model="formInline">
          <el-form-item>
            <el-input v-model="searchPath" placeholder="请输入内容" prefix-icon="el-icon-search" />
          </el-form-item>
          <el-form-item>
            <i class="el-icon-circle-plus" />
            <i class="el-icon-caret-bottom" />
          </el-form-item>
        </el-form>

        <el-submenu v-for="(tagObject, name, index) in interfaces" :key="index" :index="name">
          <template slot="title">
            <span class="item-title">{{ name }}</span>
          </template>
          <el-menu-item-group>
            <el-menu-item v-for="(item, index) in tagObject" :key="index" style="padding-left: 20px" :index="util.format('%s-%s', item.method, item.path)">
              <el-tag :type="getMethodTag(item.method)">{{ item.method }}</el-tag>
              &nbsp;&nbsp;&nbsp;&nbsp;{{ item.path }}
            </el-menu-item>
          </el-menu-item-group>
        </el-submenu>
      </el-menu>
      <multipane-resizer />
      <div class="editor-container">
        <code-mirror v-model="value" mode="yaml" theme="material" :search-base-text="baseText" :interface-text="interfaceText" />
      </div>
      <multipane-resizer />
      <div class="pane" :style="{ flexGrow: 1 }">
        <div>
          <h6 class="title is-6">Pane 3</h6>
          <p class="subtitle is-6">Takes remaining available space.</p>
          <p>
            <small>
              <strong>Configured with:</strong><br>
              flex-grow: 1<br>
            </small>
          </p>
        </div>
      </div>
    </multipane>
  </div>
</template>

<script>
import {
  getSwaggerVersions,
  getSwagger,
  createOrUpdateSwagger
} from '@/api/swagger'
import {
  Message
} from 'element-ui'
import {
  Multipane,
  MultipaneResizer
} from 'vue-multipane'
import CodeMirror from './code-mirror'
import util from 'util'
import YAML from 'yaml'
import Enforcer from 'openapi-enforcer'

const baseField = ['info', 'tags', 'servers']
export default {
  name: 'Swagger',
  components: {
    Multipane,
    MultipaneResizer,
    CodeMirror
  },
  data() {
    return {
      projectID: '',
      versions: [],
      value: 'openapi: 3.0.1',
      version: undefined,
      searchPath: '',
      buttonText: 'save',
      swaggerObject: undefined,
      initSwaggerObject: undefined,
      saveButtonDisable: 'true',
      init: false,
      baseText: '',
      interfaceText: '',
      util: util,
      sideValue: '1',
      interfaces: {}
    }
  },
  watch: {
    version: function () {
      this.getSwagger()
      this.init = false
    },
    versions: function (val) {
      if (this.version === undefined) {
        this.version = val[0]
      }
    },
    value: function (val) {
      this.swaggerObject = this.parseSwaggerYaml(val)
      if (!this.init) {
        this.initSwaggerObject = this.parseSwaggerYaml(this.value)
        this.init = true
      }
      this.triggerNewVersion()
      this.initInterfaces('')
      this.validateSwagger()
    },
    searchPath: function (val) {
      this.initInterfaces(val)
    }
  },
  mounted() {
    this.projectID = this.$route.query.id
    this.getSwaggerVersions()
  },
  methods: {
    getSwaggerVersions() {
      getSwaggerVersions({
          projectID: this.projectID
        })
        .then(res => {
          this.versions = res.data
        })
    },
    getSwagger() {
      getSwagger({
          projectID: this.projectID,
          version: this.version
        })
        .then(res => {
          this.value = res.data === '' ? 'openapi: 3.0.1' : YAML.stringify(YAML.parse(res.data))
        })
    },
    parseSwaggerYaml(val) {
      try {
        return YAML.parse(val)
      } catch (e) {
        Message({
          message: 'yaml format error',
          type: 'error',
          duration: 1000
        })
      }
    },
    triggerNewVersion() {
      if (this.swaggerObject !== undefined && 'info' in this.swaggerObject) {
        const currentVersion = this.swaggerObject.info.version
        if (this.initSwaggerObject !== undefined && 'info' in this.initSwaggerObject) {
          const intiVersion = this.initSwaggerObject.info.version
          if (intiVersion !== currentVersion) {
            this.buttonText = 'createNewVersion'
          } else {
            this.buttonText = 'save'
          }
        } else {
          this.buttonText = 'save'
        }
      }
    },
    saveOrCreate() {
      createOrUpdateSwagger({
          id: this.projectID,
          swagger: JSON.stringify(this.swaggerObject)
        })
        .then(() => {
          this.getSwagger()
          this.getSwaggerVersions()
          Message({
            message: 'save or create swagger success',
            type: 'success',
            duration: 1000
          })
        })
    },
    handleDropdown(e) {
      if (e === '1') {
        this.value = YAML.stringify(this.initSwaggerObject)
      }
    },
    initInterfaces(text) {
      this.interfaces = {}
      for (const [path, o1] of Object.entries(this.swaggerObject.paths)) {
        if (text === '') {
          this.pushInterface(path, o1)
        } else {
          if (path.includes(text)) {
            this.pushInterface(path, o1)
          }
        }
      }
      if ('components' in this.swaggerObject) {
        for (const key of Object.keys(this.swaggerObject.components.schemas)) {
          if (text === '') {
            this.pushSchemas(key)
          } else {
            if (key.includes(text)) {
              this.pushSchemas(key)
            }
          }
        }
      }
    },
    pushSchemas(key) {
      if ('Schemas' in this.interfaces) {
        this.interfaces['Schemas'].push({
          method: 'SCHEMA',
          path: key
        })
      } else {
        this.interfaces['Schemas'] = [{
          method: 'SCHEMA',
          path: key
        }]
      }
    },
    pushInterface(path, o1) {
      for (const [method, o2] of Object.entries(o1)) {
        if ('tags' in o2) {
          for (const tag of o2.tags) {
            if (tag in this.interfaces) {
              this.interfaces[tag].push({
                path,
                method
              })
            } else {
              this.interfaces[tag] = [{
                path,
                method
              }]
            }
          }
        }
      }
    },
    handleSideSelect(index) {
      if (baseField.includes(index)) {
        this.baseText = util.format('%s:', index)
      } else {
        this.interfaceText = index
      }
    },
    getMethodTag(method) {
      if (method === 'post') {
        return 'success'
      } else if (method === 'get') {
        return ''
      } else if (method === 'put') {
        return 'warning'
      } else if (method === 'delete') {
        return 'danger'
      }
      return 'info'
    },
    validateSwagger() {
      Enforcer(this.swaggerObject, {
          fullResult: true
        })
        .then(function ({
          error,
          warning
        }) {
          if (!error) {
            console.log('No errors with your document')
            if (warning) {
              console.warn(warning)
            }
          } else {
            console.error(error)
          }
        })
    }
  }
}
</script>

<style lang="scss" scoped>
.vertical-panes {
  width: 100%;
  height: 90%;
  position: absolute;
  border: 1px solid #ccc;
}

.el-collapse {
  height: 100%;
  overflow: scroll;
  width: 15%;
  padding-left: 10px
}

.vertical-panes>.pane {
  text-align: left;
  padding: 15px;
  overflow: hidden;
  background: #eee;
}

.editor-container {
  height: 100%;
  overflow: scroll;
  width: 50%;
  position: relative;
}

.item-title {
  font-weight: bold;
  padding-left: 20px;
}

.navbar {
  min-height: 10%;
  overflow: hidden;
  position: relative;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, .08);

  .hamburger-container {
    line-height: 46px;
    height: 100%;
    float: left;
    cursor: pointer;
    transition: background .3s;
    -webkit-tap-highlight-color: transparent;

    & :hover {
      background: rgba(0, 0, 0, .025)
    }

  }

  .breadcrumb-container {
    float: left;
  }

  .errLog-container {
    display: inline-block;
    vertical-align: top;
  }

  .left-menu {
    float: left;
    height: 100%;
    line-height: 50px;
    padding: 0px 0px 0px 4px;
  }

  .right-menu {
    float: right;
    height: 100%;
    line-height: 50px;
    padding: 0px 4px 0px 0px;
  }
}
</style>
