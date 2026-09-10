import { Boxes, ClipboardList, ShieldCheck, UserPlus } from 'lucide-react'
import { Link } from 'react-router-dom'
import { useSelector } from 'react-redux'

const steps = [
	{
		icon: UserPlus,
		title: '1. An admin creates your account',
		description:
			'Administrators add employees and managers with a role, office location, and profile picture. There is no public sign-up — you receive credentials from your admin.',
	},
	{
		icon: Boxes,
		title: '2. Inventory is stocked',
		description:
			'Admins register the laptops, monitors, and other equipment the company owns. Only assets that exist in inventory can be requested.',
	},
	{
		icon: ClipboardList,
		title: '3. Employees raise a request',
		description:
			'Pick what you need from the available asset list and explain why. Your request starts as pending and is tracked on your dashboard.',
	},
	{
		icon: ShieldCheck,
		title: '4. Your manager decides',
		description:
			'The manager you report to approves or rejects the request. On approval a matching asset is reserved for you automatically.',
	},
]

const roles = [
	{ role: 'Employee', can: 'Raise asset requests and track their status.' },
	{ role: 'Manager', can: 'Do everything an employee can, plus approve or reject requests from their reports.' },
	{ role: 'Admin', can: 'Create users, manage inventory, and monitor every request in the system.' },
]

function LandingPage() {
	const { token, role } = useSelector((state) => state.auth)
	const dashboardPath = role === 'admin' ? '/admin' : '/requests'

	return (
		<div className="min-h-screen bg-slate-50 text-slate-900">
			<header className="border-b border-slate-200 bg-white">
				<div className="mx-auto flex max-w-6xl items-center justify-between px-4 py-4 sm:px-6">
					<span className="text-base font-semibold tracking-tight">Asset Tracker</span>
					<nav aria-label="Account" className="flex items-center gap-2">
						{token ? (
							<Link
								className="rounded-lg bg-slate-900 px-4 py-2 text-sm font-medium text-white transition hover:bg-slate-700"
								to={dashboardPath}
							>
								Go to dashboard
							</Link>
						) : (
							<>
								<Link
									className="rounded-lg px-3 py-2 text-sm font-medium text-slate-600 transition hover:bg-slate-100 hover:text-slate-900"
									to="/admin/login"
								>
									Admin login
								</Link>
								<Link
									className="rounded-lg bg-slate-900 px-4 py-2 text-sm font-medium text-white transition hover:bg-slate-700"
									to="/login"
								>
									Sign in
								</Link>
							</>
						)}
					</nav>
				</div>
			</header>

			<main>
				<section className="mx-auto max-w-6xl px-4 py-16 sm:px-6 sm:py-20">
					<p className="text-sm font-medium text-indigo-600">Internal equipment portal</p>
					<h1 className="mt-3 max-w-3xl text-4xl font-bold tracking-tight sm:text-5xl">
						Request company equipment without chasing anyone over email
					</h1>
					<p className="mt-5 max-w-2xl text-lg leading-8 text-slate-600">
						Asset Tracker keeps equipment requests, manager approvals, and company inventory in one
						place, so employees know where their request stands and admins know where every asset went.
					</p>
					<div className="mt-8 flex flex-wrap gap-3">
						<Link
							className="rounded-lg bg-indigo-600 px-5 py-3 text-sm font-medium text-white transition hover:bg-indigo-700"
							to={token ? dashboardPath : '/login'}
						>
							{token ? 'Go to dashboard' : 'Sign in to your account'}
						</Link>
						<a
							className="rounded-lg border border-slate-300 bg-white px-5 py-3 text-sm font-medium text-slate-700 transition hover:bg-slate-100"
							href="#how-it-works"
						>
							See how it works
						</a>
					</div>
				</section>

				<section className="border-y border-slate-200 bg-white py-16" id="how-it-works">
					<div className="mx-auto max-w-6xl px-4 sm:px-6">
						<h2 className="text-2xl font-bold tracking-tight">How it works</h2>
						<p className="mt-2 text-sm text-slate-600">Four steps from a new account to equipment in hand.</p>
						<ol className="mt-8 grid gap-6 sm:grid-cols-2">
							{steps.map(({ icon: Icon, title, description }) => (
								<li className="rounded-xl border border-slate-200 p-6" key={title}>
									<div className="flex h-10 w-10 items-center justify-center rounded-lg bg-indigo-50 text-indigo-600">
										<Icon aria-hidden="true" className="h-5 w-5" />
									</div>
									<h3 className="mt-4 font-semibold">{title}</h3>
									<p className="mt-2 text-sm leading-6 text-slate-600">{description}</p>
								</li>
							))}
						</ol>
					</div>
				</section>

				<section className="mx-auto max-w-6xl px-4 py-16 sm:px-6">
					<h2 className="text-2xl font-bold tracking-tight">What each role can do</h2>
					<dl className="mt-8 space-y-4">
						{roles.map(({ role: roleName, can }) => (
							<div className="rounded-xl border border-slate-200 bg-white p-5 sm:flex sm:gap-6" key={roleName}>
								<dt className="w-32 shrink-0 font-semibold">{roleName}</dt>
								<dd className="mt-1 text-sm leading-6 text-slate-600 sm:mt-0">{can}</dd>
							</div>
						))}
					</dl>
					<div className="mt-10 rounded-xl border border-slate-200 bg-white p-6">
						<h3 className="font-semibold">Need an account?</h3>
						<p className="mt-2 text-sm leading-6 text-slate-600">
							Accounts are created by an administrator, so contact your admin to be added. Once you have
							credentials, sign in and your dashboard will show your requests and profile.
						</p>
						<Link className="mt-4 inline-block text-sm font-medium text-indigo-600 hover:text-indigo-700" to="/login">
							Sign in →
						</Link>
					</div>
				</section>
			</main>

			<footer className="border-t border-slate-200 bg-white py-6">
				<p className="mx-auto max-w-6xl px-4 text-sm text-slate-500 sm:px-6">
					Asset Tracker — internal equipment request and approval system.
				</p>
			</footer>
		</div>
	)
}

export default LandingPage
