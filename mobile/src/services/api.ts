// mobile/src/services/api.ts
const LOCAL_IP = "192.168.1.7" ; 
const BASE_URL = `http://${LOCAL_IP}:8080/api/v1`;

export const apiCall = async (endpoint: string, method: string = "GET", body?: any, token?: string) => {
  const headers: HeadersInit = {
    "Content-Type": "application/json",
  };

  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
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