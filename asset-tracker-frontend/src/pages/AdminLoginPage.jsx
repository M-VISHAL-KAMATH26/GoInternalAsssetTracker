import { useState } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { Link, Navigate, useNavigate } from 'react-router-dom'
import { login } from '../features/auth/authApi'
import { setCredentials } from '../features/auth/authSlice'

function AdminLoginPage() {
	const dispatch = useDispatch()
	const navigate = useNavigate()
	const { role, token } = useSelector((state) => state.auth)
	const [email, setEmail] = useState('')
	const [password, setPassword] = useState('')
	const [error, setError] = useState('')
	const [isSubmitting, setIsSubmitting] = useState(false)

	if (token && role === 'admin') {
		return <Navigate replace to="/admin" />
	}

	const handleSubmit = async (event) => {
		event.preventDefault()
		setError('')
		setIsSubmitting(true)

		try {
			const data = await login(email.trim(), password)
			if (data.role !== 'admin') {
				setError('This account does not have administrator access.')
				return
			}
			dispatch(setCredentials(data.token))
			navigate('/admin')
		} catch {
			setError('Invalid admin email or password.')
		} finally {
			setIsSubmitting(false)
		}
	}

	return (
		<main className="flex min-h-screen items-center justify-center bg-slate-950 px-4 py-8">
			<div className="w-full max-w-md">
				<div className="mb-6 text-white">
					<p className="text-sm font-medium text-indigo-300">Asset Tracker Administration</p>
					<h1 className="mt-1 text-3xl font-semibold tracking-tight">Admin sign in</h1>
					<p className="mt-2 text-sm text-slate-300">Only administrator accounts can access this panel.</p>
					<Link className="mt-3 inline-block text-sm font-medium text-slate-400 hover:text-white" to="/">
						← Back to home
					</Link>
				</div>
				<form className="space-y-5 rounded-xl bg-white p-6 shadow-xl sm:p-8" onSubmit={handleSubmit}>
					<label className="block text-sm font-medium text-slate-700" htmlFor="admin-email">
						Admin email
						<input
							autoComplete="email"
							className="mt-2 w-full rounded-lg border border-slate-300 px-3 py-2.5 outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-100"
							id="admin-email"
							onChange={(event) => setEmail(event.target.value)}
							required
							type="email"
							value={email}
						/>
					</label>
					<label className="block text-sm font-medium text-slate-700" htmlFor="admin-password">
						Password
						<input
							autoComplete="current-password"
							className="mt-2 w-full rounded-lg border border-slate-300 px-3 py-2.5 outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-100"
							id="admin-password"
							onChange={(event) => setPassword(event.target.value)}
							required
							type="password"
							value={password}
						/>
					</label>
					{error && <p className="text-sm text-red-600" role="alert">{error}</p>}
					<button
						className="w-full rounded-lg bg-indigo-600 px-4 py-2.5 font-medium text-white hover:bg-indigo-700 disabled:opacity-60"
						disabled={isSubmitting}
						type="submit"
					>
						{isSubmitting ? 'Signing in...' : 'Open admin panel'}
					</button>
				</form>
			</div>
		</main>
	)
}

export default AdminLoginPage
