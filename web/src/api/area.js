import service from '@/utils/request'

export const createArea = (data) => {
    return service({
        url: '/area',
        method: 'post',
        data
    })
}

export const editArea = (data) => {
    return service({
      url: '/area',
      method: 'put',
      data
    })
}

export const getAreaList = (data) => {
    return service({
      url: '/area/list',
      method: 'post',
      data
    })
}