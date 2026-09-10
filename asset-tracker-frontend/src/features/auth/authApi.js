import userServiceClient from '../../api/userServiceClient'

export const login = async (email, password) => {
	const { data } = await userServiceClient.post('/login', { email, password })
	return data
}
