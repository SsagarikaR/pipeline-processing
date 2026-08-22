import axios from 'axios';
import toast from 'react-hot-toast';

const axiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081/api/v1',
  headers: {
    'Content-Type': 'application/json',
    'X-API-Key': import.meta.env.VITE_API_KEY || 'secret-pipeline-key',
  },
});

// Add a response interceptor
axiosInstance.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    // Show a toast message for errors. The backend returns errors via
    // http.Error, i.e. a plain-text body, not a JSON { message } shape.
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
