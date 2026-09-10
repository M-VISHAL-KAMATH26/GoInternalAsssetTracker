import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { listAssetCatalog } from '../features/assets/assetsApi'
import { createRequest } from '../features/requests/requestsApi'

const initialForm = {
  asset: '',
  justification: '',
}

// The dropdown stores type and category in a single value, since a request
// must match an existing inventory pair to be approvable.
const toAssetValue = (entry) => `${entry.type}||${entry.category}`

function NewRequestPage() {
  const navigate = useNavigate()
  const [form, setForm] = useState(initialForm)
  const [catalog, setCatalog] = useState([])
  const [isLoadingCatalog, setIsLoadingCatalog] = useState(true)
  const [catalogError, setCatalogError] = useState('')
  const [errors, setErrors] = useState({})
  const [submitError, setSubmitError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  useEffect(() => {
    let isMounted = true

    listAssetCatalog()
      .then((response) => {
        if (isMounted) {
          setCatalog(Array.isArray(response.data) ? response.data : [])
        }
      })
      .catch(() => {
        if (isMounted) {
          setCatalogError('Unable to load the available assets. Please try again.')
        }
      })
      .finally(() => {
        if (isMounted) {
          setIsLoadingCatalog(false)
        }
      })

    return () => {
      isMounted = false
    }
  }, [])

  const handleChange = (event) => {
    const { name, value } = event.target

    setForm((currentForm) => ({ ...currentForm, [name]: value }))
    setErrors((currentErrors) => ({ ...currentErrors, [name]: '' }))
    setSubmitError('')
  }

  const validate = () => {
    const nextErrors = {}

    if (!form.asset) {
      nextErrors.asset = 'Select the asset you need.'
    }

    if (!form.justification.trim()) {
      nextErrors.justification = 'Justification is required.'
    }

    return nextErrors
  }

  const handleSubmit = async (event) => {
    event.preventDefault()

    const nextErrors = validate()

    if (Object.keys(nextErrors).length > 0) {
      setErrors(nextErrors)
      return
    }

    setIsSubmitting(true)
    setSubmitError('')

    const [assetType, category] = form.asset.split('||')

    try {
      await createRequest({
        assetType,
        category,
        justification: form.justification.trim(),
      })
      navigate('/requests')
    } catch {
      setSubmitError('Unable to submit your request. Please try again.')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <main className="min-h-screen bg-slate-50 px-4 py-8 text-slate-900 sm:px-6">
      <div className="mx-auto max-w-2xl">
        <div className="mb-6">
          <p className="text-sm font-medium text-slate-500">Employee portal</p>
          <h1 className="mt-1 text-3xl font-semibold tracking-tight">New request</h1>
          <p className="mt-2 text-sm text-slate-600">
            Tell us what asset you need and why it is required.
          </p>
        </div>

        <form
          className="space-y-6 rounded-xl border border-slate-200 bg-white p-6 shadow-sm sm:p-8"
          onSubmit={handleSubmit}
          noValidate
        >
          <div>
            <label className="mb-2 block text-sm font-medium text-slate-700" htmlFor="asset">
              Asset
            </label>
            <select
              aria-describedby={errors.asset ? 'asset-error' : undefined}
              aria-invalid={Boolean(errors.asset)}
              className="w-full rounded-lg border border-slate-300 bg-white px-3 py-2.5 text-sm outline-none transition focus:border-slate-500 focus:ring-2 focus:ring-slate-200 disabled:bg-slate-50"
              disabled={isLoadingCatalog || catalog.length === 0}
              id="asset"
              name="asset"
              onChange={handleChange}
              value={form.asset}
            >
              <option value="">
                {isLoadingCatalog ? 'Loading available assets...' : 'Select an asset'}
              </option>
              {catalog.map((entry) => (
                <option key={toAssetValue(entry)} value={toAssetValue(entry)}>
                  {entry.type} — {entry.category} ({entry.available_count} available)
                </option>
              ))}
            </select>
            {errors.asset && <p className="mt-1 text-sm text-red-600" id="asset-error">{errors.asset}</p>}
            {catalogError && <p className="mt-1 text-sm text-red-600" role="alert">{catalogError}</p>}
            {!isLoadingCatalog && !catalogError && catalog.length === 0 && (
              <p className="mt-1 text-sm text-slate-500">
                No assets are available right now. Ask an administrator to add inventory.
              </p>
            )}
          </div>

          <div>
            <label className="mb-2 block text-sm font-medium text-slate-700" htmlFor="justification">
              Justification
            </label>
            <textarea
              aria-describedby={errors.justification ? 'justification-error' : undefined}
              aria-invalid={Boolean(errors.justification)}
              className="min-h-32 w-full resize-y rounded-lg border border-slate-300 px-3 py-2.5 text-sm outline-none transition placeholder:text-slate-400 focus:border-slate-500 focus:ring-2 focus:ring-slate-200"
              id="justification"
              name="justification"
              onChange={handleChange}
              placeholder="Explain how this asset will support your work"
              value={form.justification}
            />
            {errors.justification && <p className="mt-1 text-sm text-red-600" id="justification-error">{errors.justification}</p>}
          </div>

          {submitError && <p className="text-sm text-red-600" role="alert">{submitError}</p>}

          <div className="flex items-center justify-end gap-3 border-t border-slate-100 pt-5">
            <button
              className="rounded-lg px-4 py-2.5 text-sm font-medium text-slate-600 transition hover:bg-slate-100"
              disabled={isSubmitting}
              onClick={() => navigate('/requests')}
              type="button"
            >
              Cancel
            </button>
            <button
              className="rounded-lg bg-slate-900 px-4 py-2.5 text-sm font-medium text-white transition hover:bg-slate-700 disabled:cursor-not-allowed disabled:opacity-60"
              disabled={isSubmitting}
              type="submit"
            >
              {isSubmitting ? 'Submitting...' : 'Submit request'}
            </button>
          </div>
        </form>
      </div>
    </main>
  )
}

export default NewRequestPage