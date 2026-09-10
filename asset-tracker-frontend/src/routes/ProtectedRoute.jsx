import { Navigate, Outlet, useLocation } from 'react-router-dom'
import { useSelector } from 'react-redux'

function ProtectedRoute({ allowedRoles, children, redirectTo = '/login' }) {
	const location = useLocation()
	const { token, role } = useSelector((state) => state.auth)

	if (!token) {
		return <Navigate replace state={{ from: location }} to={redirectTo} />
	}

	if (allowedRoles && !allowedRoles.includes(role)) {
		return <Navigate replace to="/not-authorized" />
	}

	return children ?? <Outlet />
}

export default ProtectedRoute
