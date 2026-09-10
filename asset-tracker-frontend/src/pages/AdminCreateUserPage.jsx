import { useState } from 'react'
import { Link } from 'react-router-dom'
import { createEmployee } from '../features/admin/adminApi'

const initialForm = {
	name: '',
	email: '',
	password: '',
	role: 'employee',
	managerEmail: '',
	managerName: '',
	managerPassword: '',
}

function AdminCreateUserPage() {
	const [form, setForm] = useState(initialForm)
	const [isSubmitting, setIsSubmitting] = useState(false)
	const [error, setError] = useState('')
	const [success, setSuccess] = useState('')

	const updateField = (event) => {
		const { name, value } = event.target
		setForm((current) => ({ ...current, [name]: value }))
		setError('')
		setSuccess('')
	}

	const handleSubmit = async (event) => {
		event.preventDefault()
		setIsSubmitting(true)
		setError('')
		setSuccess('')

		try {
			const employee = await createEmployee({
				name: form.name.trim(),
				email: form.email.trim(),
				password: form.password,
				role: form.role,
				manager_email: form.role === 'admin' ? '' : form.managerEmail.trim(),
				manager_name: form.role === 'admin' ? '' : form.managerName.trim(),
				manager_password: form.role === 'admin' ? '' : form.managerPassword,
			})
			setSuccess(`${employee.name} was created and can now sign in.`)
			setForm(initialForm)
		} catch (requestError) {
			setError(requestError.response?.data?.error ?? 'Unable to create the user. Please try again.')
		} finally {
			setIsSubmitting(false)
		}
	}

	const needsManager = form.role !== 'admin'

	return (
		<div className="mx-auto max-w-2xl">
			<div className="mb-6">
				<Link className="text-sm font-medium text-indigo-600 hover:text-indigo-700" to="/admin">← Back to admin panel</Link>
				<h1 className="mt-3 text-3xl font-bold tracking-tight">Add a new user</h1>
				<p className="mt-2 text-sm text-slate-600">Create an employee, manager, or another administrator.</p>
			</div>

			<form className="space-y-6 rounded-xl border border-slate-200 bg-white p-6 shadow-sm sm:p-8" onSubmit={handleSubmit}>
				<div className="grid gap-5 sm:grid-cols-2">
					<label className="text-sm font-medium text-slate-700">Name
						<input className="mt-2 w-full rounded-lg border px-3 py-2.5" name="name" onChange={updateField} required value={form.name} />
					</label>
					<label className="text-sm font-medium text-slate-700">Role
						<select className="mt-2 w-full rounded-lg border bg-white px-3 py-2.5" name="role" onChange={updateField} value={form.role}>
							<option value="employee">Employee</option>
							<option value="manager">Manager</option>
							<option value="admin">Admin</option>
						</select>
					</label>
					<label className="text-sm font-medium text-slate-700">Email
						<input autoComplete="email" className="mt-2 w-full rounded-lg border px-3 py-2.5" name="email" onChange={updateField} required type="email" value={form.email} />
					</label>
					<label className="text-sm font-medium text-slate-700">Temporary password
						<input autoComplete="new-password" className="mt-2 w-full rounded-lg border px-3 py-2.5" minLength={8} name="password" onChange={updateField} required type="password" value={form.password} />
					</label>
				</div>

				{needsManager && (
					<fieldset className="space-y-4 rounded-lg border border-slate-200 bg-slate-50 p-5">
						<legend className="px-2 text-sm font-semibold text-slate-800">Approving manager</legend>
						<p className="text-sm text-slate-600">Enter an existing manager’s email. If the manager is new, also provide their name and temporary password; the new manager will report to the current admin.</p>
						<label className="block text-sm font-medium text-slate-700">Manager email
							<input className="mt-2 w-full rounded-lg border bg-white px-3 py-2.5" name="managerEmail" onChange={updateField} required type="email" value={form.managerEmail} />
						</label>
						<div className="grid gap-5 sm:grid-cols-2">
							<label className="text-sm font-medium text-slate-700">New manager name <span className="font-normal text-slate-400">(if new)</span>
								<input className="mt-2 w-full rounded-lg border bg-white px-3 py-2.5" name="managerName" onChange={updateField} value={form.managerName} />
							</label>
							<label className="text-sm font-medium text-slate-700">New manager password <span className="font-normal text-slate-400">(if new)</span>
								<input className="mt-2 w-full rounded-lg border bg-white px-3 py-2.5" minLength={8} name="managerPassword" onChange={updateField} type="password" value={form.managerPassword} />
							</label>
						</div>
					</fieldset>
				)}

				{error && <p className="text-sm text-red-600" role="alert">{error}</p>}
				{success && <p className="text-sm text-emerald-700" role="status">{success}</p>}

				<div className="flex justify-end">
					<button className="rounded-lg bg-indigo-600 px-5 py-2.5 text-sm font-medium text-white hover:bg-indigo-700 disabled:opacity-60" disabled={isSubmitting} type="submit">
						{isSubmitting ? 'Creating user...' : 'Create user'}
					</button>
				</div>
			</form>
		</div>
	)
}

export default AdminCreateUserPage
