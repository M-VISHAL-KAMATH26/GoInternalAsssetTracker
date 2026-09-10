import { useCallback, useEffect, useState } from 'react'
import Toast from '../components/Toast'
import { approveRequest, listPendingApprovals, rejectRequest } from '../features/requests/approvalsApi'

const getRequestList = (data) => {
  if (Array.isArray(data)) {
    return data
  }

  return data?.requests ?? data?.items ?? []
}

const getRequestId = (request) => request.id ?? request.requestId

const fetchPendingRequests = async () => {
  const response = await listPendingApprovals()
  return getRequestList(response.data)
}

function ApprovalsPage() {
  const [requests, setRequests] = useState([])
  const [isLoading, setIsLoading] = useState(true)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState('')
  const [feedback, setFeedback] = useState(null)
  const [selectedAction, setSelectedAction] = useState(null)
  const [comment, setComment] = useState('')

  const loadPendingRequests = useCallback(async (showLoading = true) => {
    if (showLoading) {
      setIsLoading(true)
    }

    setError('')

    try {
      setRequests(await fetchPendingRequests())
    } catch {
      setError('Unable to load pending requests. Please try again.')
    } finally {
      setIsLoading(false)
    }
  }, [])

  useEffect(() => {
    let isMounted = true

    const loadInitialRequests = async () => {
      try {
        setRequests(await fetchPendingRequests())
      } catch {
        if (isMounted) {
          setError('Unable to load pending requests. Please try again.')
        }
      } finally {
        if (isMounted) {
          setIsLoading(false)
        }
      }
    }

    loadInitialRequests()

    return () => {
      isMounted = false
    }
  }, [])

  const openAction = (request, action) => {
    setSelectedAction({ action, request })
    setComment('')
    setFeedback(null)
  }

  const closeAction = () => {
    if (!isSubmitting) {
      setSelectedAction(null)
      setComment('')
    }
  }

  const handleAction = async (event) => {
    event.preventDefault()

    if (!selectedAction) {
      return
    }

    const requestId = getRequestId(selectedAction.request)
    const actionRequest = selectedAction.action === 'approve'
      ? approveRequest
      : rejectRequest

    setIsSubmitting(true)
    setFeedback(null)

    try {
      await actionRequest(requestId, comment.trim())
      setSelectedAction(null)
      setComment('')
      setFeedback({
        message: `Request ${selectedAction.action}d successfully.`,
        type: 'success',
      })
      await loadPendingRequests(false)
    } catch {
      setFeedback({
        message: `Unable to ${selectedAction.action} this request. Please try again.`,
        type: 'error',
      })
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <main className="min-h-screen bg-slate-50 px-4 py-8 text-slate-900 sm:px-6">
      <Toast
        message={feedback?.message}
        onClose={() => setFeedback(null)}
        type={feedback?.type}
      />
      <div className="mx-auto max-w-6xl">
        <div className="mb-6">
          <p className="text-sm font-medium text-slate-500">Manager portal</p>
          <h1 className="mt-1 text-3xl font-semibold tracking-tight">Request approvals</h1>
          <p className="mt-2 text-sm text-slate-600">
            Review pending employee requests and record a decision.
          </p>
        </div>

        <section className="overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm">
          {isLoading && (
            <div className="p-8 text-center text-sm text-slate-500">Loading pending requests...</div>
          )}

          {!isLoading && error && (
            <div className="p-8 text-center text-sm text-red-600">{error}</div>
          )}

          {!isLoading && !error && requests.length === 0 && (
            <div className="p-10 text-center">
              <h2 className="text-lg font-medium text-slate-900">No pending requests</h2>
              <p className="mt-2 text-sm text-slate-500">
                New employee requests will appear here for review.
              </p>
            </div>
          )}

          {!isLoading && !error && requests.length > 0 && (
            <div className="divide-y divide-slate-100">
              {requests.map((request) => {
                const requestId = getRequestId(request)
                const isSelected = selectedAction?.request === request

                return (
                  <article className="p-6" key={requestId}>
                    <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
                      <div>
                        <h2 className="text-base font-semibold text-slate-900">
                          {request.assetType ?? request.asset_type ?? 'Asset request'}
                        </h2>
                        <dl className="mt-2 grid gap-x-8 gap-y-1 text-sm text-slate-600 sm:grid-cols-3">
                          <div>
                            <dt className="text-xs uppercase tracking-wide text-slate-400">Category</dt>
                            <dd>{request.category ?? '-'}</dd>
                          </div>
                          <div>
                            <dt className="text-xs uppercase tracking-wide text-slate-400">Employee</dt>
                            <dd>{request.employeeId ?? request.employee_id ?? '-'}</dd>
                          </div>
                          <div>
                            <dt className="text-xs uppercase tracking-wide text-slate-400">Request ID</dt>
                            <dd>{requestId ?? '-'}</dd>
                          </div>
                        </dl>
                        <p className="mt-4 max-w-3xl text-sm leading-6 text-slate-600">
                          {request.justification ?? 'No justification provided.'}
                        </p>
                      </div>
                      <div className="flex shrink-0 gap-2">
                        <button
                          className="rounded-lg bg-emerald-600 px-3 py-2 text-sm font-medium text-white transition hover:bg-emerald-700 disabled:cursor-not-allowed disabled:opacity-60"
                          disabled={isSubmitting}
                          onClick={() => openAction(request, 'approve')}
                          type="button"
                        >
                          Approve
                        </button>
                        <button
                          className="rounded-lg border border-red-200 px-3 py-2 text-sm font-medium text-red-700 transition hover:bg-red-50 disabled:cursor-not-allowed disabled:opacity-60"
                          disabled={isSubmitting}
                          onClick={() => openAction(request, 'reject')}
                          type="button"
                        >
                          Reject
                        </button>
                      </div>
                    </div>

                    {isSelected && (
                      <form className="mt-5 rounded-lg border border-slate-200 bg-slate-50 p-4" onSubmit={handleAction}>
                        <label className="block text-sm font-medium text-slate-700" htmlFor={`comment-${requestId}`}>
                          Comment <span className="font-normal text-slate-500">(optional)</span>
                          <textarea
                            className="mt-2 min-h-24 w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm outline-none transition focus:border-slate-500 focus:ring-2 focus:ring-slate-200"
                            disabled={isSubmitting}
                            id={`comment-${requestId}`}
                            onChange={(event) => setComment(event.target.value)}
                            placeholder="Add context for the employee"
                            value={comment}
                          />
                        </label>
                        <div className="mt-3 flex justify-end gap-2">
                          <button
                            className="rounded-lg px-3 py-2 text-sm font-medium text-slate-600 hover:bg-slate-200 disabled:opacity-60"
                            disabled={isSubmitting}
                            onClick={closeAction}
                            type="button"
                          >
                            Cancel
                          </button>
                          <button
                            className="rounded-lg bg-slate-900 px-3 py-2 text-sm font-medium text-white hover:bg-slate-700 disabled:cursor-not-allowed disabled:opacity-60"
                            disabled={isSubmitting}
                            type="submit"
                          >
                            {isSubmitting ? 'Saving...' : `Confirm ${selectedAction.action}`}
                          </button>
                        </div>
                      </form>
                    )}
                  </article>
                )
              })}
            </div>
          )}
        </section>
      </div>
    </main>
  )
}

export default ApprovalsPage