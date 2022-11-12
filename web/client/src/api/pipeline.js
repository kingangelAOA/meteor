import request from '@/utils/axiosReq'

export function getPiplines(params) {
    return request({
        url: '/pipeline/list',
        method: 'get',
        params
    })
}

export function createPipeline(data) {
    return request({
        url: '/pipeline',
        method: 'post',
        data
    })
}

export function updatePipeline(data) {
    return request({
        url: '/pipeline',
        method: 'put',
        data
    })
}

export function deletePipeline(params) {
    return request({
        url: '/pipeline',
        method: 'delete',
        params
    })
}

export function updateFlow(data) {
    return request({
        url: '/pipeline/flow',
        method: 'put',
        data
    })
}

export function getFlow(params) {
    return request({
        url: '/pipeline/flow',
        method: 'get',
        params
    })
}

export function run(params) {
    return request({
        url: '/pipeline/run',
        method: 'post',
        params
    })
}

export function getTasks(params) {
    return request({
        url: '/pipeline/tasks',
        method: 'get',
        params
    })
}

export function getTypes() {
    return request({
        url: '/pipeline/types',
        method: 'get',
    })
}