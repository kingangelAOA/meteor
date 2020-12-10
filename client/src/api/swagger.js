import request from '@/utils/request'

export function getSwaggerVersions(params) {
  return request({
    url: '/swagger/versions',
    method: 'get',
    params
  })
}

export function getSwagger(params) {
  return request({
    url: '/swagger',
    method: 'get',
    params
  })
}

export function createOrUpdateSwagger(data) {
  return request({
    headers: {
      'Content-Type': 'application/json;charset=utf-8'
    },
    url: '/swagger',
    method: 'post',
    data: JSON.stringify(data)
  })
}
