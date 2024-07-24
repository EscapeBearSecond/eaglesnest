import service from '@/utils/request'

export const createArea = (data) => {
    return service({
      url: '/area',
      method: 'post',
      data
    })
  }