import request from '@/utils/axiosReq'

export function loginReq(data) {
  console.log(data)
  return request({
    url: '/user/login',
    method: 'post',
    data
  })
}

export function getInfoReq() {
  console.log('get user info')
  return request({
    url: '/user/info',
    method: 'get'
  })
}

export function logoutReq() {
  return request({
    url: '/integration-front/user/loginOut',
    method: 'post'
  })
}
