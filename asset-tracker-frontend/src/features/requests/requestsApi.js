import axiosClient from '../../api/axios'

export const createRequest = ({ assetType, category, justification }) =>
	axiosClient.post('/requests', {
		asset_type: assetType,
		category,
		justification,
	})

export const listMyRequests = () => axiosClient.get('/requests')

export const getRequest = (id) => axiosClient.get(`/requests/${id}`)