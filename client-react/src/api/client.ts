import axios from "axios";

const baseURL = import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8000";

export const apiClient = axios.create({
  baseURL,
  withCredentials: true,
  headers: {
    "Content-Type": "application/x-www-form-urlencoded"
  }
});

const toFormData = (data?: Record<string, unknown>) => {
  const form = new URLSearchParams();

  if (!data) {
    return form;
  }

  Object.entries(data).forEach(([key, value]) => {
    if (value === undefined || value === null) {
      return;
    }

    if (Array.isArray(value)) {
      value.forEach((item) => {
        form.append(key, String(item));
      });
      return;
    }

    form.append(key, String(value));
  });

  return form;
};

export const postForm = <T>(url: string, data?: Record<string, unknown>) => {
  return apiClient.post<T>(url, toFormData(data));
};
