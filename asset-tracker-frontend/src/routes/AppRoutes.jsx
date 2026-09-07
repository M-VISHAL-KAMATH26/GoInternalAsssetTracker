import { Route, Routes } from 'react-router-dom'
import ApprovalsPage from '../pages/ApprovalsPage'
import AssetsPage from '../pages/AssetsPage'
import DevLoginPage from '../pages/DevLoginPage'
import HomePage from '../pages/HomePage'
import NewRequestPage from '../pages/NewRequestPage'
import RequestsListPage from '../pages/RequestsListPage'
import UnauthorizedPage from '../pages/UnauthorizedPage'
import Layout from '../components/Layout'
import ProtectedRoute from './ProtectedRoute'

function AppRoutes() {
	return (
		<Routes>
			<Route element={<DevLoginPage />} path="/login" />
			<Route element={<UnauthorizedPage />} path="/not-authorized" />

			<Route element={<ProtectedRoute />}>
				<Route element={<Layout />}>
					<Route element={<HomePage />} path="/" />
					<Route element={<RequestsListPage />} path="/requests" />
					<Route element={<NewRequestPage />} path="/requests/new" />
				</Route>
			</Route>

			<Route element={<ProtectedRoute allowedRoles={['manager', 'admin']} />}>
				<Route element={<Layout />}>
					<Route element={<ApprovalsPage />} path="/approvals" />
				</Route>
			</Route>

			<Route element={<ProtectedRoute allowedRoles={['admin']} />}>
				<Route element={<Layout />}>
					<Route element={<AssetsPage />} path="/assets" />
				</Route>
			</Route>
		</Routes>
	)
}

export default AppRoutes
