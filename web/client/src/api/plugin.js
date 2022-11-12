import request from '@/utils/axiosReq'

export function getPlugins(params) {
    return request({
        url: '/plugin/list',
        method: 'get',
        params
    })
}

export function createPlugin(data) {
    return request({
        url: '/plugin',
        method: 'post',
        data
    })
}

export function updatePlugin(data) {
    return request({
        url: '/plugin',
        method: 'put',
        data
    })
}

export function deletePlugin(params) {
    return request({
        url: '/plugin',
        method: 'delete',
        params
    })
}

export function getPluginByID(params) {
    return request({
        url: '/plugin',
        method: 'get',
        params
    })
}

export function uploadPluginFile(data) {
    return request({
        url: '/plugin/file',
        method: 'put',
        data
    })
}

export function debugPlugin(data) {
    return request({
        url: '/plugin/debug',
        method: 'post',
        data
    })
}

export function getSelectList() {
    return request({
        url: '/plugin/selectList',
        method: 'get'
    })
}

export function checkPluginFileStatus(params) {
    return request({
        url: '/plugin/file/status',
        method: 'get',
        params
    })
}

export function publishPlugin(data) {
    return request({
        url: '/plugin/publish',
        method: 'post',
        data
    })
}
