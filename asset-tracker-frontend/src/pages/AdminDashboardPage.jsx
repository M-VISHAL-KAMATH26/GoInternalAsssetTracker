import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { listAllRequests, listEmployees } from '../features/admin/adminApi'

const statusStyles = {
	pending: 'bg-amber-100 text-amber-700',
	approved: 'bg-emerald-100 text-emerald-700',
	rejected: 'bg-red-100 text-red-700',
	fulfilled: 'bg-blue-100 text-blue-700',
}

const formatDate = (value) => value
	? new Intl.DateTimeFormat(undefined, { dateStyle: 'medium' }).format(new Date(value))
	: '-'

function AdminDashboardPage() {
	const [requests, setRequests] = useState([])
	const [employees, setEmployees] = useState([])
	const [isLoading, setIsLoading] = useState(true)
	const [error, setError] = useState('')

	useEffect(() => {
		let mounted = true

		Promise.all([listAllRequests(), listEmployees()])
			.then(([requestData, employeeData]) => {
				if (mounted) {
					setRequests(Array.isArray(requestData) ? requestData : [])
					setEmployees(Array.isArray(employeeData) ? employeeData : [])
				}
			})
			.catch(() => mounted && setError('Unable to load admin data. Make sure user-service and request-service are running.'))
			.finally(() => mounted && setIsLoading(false))

		return () => {
			mounted = false
		}
	}, [])

	const employeeNames = new Map(employees.map((employee) => [employee.id, employee.name]))
	const pendingCount = requests.filter((request) => request.status === 'pending').length

	return (
		<div className="space-y-8">
			<header className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
				<div>
					<p className="text-sm font-medium text-indigo-600">Administration</p>
					<h1 className="mt-1 text-3xl font-bold tracking-tight">Admin panel</h1>
					<p className="mt-2 text-sm text-slate-600">Manage users and monitor asset requests.</p>
				</div>
				<Link className="rounded-lg bg-indigo-600 px-4 py-2.5 text-sm font-medium text-white hover:bg-indigo-700" to="/admin/users/new">
					Add new user
				</Link>
			</header>

			<section className="grid gap-4 sm:grid-cols-3">
				<div className="rounded-xl border bg-white p-5 shadow-sm"><p className="text-sm text-slate-500">Total requests</p><p className="mt-2 text-3xl font-bold">{requests.length}</p></div>
				<div className="rounded-xl border bg-white p-5 shadow-sm"><p className="text-sm text-slate-500">Pending requests</p><p className="mt-2 text-3xl font-bold">{pendingCount}</p></div>
				<div className="rounded-xl border bg-white p-5 shadow-sm"><p className="text-sm text-slate-500">Registered users</p><p className="mt-2 text-3xl font-bold">{employees.length}</p></div>
			</section>

			<section className="overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm">
				<div className="border-b px-5 py-4"><h2 className="font-semibold">All requests</h2></div>
				{isLoading && <p className="p-8 text-center text-sm text-slate-500">Loading requests...</p>}
				{!isLoading && error && <p className="p-8 text-center text-sm text-red-600">{error}</p>}
				{!isLoading && !error && requests.length === 0 && <p className="p-8 text-center text-sm text-slate-500">No requests have been created.</p>}
				{!isLoading && !error && requests.length > 0 && (
					<div className="overflow-x-auto">
						<table className="w-full text-left text-sm">
							<thead className="bg-slate-50 text-xs uppercase tracking-wide text-slate-500">
								<tr><th className="px-5 py-3">Employee</th><th className="px-5 py-3">Asset</th><th className="px-5 py-3">Category</th><th className="px-5 py-3">Status</th><th className="px-5 py-3">Created</th></tr>
							</thead>
							<tbody className="divide-y divide-slate-100">
								{requests.map((request) => {
									const employeeID = request.employee_id ?? request.employeeId
									const status = request.status?.toLowerCase()
									return (
										<tr key={request.id}>
											<td className="px-5 py-4"><p className="font-medium">{employeeNames.get(employeeID) ?? 'Unknown user'}</p><p className="text-xs text-slate-400">{employeeID}</p></td>
											<td className="px-5 py-4">{request.asset_type ?? request.assetType}</td>
											<td className="px-5 py-4">{request.category}</td>
											<td className="px-5 py-4"><span className={`rounded-full px-2.5 py-1 text-xs font-medium ${statusStyles[status] ?? 'bg-slate-100 text-slate-600'}`}>{status}</span></td>
											<td className="px-5 py-4 text-slate-500">{formatDate(request.created_at ?? request.createdAt)}</td>
										</tr>
									)
								})}
							</tbody>
						</table>
					</div>
				)}
			</section>
		</div>
	)
}

export default AdminDashboardPage
