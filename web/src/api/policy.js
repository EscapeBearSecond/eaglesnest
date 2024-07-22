import service from '@/utils/request'

export const getPolicyList = (data) => {
    return service({
      url: '/policy/list',
      method: 'post',
      data
    })
}


export const createPolicy = (data) => {
  return service({
    url: '/policy',
    method: 'get',
    data
  })
}