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

// The inventory API expects snake_case keys, so map the camelCase form
// fields the UI works with before sending them.
const toAssetPayload = ({ name, type, category, serialNumber, status }) => ({
	name,
	type,
	category,
	serial_number: serialNumber,
	...(status ? { status } : {}),
})

export const listAssets = () => inventoryClient.get('/assets')

export const listAssetCatalog = () => inventoryClient.get('/assets/catalog')

export const createAsset = (asset) => inventoryClient.post('/assets', toAssetPayload(asset))

export const updateAsset = (id, asset) => inventoryClient.put(`/assets/${id}`, toAssetPayload(asset))

export const retireAsset = (id) => inventoryClient.patch(`/assets/${id}/retire`)