import axios from 'axios'

const axiosClient = axios.create({
	baseURL: import.meta.env.VITE_REQUEST_SERVICE_URL,
})

axiosClient.interceptors.request.use((config) => {
	const token = localStorage.getItem('token')

	if (token) {
		config.headers.Authorization = `Bearer ${token}`
	}

	return config
})

export default axiosClient
