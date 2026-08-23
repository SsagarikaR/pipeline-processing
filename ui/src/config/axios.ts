import axios from 'axios';
import toast from 'react-hot-toast';

/**
 * Shared axios client for every API call in the app. Points at the
 * backend's base URL and sends the API key it requires on every
 * request; both are overridable via env vars for non-local environments.
 */
const axiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081/api/v1',
  headers: {
    'Content-Type': 'application/json',
    'X-API-Key': import.meta.env.VITE_API_KEY || 'secret-pipeline-key',
  },
});

/**
 * On any failed request, surface the backend's error message as a toast
 * so the user always sees why something went wrong, without every
 * caller having to handle it individually.
 */
axiosInstance.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    // The backend returns errors via http.Error, i.e. a plain-text
    // body, not a JSON { message } shape.
    const data = error.response?.data;
    const errorMessage =
      (typeof data === 'string' && data.trim()) ||
      data?.message ||
      error.message ||
      'An unexpected error occurred';
    toast.error(errorMessage);
    return Promise.reject(error);
  }
);

export default axiosInstance;
