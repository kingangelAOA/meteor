import request from '@/utils/axiosReq'

export function getNodes(params) {
  return request({
    url: '/node/list',
    method: 'get',
    params
  })
}

export function getAllNodes(params) {
  return request({
    url: '/node/all',
    method: 'get',
    params
  })
}

export function getBindingPluginNodes(params) {
  return request({
    url: '/node/binding/plugin',
    method: 'get',
    params
  })
}

export function createNode(data) {
  return request({
    url: '/node',
    method: 'post',
    data
  })
}

export function updateNode(data) {
  return request({
    url: '/node',
    method: 'put',
    data
  })
}

export function deleteNode(params) {
  return request({
    url: '/node',
    method: 'delete',
    params
  })
}
