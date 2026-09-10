import axiosClient from '../../api/axios'

export const listPendingApprovals = () => axiosClient.get('/approvals/pending')

export const approveRequest = (id, comment = '') =>
	axiosClient.patch(`/requests/${id}/approve`, comment ? { comment } : {})

export const rejectRequest = (id, comment = '') =>
	axiosClient.patch(`/requests/${id}/reject`, comment ? { comment } : {})