import axios from 'axios'

const userServiceClient = axios.create({
	baseURL: import.meta.env.VITE_USER_SERVICE_URL,
})

export default userServiceClient
