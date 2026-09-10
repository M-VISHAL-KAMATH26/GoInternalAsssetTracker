import userServiceClient from '../../api/userServiceClient'

export const login = async (email, password) => {
	const { data } = await userServiceClient.post('/login', { email, password })
	return data
}

export const getMyProfile = async () => {
	const { data } = await userServiceClient.get('/employees/me', {
		headers: {
			Authorization: `Bearer ${localStorage.getItem('token')}`,
		},
	})
	return data
}
