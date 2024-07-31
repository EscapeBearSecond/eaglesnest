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
    method: 'post',
    data
  })
}

export const updatePolicy = (data) => {
  return service({
    url: '/policy',
    method: 'put',
    data
  })
}

export const deletePolicy = (data) => {
  return service({
    url: `/policy/${data.id}`,
    method: 'delete',
    data
  })
}

export const getPolicyId = (param) => {
  return service({
    url: `/policy/${param.id}`,
    method: 'get',
    param
  })
}
