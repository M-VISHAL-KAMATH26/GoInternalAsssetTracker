import axiosClient from '../../api/axios'
import userServiceClient from '../../api/userServiceClient'

const authConfig = () => ({
	headers: {
		Authorization: `Bearer ${localStorage.getItem('token')}`,
	},
})

export const listAllRequests = async () => {
	const { data } = await axiosClient.get('/requests/all')
	return data
}

export const createEmployee = async (employee) => {
	const { data } = await userServiceClient.post('/employees', employee, authConfig())
	return data
}

export const listEmployees = async () => {
	const { data } = await userServiceClient.get('/employees', authConfig())
	return data
}
