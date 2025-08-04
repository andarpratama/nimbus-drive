export const API_BASE_URL = (import.meta as any).env.VITE_API_URL || 'http://localhost:8080'
export const getAuthHeaders = (): Record<string, string> => {
  const token = localStorage.getItem('token')
  return {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
}
export const apiRequest = async <T>(
  endpoint: string, 
  options: RequestInit = {}
): Promise<T> => {
  const url = `${API_BASE_URL}${endpoint}`
  const response = await fetch(url, {
    ...options,
    headers: {
      ...getAuthHeaders(),
      ...options.headers
    }
  })
  if (!response.ok) {
    throw new Error(`API request failed: ${response.statusText}`)
  }
  const data = await response.json()
  return data
} 