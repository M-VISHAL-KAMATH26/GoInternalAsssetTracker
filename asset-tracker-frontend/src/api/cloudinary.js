import axios from 'axios'

const MAX_IMAGE_SIZE = 2 * 1024 * 1024

export const uploadProfileImage = async (file) => {
	if (!file) {
		return ''
	}
	if (!file.type.startsWith('image/')) {
		throw new Error('Profile picture must be an image file.')
	}
	if (file.size > MAX_IMAGE_SIZE) {
		throw new Error('Profile picture must be smaller than 2 MB.')
	}

	const cloudName = import.meta.env.VITE_CLOUDINARY_CLOUD_NAME
	const uploadPreset = import.meta.env.VITE_CLOUDINARY_UPLOAD_PRESET
	if (!cloudName || !uploadPreset) {
		throw new Error('Cloudinary is not configured. Add the cloud name and unsigned upload preset to the frontend environment.')
	}

	const body = new FormData()
	body.append('file', file)
	body.append('upload_preset', uploadPreset)

	const apiKey = import.meta.env.VITE_CLOUDINARY_API_KEY
	if (apiKey) {
		body.append('api_key', apiKey)
	}

	const { data } = await axios.post(
		`https://api.cloudinary.com/v1_1/${cloudName}/image/upload`,
		body,
	)
	return data.secure_url
}
