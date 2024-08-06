import service from '@/utils/request'

export const createTemplate = (data) => {
    return service({
        url: '/template',
        method: 'post',
        data
    })
}

export const getTemplateList = (data) => {
    return service({
      url: '/template/list',
      method: 'post',
      data
    })
}
export const postTemplateImports = (data) => {
    return service({
      url: '/template/imports',
      method: 'post',
      data
    })
}

export const getTemplateTagList = (params) => {
    return service({
      url: '/template/tags',
      method: 'get',
      params
    })
}


export const getTemplate = (param) => {
    return service({
      url: `/template/${param.id}`,
      method: 'get',
      param
    })
}

export const updateTemplate = (data) => {
    return service({
      url: `/template`,
      method: 'put',
      data
    })
}

export const delTemplate = (data) => {
  return service({
    url: `/template/${data.id}`,
    method: 'delete',
    data
  })
}

