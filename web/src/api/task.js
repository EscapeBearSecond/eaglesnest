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