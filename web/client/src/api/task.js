import request from '@/utils/axiosReq'


export function getRunningFlowInfo(params) {
    return request({
        url: '/task/running/flow/info',
        method: 'get',
        params
    })
}

export function getRunningMonitorInfo(params) {
    return request({
        url: '/task/running/monitor/info',
        method: 'get',
        params
    })
}

export function stopTask(data) {
    return request({
        url: '/task/stop',
        method: 'post',
        data
    })
}

export function resetQPS(data) {
    return request({
        url: '/task/reset/qps',
        method: 'post',
        data
    })
}

export function resetUsers(data) {
    return request({
        url: '/task/reset/users',
        method: 'post',
        data
    })
}

export function getTaskPipelineUsers(params) {
    return request({
        url: '/task/pipeline/users',
        method: 'get',
        params
    })
}

export function getTaskStatus(params) {
    return request({
        url: '/task/status',
        method: 'get',
        params
    })
}

export function getTaskWorks(params) {
    return request({
        url: '/task/works',
        method: 'get',
        params
    })
}

