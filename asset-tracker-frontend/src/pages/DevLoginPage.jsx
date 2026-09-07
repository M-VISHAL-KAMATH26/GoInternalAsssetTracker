import { useState } from 'react'
import { useDispatch } from 'react-redux'
import { setCredentials } from '../features/auth/authSlice'

function DevLoginPage() {
  const [token, setToken] = useState('')
  const dispatch = useDispatch()

  const handleSubmit = (event) => {
    event.preventDefault()

    const trimmedToken = token.trim()

    if (trimmedToken) {
      dispatch(setCredentials(trimmedToken))
    }
  }

  return (
    <main className="mx-auto max-w-xl p-6">
      <h1 className="mb-4 text-2xl font-semibold">Developer Login</h1>
      <form className="space-y-4" onSubmit={handleSubmit}>
        <label className="block" htmlFor="dev-jwt">
          <span className="mb-2 block text-sm font-medium">JWT token</span>
          <textarea
            className="min-h-40 w-full rounded border p-3 font-mono text-sm"
            id="dev-jwt"
            onChange={(event) => setToken(event.target.value)}
            placeholder="Paste a generated JWT"
            value={token}
          />
        </label>
        <button
          className="rounded bg-slate-900 px-4 py-2 text-white"
          type="submit"
        >
          Set credentials
        </button>
      </form>
    </main>
  )
}

export default DevLoginPage