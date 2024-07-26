import service from '@/utils/request'

export const createApi = (data) => {
    return service({
        url: '/scan/add',
        method: 'post',
        data
    })
}

export const editApi = (data) => {
    return service({
      url: '/scan',
      method: 'put',
      data
    })
}

export const getListApi = (data) => {
    return service({
      url: '/scan/list',
      method: 'post',
      data
    })
}

export const delApi = (data) => {
  return service({
    url: `/scan/${data.id}`,
    method: 'del',
    data
  })
}