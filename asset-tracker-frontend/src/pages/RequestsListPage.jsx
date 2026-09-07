import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { listMyRequests } from '../features/requests/requestsApi'

const getRequestList = (data) => {
  if (Array.isArray(data)) {
    return data
  }

  return data?.requests ?? data?.items ?? []
}

const formatDate = (value) => {
  if (!value) {
    return '-'
  }

  const date = new Date(value)

  return Number.isNaN(date.getTime())
    ? '-'
    : new Intl.DateTimeFormat(undefined, {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
      }).format(date)
}

function RequestsListPage() {
  const [requests, setRequests] = useState([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    const loadRequests = async () => {
      try {
        const response = await listMyRequests()
        setRequests(getRequestList(response.data))
      } catch {
        setError('Unable to load your requests. Please try again.')
      } finally {
        setIsLoading(false)
      }
    }

    loadRequests()
  }, [])

  return (
    <main className="min-h-screen bg-slate-50 px-4 py-8 text-slate-900 sm:px-6">
      <div className="mx-auto max-w-6xl">
        <div className="mb-6 flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <p className="text-sm font-medium text-slate-500">Employee portal</p>
            <h1 className="mt-1 text-3xl font-semibold tracking-tight">My requests</h1>
            <p className="mt-2 text-sm text-slate-600">
              Track the status of your submitted asset requests.
            </p>
          </div>
          <Link
            className="inline-flex items-center justify-center rounded-lg bg-slate-900 px-4 py-2.5 text-sm font-medium text-white transition hover:bg-slate-700"
            to="/requests/new"
          >
            New request
          </Link>
        </div>

        <section className="overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm">
          {isLoading && (
            <div className="p-8 text-center text-sm text-slate-500">Loading requests...</div>
          )}

          {!isLoading && error && (
            <div className="p-8 text-center text-sm text-red-600">{error}</div>
          )}

          {!isLoading && !error && requests.length === 0 && (
            <div className="p-10 text-center">
              <h2 className="text-lg font-medium text-slate-900">No requests yet</h2>
              <p className="mt-2 text-sm text-slate-500">
                Create a request when you need a new asset.
              </p>
            </div>
          )}

          {!isLoading && !error && requests.length > 0 && (
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-slate-200 text-left text-sm">
                <thead className="bg-slate-50 text-xs uppercase tracking-wide text-slate-500">
                  <tr>
                    <th className="px-6 py-3 font-medium" scope="col">Asset type</th>
                    <th className="px-6 py-3 font-medium" scope="col">Category</th>
                    <th className="px-6 py-3 font-medium" scope="col">Status</th>
                    <th className="px-6 py-3 font-medium" scope="col">Created</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {requests.map((request) => (
                    <tr className="text-slate-700" key={request.id ?? request.requestId}>
                      <td className="whitespace-nowrap px-6 py-4 font-medium text-slate-900">
                        {request.assetType ?? request.asset_type ?? '-'}
                      </td>
                      <td className="whitespace-nowrap px-6 py-4">
                        {request.category ?? '-'}
                      </td>
                      <td className="whitespace-nowrap px-6 py-4">
                        {request.status ?? '-'}
                      </td>
                      <td className="whitespace-nowrap px-6 py-4">
                        {formatDate(request.createdAt ?? request.created_at)}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </section>
      </div>
    </main>
  )
}

export default RequestsListPage