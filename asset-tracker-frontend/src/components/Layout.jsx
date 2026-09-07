import { Outlet } from 'react-router-dom'
import NavBar from './NavBar'

function Layout() {
	return (
		<div className="min-h-screen bg-slate-50 text-slate-900">
			<NavBar />
			<main className="mx-auto w-full max-w-7xl px-4 py-8 sm:px-6">
				<Outlet />
			</main>
		</div>
	)
}

export default Layout