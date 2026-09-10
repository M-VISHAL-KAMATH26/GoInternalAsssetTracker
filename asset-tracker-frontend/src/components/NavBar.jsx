import { NavLink, useNavigate } from 'react-router-dom'
import { useDispatch, useSelector } from 'react-redux'
import { logout } from '../features/auth/authSlice'

const linkClassName = ({ isActive }) =>
	`rounded-md px-3 py-2 text-sm font-medium transition ${isActive
		? 'bg-slate-100 text-slate-900'
		: 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'}`

function NavBar() {
	const dispatch = useDispatch()
	const navigate = useNavigate()
	const role = useSelector((state) => state.auth.role)?.toLowerCase()

	const handleLogout = () => {
		dispatch(logout())
		navigate('/')
	}

	return (
		<header className="border-b border-slate-200 bg-white">
			<div className="mx-auto flex max-w-7xl flex-wrap items-center gap-4 px-4 py-3 sm:px-6">
				<NavLink className="mr-2 text-base font-semibold tracking-tight text-slate-900" to={role === 'admin' ? '/admin' : '/requests'}>
					Asset Tracker
				</NavLink>
				<nav aria-label="Main navigation" className="flex flex-1 flex-wrap items-center gap-1">
					{role === 'admin' && <NavLink className={linkClassName} to="/admin">Admin panel</NavLink>}
					<NavLink className={linkClassName} end to="/requests">Dashboard</NavLink>
					<NavLink className={linkClassName} to="/requests/history">My requests</NavLink>
					<NavLink className={linkClassName} to="/requests/new">New request</NavLink>
					{(role === 'manager' || role === 'admin') && (
						<NavLink className={linkClassName} to="/approvals">Approvals</NavLink>
					)}
					{role === 'admin' && (
						<>
							<NavLink className={linkClassName} to="/admin/users/new">Add user</NavLink>
							<NavLink className={linkClassName} to="/assets">Assets</NavLink>
						</>
					)}
				</nav>
				<div className="flex items-center gap-3">
					<span className="hidden text-xs font-medium uppercase tracking-wide text-slate-400 sm:inline">
						{role ?? 'user'}
					</span>
					<button
						className="rounded-md border border-slate-200 px-3 py-2 text-sm font-medium text-slate-600 transition hover:bg-slate-50 hover:text-slate-900"
						onClick={handleLogout}
						type="button"
					>
						Logout
					</button>
				</div>
			</div>
		</header>
	)
}

export default NavBar