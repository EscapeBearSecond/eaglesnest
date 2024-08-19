import service from '@/utils/request'

export const taaskStatistics = (param) => {
    return service({
        url: '/task/statistics',
        method: 'get',
        param
    })
}

export const vulnStatistics = (param) => {
    return service({
        url: '/vuln/statistics',
        method: 'get',
        param
    })
}

export const getVulnDCommon = (param) => {
    return service({
        url: '/vuln/common',
        method: 'get',
        param
    })
}

export const getAssethighrisk = (param) => {
    return service({
        url: '/asset/highrisk',
        method: 'get',
        param
    })
}
