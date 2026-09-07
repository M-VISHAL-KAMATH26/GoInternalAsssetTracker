import { createSlice } from '@reduxjs/toolkit'

const TOKEN_STORAGE_KEY = 'token'

const decodeJwt = (token) => {
	try {
		const payload = token.split('.')[1]

		if (!payload) {
			return null
		}

		const base64 = payload.replace(/-/g, '+').replace(/_/g, '/')
		const paddedBase64 = base64.padEnd(Math.ceil(base64.length / 4) * 4, '=')

		return JSON.parse(atob(paddedBase64))
	} catch {
		return null
	}
}

const getAuthState = () => {
	const token = localStorage.getItem(TOKEN_STORAGE_KEY)
	const claims = token ? decodeJwt(token) : null

	return {
		token,
		role: claims?.role ?? null,
		employeeId: claims?.employeeId ?? claims?.employee_id ?? null,
	}
}

const authSlice = createSlice({
	name: 'auth',
	initialState: getAuthState(),
	reducers: {
		setCredentials: (state, action) => {
			const token = action.payload
			const claims = decodeJwt(token)

			state.token = token
			state.role = claims?.role ?? null
			state.employeeId = claims?.employeeId ?? claims?.employee_id ?? null
			localStorage.setItem(TOKEN_STORAGE_KEY, token)
		},
		logout: (state) => {
			state.token = null
			state.role = null
			state.employeeId = null
			localStorage.removeItem(TOKEN_STORAGE_KEY)
		},
	},
})

export const { setCredentials, logout } = authSlice.actions
export default authSlice.reducer
