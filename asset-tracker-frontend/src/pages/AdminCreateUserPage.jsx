import { useState } from 'react'
import { Link } from 'react-router-dom'
import { uploadProfileImage } from '../api/cloudinary'
import { createEmployee } from '../features/admin/adminApi'

const initialForm = {
	name: '',
	email: '',
	password: '',
	role: 'employee',
	officeLocation: '',
	managerEmail: '',
	managerName: '',
	managerPassword: '',
	managerOfficeLocation: '',
}

function AdminCreateUserPage() {
	const [form, setForm] = useState(initialForm)
	const [isSubmitting, setIsSubmitting] = useState(false)
	const [error, setError] = useState('')
	const [success, setSuccess] = useState('')
	const [profileImage, setProfileImage] = useState(null)
	const [managerProfileImage, setManagerProfileImage] = useState(null)
	const [fileInputKey, setFileInputKey] = useState(0)

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
			const [avatarURL, managerAvatarURL] = await Promise.all([
				uploadProfileImage(profileImage),
				form.role === 'admin' ? Promise.resolve('') : uploadProfileImage(managerProfileImage),
			])

			const employee = await createEmployee({
				name: form.name.trim(),
				email: form.email.trim(),
				password: form.password,
				role: form.role,
				office_location: form.officeLocation.trim(),
				avatar_url: avatarURL,
				manager_email: form.role === 'admin' ? '' : form.managerEmail.trim(),
				manager_name: form.role === 'admin' ? '' : form.managerName.trim(),
				manager_password: form.role === 'admin' ? '' : form.managerPassword,
				manager_office_location: form.role === 'admin' ? '' : form.managerOfficeLocation.trim(),
				manager_avatar_url: form.role === 'admin' ? '' : managerAvatarURL,
			})
			setSuccess(`${employee.name} was created and can now sign in.`)
			setForm(initialForm)
			setProfileImage(null)
			setManagerProfileImage(null)
			setFileInputKey((current) => current + 1)
		} catch (requestError) {
			setError(requestError.response?.data?.error ?? requestError.message ?? 'Unable to create the user. Please try again.')
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
					<label className="text-sm font-medium text-slate-700">Office location
						<input className="mt-2 w-full rounded-lg border px-3 py-2.5" name="officeLocation" onChange={updateField} placeholder="Bangalore" required value={form.officeLocation} />
					</label>
					<label className="text-sm font-medium text-slate-700">Profile picture <span className="font-normal text-slate-400">(optional)</span>
						<input
							accept="image/*"
							className="mt-2 block w-full text-sm text-slate-600 file:mr-3 file:rounded-lg file:border-0 file:bg-indigo-50 file:px-3 file:py-2 file:text-indigo-700"
							key={`employee-${fileInputKey}`}
							onChange={(event) => setProfileImage(event.target.files?.[0] ?? null)}
							type="file"
						/>
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
							<label className="text-sm font-medium text-slate-700">New manager office <span className="font-normal text-slate-400">(if new)</span>
								<input className="mt-2 w-full rounded-lg border bg-white px-3 py-2.5" name="managerOfficeLocation" onChange={updateField} placeholder="Mumbai" value={form.managerOfficeLocation} />
							</label>
							<label className="text-sm font-medium text-slate-700">New manager picture <span className="font-normal text-slate-400">(optional)</span>
								<input
									accept="image/*"
									className="mt-2 block w-full text-sm text-slate-600 file:mr-3 file:rounded-lg file:border-0 file:bg-white file:px-3 file:py-2 file:text-indigo-700"
									key={`manager-${fileInputKey}`}
									onChange={(event) => setManagerProfileImage(event.target.files?.[0] ?? null)}
									type="file"
								/>
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
