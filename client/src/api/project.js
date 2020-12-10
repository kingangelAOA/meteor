import request from '@/utils/request'

export function createOrUpdateProject(data) {
  return request({
    headers: {
      'Content-Type': 'application/json;charset=utf-8'
    },
    url: '/project',
    method: 'post',
    data: JSON.stringify(data)
  })
}

export function getAllProjects(query) {
  return request({
    url: '/project/all',
    method: 'get',
    params: query
  })
}
