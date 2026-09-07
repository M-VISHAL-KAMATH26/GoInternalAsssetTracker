import axiosClient from '../../api/axios'

export const createRequest = (request) => axiosClient.post('/requests', request)

export const listMyRequests = () => axiosClient.get('/requests')

export const getRequest = (id) => axiosClient.get(`/requests/${id}`)