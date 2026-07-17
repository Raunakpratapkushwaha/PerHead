const BASE_URL = process.env.EXPO_PUBLIC_API_URL;

export const apiCall = async (endpoint: string, method: string = "GET", body?: any, token?: string) => {
  const headers: HeadersInit = {
    "Content-Type": "application/json",
  };

  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  // Ensure BASE_URL exists just in case the .env fails to load
  if (!BASE_URL) {
    throw new Error("API URL is not defined in environment variables");
  }

  const response = await fetch(`${BASE_URL}${endpoint}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  });

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.error || "Something went wrong");
  }

  return data;
};