import { useState } from 'react'
import { useDispatch } from 'react-redux'
import { Link, useNavigate } from 'react-router-dom'
import { login } from '../features/auth/authApi'
import { setCredentials } from '../features/auth/authSlice'

function LoginPage() {
	const dispatch = useDispatch()
	const navigate = useNavigate()
	const [email, setEmail] = useState('')
	const [password, setPassword] = useState('')
	const [error, setError] = useState('')
	const [isSubmitting, setIsSubmitting] = useState(false)

	const handleSubmit = async (event) => {
		event.preventDefault()
		setError('')
		setIsSubmitting(true)

		try {
			const data = await login(email.trim(), password)
			dispatch(setCredentials(data.token))
			navigate(data.role === 'admin' ? '/admin' : '/requests')
		} catch {
			setError('Invalid email or password.')
		} finally {
			setIsSubmitting(false)
		}
	}

	return (
		<main className="flex min-h-screen items-center justify-center bg-slate-50 px-4 py-8 text-slate-900 sm:px-6">
			<div className="w-full max-w-md">
				<div className="mb-6">
					<p className="text-sm font-medium text-slate-500">Asset Tracker</p>
					<h1 className="mt-1 text-3xl font-semibold tracking-tight">Sign in</h1>
					<p className="mt-2 text-sm text-slate-600">
						Use your work email and password to continue.
					</p>
					<Link className="mt-3 inline-block text-sm font-medium text-slate-500 hover:text-slate-900" to="/">
						← Back to home
					</Link>
				</div>

				<form
					className="space-y-6 rounded-xl border border-slate-200 bg-white p-6 shadow-sm sm:p-8"
					noValidate
					onSubmit={handleSubmit}
				>
					<div>
						<label className="mb-2 block text-sm font-medium text-slate-700" htmlFor="email">
							Email
						</label>
						<input
							autoComplete="email"
							className="w-full rounded-lg border border-slate-300 px-3 py-2.5 text-sm outline-none transition placeholder:text-slate-400 focus:border-slate-500 focus:ring-2 focus:ring-slate-200"
							id="email"
							name="email"
							onChange={(event) => {
								setEmail(event.target.value)
								setError('')
							}}
							placeholder="you@company.com"
							required
							type="email"
							value={email}
						/>
					</div>

					<div>
						<label className="mb-2 block text-sm font-medium text-slate-700" htmlFor="password">
							Password
						</label>
						<input
							autoComplete="current-password"
							className="w-full rounded-lg border border-slate-300 px-3 py-2.5 text-sm outline-none transition placeholder:text-slate-400 focus:border-slate-500 focus:ring-2 focus:ring-slate-200"
							id="password"
							name="password"
							onChange={(event) => {
								setPassword(event.target.value)
								setError('')
							}}
							required
							type="password"
							value={password}
						/>
					</div>

					{error && (
						<p className="text-sm text-red-600" role="alert">
							{error}
						</p>
					)}

					<button
						className="w-full rounded-lg bg-slate-900 px-4 py-2.5 text-sm font-medium text-white transition hover:bg-slate-700 disabled:cursor-not-allowed disabled:opacity-60"
						disabled={isSubmitting}
						type="submit"
					>
						{isSubmitting ? 'Signing in...' : 'Sign in'}
					</button>
				</form>
			</div>
		</main>
	)
}

export default LoginPage
