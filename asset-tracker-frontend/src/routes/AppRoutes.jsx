import { Route, Routes } from 'react-router-dom'
import AdminCreateUserPage from '../pages/AdminCreateUserPage'
import AdminDashboardPage from '../pages/AdminDashboardPage'
import AdminLoginPage from '../pages/AdminLoginPage'
import ApprovalsPage from '../pages/ApprovalsPage'
import AssetsPage from '../pages/AssetsPage'
import HomePage from '../pages/HomePage'
import LoginPage from '../pages/LoginPage'
import NewRequestPage from '../pages/NewRequestPage'
import RequestsListPage from '../pages/RequestsListPage'
import UnauthorizedPage from '../pages/UnauthorizedPage'
import Layout from '../components/Layout'
import ProtectedRoute from './ProtectedRoute'

function AppRoutes() {
	return (
		<Routes>
			<Route element={<LoginPage />} path="/login" />
			<Route element={<AdminLoginPage />} path="/admin/login" />
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

			<Route element={<ProtectedRoute allowedRoles={['admin']} redirectTo="/admin/login" />}>
				<Route element={<Layout />}>
					<Route element={<AdminDashboardPage />} path="/admin" />
					<Route element={<AdminCreateUserPage />} path="/admin/users/new" />
				</Route>
			</Route>
		</Routes>
	)
}

export default AppRoutes
