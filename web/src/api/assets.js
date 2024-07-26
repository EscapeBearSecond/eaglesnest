import service from '@/utils/request'

export const createApi = (data) => {
    return service({
        url: '/asset/add',
        method: 'post',
        data
    })
}

export const editApi = (data) => {
    return service({
      url: '/asset',
      method: 'put',
      data
    })
}

export const getListApi = (data) => {
    return service({
      url: '/asset/list',
      method: 'post',
      data
    })
}

export const delApi = (data) => {
  return service({
    url: `/asset/${data.id}`,
    method: 'del',
    data
  })
}