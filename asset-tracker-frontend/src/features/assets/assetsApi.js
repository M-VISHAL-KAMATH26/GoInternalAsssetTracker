import axios from 'axios'

const inventoryClient = axios.create({
	baseURL: import.meta.env.VITE_INVENTORY_SERVICE_URL,
})

inventoryClient.interceptors.request.use((config) => {
	const token = localStorage.getItem('token')

	if (token) {
		config.headers.Authorization = `Bearer ${token}`
	}

	return config
})

export const listAssets = () => inventoryClient.get('/assets')

export const createAsset = (asset) => inventoryClient.post('/assets', asset)

export const updateAsset = (id, asset) => inventoryClient.put(`/assets/${id}`, asset)

export const retireAsset = (id) => inventoryClient.patch(`/assets/${id}/retire`)