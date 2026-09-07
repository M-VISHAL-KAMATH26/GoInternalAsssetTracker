import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { createRequest } from '../features/requests/requestsApi'

const initialForm = {
  assetType: '',
  category: '',
  justification: '',
}

function NewRequestPage() {
  const navigate = useNavigate()
  const [form, setForm] = useState(initialForm)
  const [errors, setErrors] = useState({})
  const [submitError, setSubmitError] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  const handleChange = (event) => {
    const { name, value } = event.target

    setForm((currentForm) => ({ ...currentForm, [name]: value }))
    setErrors((currentErrors) => ({ ...currentErrors, [name]: '' }))
    setSubmitError('')
  }

  const validate = () => {
    const nextErrors = {}

    if (!form.assetType.trim()) {
      nextErrors.assetType = 'Asset type is required.'
    }

    if (!form.category.trim()) {
      nextErrors.category = 'Category is required.'
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

    try {
      await createRequest({
        assetType: form.assetType.trim(),
        category: form.category.trim(),
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
            <label className="mb-2 block text-sm font-medium text-slate-700" htmlFor="assetType">
              Asset type
            </label>
            <input
              aria-describedby={errors.assetType ? 'assetType-error' : undefined}
              aria-invalid={Boolean(errors.assetType)}
              className="w-full rounded-lg border border-slate-300 px-3 py-2.5 text-sm outline-none transition placeholder:text-slate-400 focus:border-slate-500 focus:ring-2 focus:ring-slate-200"
              id="assetType"
              name="assetType"
              onChange={handleChange}
              placeholder="For example, laptop"
              value={form.assetType}
            />
            {errors.assetType && <p className="mt-1 text-sm text-red-600" id="assetType-error">{errors.assetType}</p>}
          </div>

          <div>
            <label className="mb-2 block text-sm font-medium text-slate-700" htmlFor="category">
              Category
            </label>
            <input
              aria-describedby={errors.category ? 'category-error' : undefined}
              aria-invalid={Boolean(errors.category)}
              className="w-full rounded-lg border border-slate-300 px-3 py-2.5 text-sm outline-none transition placeholder:text-slate-400 focus:border-slate-500 focus:ring-2 focus:ring-slate-200"
              id="category"
              name="category"
              onChange={handleChange}
              placeholder="For example, hardware"
              value={form.category}
            />
            {errors.category && <p className="mt-1 text-sm text-red-600" id="category-error">{errors.category}</p>}
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