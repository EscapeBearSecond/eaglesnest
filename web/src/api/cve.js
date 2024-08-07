import service from '@/utils/request'

export const getCveList = (data) => {
    return service({
        url: '/vuln/list',
        method: 'post',
        data
    })
}