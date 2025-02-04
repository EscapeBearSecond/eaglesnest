import service from '@/utils/request'

export const createTask = (data) => {
    return service({
        url: '/task',
        method: 'post',
        data
    })
}

export const stopTask = (data) => {
    return service({
      url: `/task/stop/${data.id}`,
      method: 'post',
      data
    })
}

export const getTaskList = (data) => {
    return service({
      url: '/task/list',
      method: 'post',
      data
    })
}

export const reportTask = (data) => {
    return service({
      url: '/task/report',
      method: 'post',
      data,
      responseType: 'blob',
      headers: {'Content-Type': 'multipart/form-data'}
    })
}

export const reportTaskDoc = (data) => {
    return service({
      url: '/task/docs',
      method: 'post',
      data,
      responseType: 'blob',
      headers: {'Content-Type': 'multipart/form-data'}
    })
}

export const startTask = (data) => {
    return service({
      url: `/task/execute/${data.id}`,
      method: 'post',
      data
    })
}

export const delTask = (data) => {
  return service({
    url: `/task/${data.id}`,
    method: 'delete',
    data
  })
}


export const getTask = (param) => {
  return service({
    url: `/task/${param.id}`,
    method: 'get',
    param
  })
}
export const getTaskStage = (param) => {
  return service({
    url: `/task/stage/${param.id}`,
    method: 'get',
    param
  })
}