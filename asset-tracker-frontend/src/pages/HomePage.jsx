import {
	Boxes,
	CheckCircle,
	ClipboardList,
	Mail,
	MapPin,
	Package,
	Plus,
	ShieldCheck,
	UserRound,
	Users,
} from 'lucide-react'
import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { useSelector } from 'react-redux'
import { listAssets } from '../features/assets/assetsApi'
import { getMyProfile } from '../features/auth/authApi'
import { listMyRequests } from '../features/requests/requestsApi'

const getList = (data, key) => {
	if (Array.isArray(data)) {
		return data
	}

	return data?.[key] ?? data?.items ?? []
}

const getStatus = (request) => (request.status ?? 'pending').toLowerCase()

const statusStyles = {
	pending: 'bg-amber-100 text-amber-700',
	approved: 'bg-emerald-100 text-emerald-700',
	rejected: 'bg-red-100 text-red-700',
	fulfilled: 'bg-blue-100 text-blue-700',
}

const formatDate = (value) => {
	if (!value) {
		return '-'
	}

	const date = new Date(value)

	return Number.isNaN(date.getTime())
		? '-'
		: new Intl.DateTimeFormat(undefined, {
				month: 'short',
				day: 'numeric',
				year: 'numeric',
		  }).format(date)
}

const getCreatedDate = (request) => request.createdAt ?? request.created_at

const sortRecent = (requests) => [...requests]
	.sort((first, second) => {
		const firstDate = new Date(getCreatedDate(first) ?? 0).getTime()
		const secondDate = new Date(getCreatedDate(second) ?? 0).getTime()

		return secondDate - firstDate
	})
	.slice(0, 5)

function StatCard({ icon: Icon, label, value, iconClassName }) {
	return (
		<div className={`rounded-xl border-l-4 bg-white p-5 shadow-sm ${iconClassName}`}>
			<div className="flex items-center justify-between">
				<Icon aria-hidden="true" className="h-5 w-5 text-slate-500" />
				<span className="text-3xl font-bold tracking-tight text-slate-900">{value}</span>
			</div>
			<p className="mt-3 text-sm font-medium text-slate-500">{label}</p>
		</div>
	)
}

function QuickAction({ description, icon: Icon, title, to }) {
	return (
		<Link
			className="group rounded-xl border border-slate-200 bg-white p-5 shadow-sm transition hover:-translate-y-0.5 hover:border-indigo-200 hover:shadow-md"
			to={to}
		>
			<div className="flex h-10 w-10 items-center justify-center rounded-lg bg-indigo-50 text-indigo-600">
				<Icon aria-hidden="true" className="h-5 w-5" />
			</div>
			<h3 className="mt-4 font-semibold text-slate-900 group-hover:text-indigo-600">{title}</h3>
			<p className="mt-1 text-sm leading-5 text-slate-500">{description}</p>
		</Link>
	)
}

function DashboardSkeleton() {
	return (
		<div aria-label="Loading dashboard" className="animate-pulse space-y-8" role="status">
			<div className="h-10 w-64 rounded bg-slate-200" />
			<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-5">
				{[1, 2, 3, 4, 5].map((item) => <div className="h-32 rounded-xl bg-slate-200" key={item} />)}
			</div>
			<div className="h-44 rounded-xl bg-slate-200" />
			<div className="h-56 rounded-xl bg-slate-200" />
		</div>
	)
}

function ProfileCard({ profile, role }) {
	const name = profile?.name ?? 'User'
	const initials = name
		.split(/\s+/)
		.filter(Boolean)
		.slice(0, 2)
		.map((part) => part[0]?.toUpperCase())
		.join('')

	return (
		<aside className="h-fit rounded-xl border border-slate-200 bg-white p-6 text-center shadow-sm lg:sticky lg:top-8">
			{profile?.avatar_url ? (
				<img
					alt={`${name}'s profile`}
					className="mx-auto h-28 w-28 rounded-full object-cover ring-4 ring-indigo-50"
					src={profile.avatar_url}
				/>
			) : (
				<div className="mx-auto flex h-28 w-28 items-center justify-center rounded-full bg-indigo-100 text-3xl font-bold text-indigo-700 ring-4 ring-indigo-50">
					{initials || <UserRound className="h-10 w-10" />}
				</div>
			)}
			<h2 className="mt-4 text-xl font-semibold text-slate-900">{name}</h2>
			<span className="mt-2 inline-flex rounded-full bg-indigo-50 px-3 py-1 text-xs font-semibold uppercase tracking-wide text-indigo-700">
				{profile?.role ?? role}
			</span>
			<div className="mt-5 space-y-3 border-t border-slate-100 pt-5 text-left text-sm text-slate-600">
				<p className="flex items-start gap-2"><Mail className="mt-0.5 h-4 w-4 shrink-0 text-slate-400" /><span className="break-all">{profile?.email ?? 'Email not available'}</span></p>
				<p className="flex items-start gap-2"><MapPin className="mt-0.5 h-4 w-4 shrink-0 text-slate-400" /><span>{profile?.office_location || 'Location not set'}</span></p>
			</div>
		</aside>
	)
}

function HomePage() {
	const role = useSelector((state) => state.auth.role)?.toLowerCase() ?? 'employee'
	const [requests, setRequests] = useState([])
	const [assets, setAssets] = useState([])
	const [profile, setProfile] = useState(null)
	const [isLoading, setIsLoading] = useState(true)
	const [error, setError] = useState('')

	useEffect(() => {
		let isMounted = true

		const loadDashboard = async () => {
			try {
				const [requestsResponse, profileData] = await Promise.all([
					listMyRequests(),
					getMyProfile(),
				])
				const requestList = getList(requestsResponse.data, 'requests')

				if (isMounted) {
					setRequests(requestList)
					setProfile(profileData)
				}

				if (role === 'admin') {
					const assetsResponse = await listAssets()

					if (isMounted) {
						setAssets(getList(assetsResponse.data, 'assets'))
					}
				}
			} catch {
				if (isMounted) {
					setError('Some dashboard data could not be loaded. Please refresh and try again.')
				}
			} finally {
				if (isMounted) {
					setIsLoading(false)
				}
			}
		}

		loadDashboard()

		return () => {
			isMounted = false
		}
	}, [role])

	const pendingCount = requests.filter((request) => getStatus(request) === 'pending').length
	// TODO: Replace this placeholder count with a manager-specific pending endpoint.
	const pendingApprovals = role === 'manager' || role === 'admin' ? pendingCount : 0
	const quantityOf = (asset) => {
		const quantity = Number(asset.quantity)
		return Number.isFinite(quantity) ? quantity : 0
	}
	const totalAssetUnits = assets.reduce((sum, asset) => sum + quantityOf(asset), 0)
	const availableAssets = assets
		.filter((asset) => (asset.status ?? '').toLowerCase() === 'available')
		.reduce((sum, asset) => sum + quantityOf(asset), 0)
	const recentRequests = sortRecent(requests)

	const stats = [
		{ icon: ClipboardList, label: 'My Requests', value: requests.length, iconClassName: 'border-indigo-500' },
		{ icon: CheckCircle, label: 'Pending', value: pendingCount, iconClassName: 'border-amber-500' },
	]

	if (role === 'manager' || role === 'admin') {
		stats.push({ icon: Users, label: 'Pending Approvals', value: pendingApprovals, iconClassName: 'border-blue-500' })
	}

	if (role === 'admin') {
		stats.push(
			{ icon: Boxes, label: 'Total Assets', value: totalAssetUnits, iconClassName: 'border-indigo-500' },
			{ icon: Package, label: 'Available Assets', value: availableAssets, iconClassName: 'border-emerald-500' },
		)
	}

	const quickActions = [
		{ icon: Plus, title: 'New Request', description: 'Request an asset for your work.', to: '/requests/new' },
		{ icon: ClipboardList, title: 'My Requests', description: 'Track your submitted requests.', to: '/requests/history' },
	]

	if (role === 'manager' || role === 'admin') {
		quickActions.push({ icon: ShieldCheck, title: 'Review Approvals', description: 'Review requests waiting for a decision.', to: '/approvals' })
	}

	if (role === 'admin') {
		quickActions.push({ icon: Package, title: 'Manage Assets', description: 'View and maintain inventory.', to: '/assets' })
	}

	return (
		<div className="grid items-start gap-8 lg:grid-cols-[minmax(0,1fr)_280px]">
			<div className="space-y-8">
				{isLoading ? <DashboardSkeleton /> : (
					<>
					<header>
						<p className="text-sm font-medium text-slate-500">Dashboard</p>
						<h1 className="mt-1 text-3xl font-bold tracking-tight text-slate-900">Welcome back, {profile?.name ?? role}</h1>
						{error && <p className="mt-2 text-sm text-red-600" role="alert">{error}</p>}
					</header>

					<section aria-label="Request and asset statistics" className="grid gap-4 sm:grid-cols-2 lg:grid-cols-5">
						{stats.map((stat) => <StatCard {...stat} key={stat.label} />)}
					</section>

					<section>
						<div className="mb-4">
							<h2 className="text-xl font-semibold text-slate-900">Quick Actions</h2>
							<p className="mt-1 text-sm text-slate-500">Jump into the workflows you use most.</p>
						</div>
						<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
							{quickActions.map((action) => <QuickAction {...action} key={action.title} />)}
						</div>
					</section>

					<section className="rounded-xl border border-slate-200 bg-white shadow-sm">
						<div className="flex items-center justify-between border-b border-slate-100 px-5 py-4">
							<h2 className="font-semibold text-slate-900">Recent Activity</h2>
							<Link className="text-sm font-medium text-indigo-600 hover:text-indigo-700" to="/requests/history">View all</Link>
						</div>
						{recentRequests.length === 0 ? (
							<p className="px-5 py-8 text-center text-sm text-slate-500">No recent requests.</p>
						) : (
							<ul className="divide-y divide-slate-100">
								{recentRequests.map((request) => {
									const status = getStatus(request)

									return (
										<li className="flex flex-col gap-2 px-5 py-4 sm:flex-row sm:items-center sm:justify-between" key={request.id ?? request.requestId}>
											<span className="font-medium text-slate-800">{request.assetType ?? request.asset_type ?? 'Asset request'}</span>
											<div className="flex items-center gap-3 text-sm">
												<span className={`rounded-full px-2.5 py-1 text-xs font-medium ${statusStyles[status] ?? 'bg-slate-100 text-slate-600'}`}>{status}</span>
												<span className="text-slate-500">{formatDate(getCreatedDate(request))}</span>
											</div>
										</li>
									)
								})}
							</ul>
						)}
					</section>
					</>
				)}
			</div>
			<ProfileCard profile={profile} role={role} />
		</div>
	)
}

export default HomePage