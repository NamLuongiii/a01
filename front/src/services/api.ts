/**
 * src/services/api.ts
 * Đơn giản: wrapper nhỏ quanh fetch với helpers get/post/put/delete
 */

const BASE_URL = import.meta.env.VITE_API_BASE_URL || "";

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    headers: {
      "Content-Type": "application/json",
      ...(options.headers || {}),
    },
    ...options,
  });

  if (!res.ok) {
    const text = await res.text().catch(() => "");
    const message = text || res.statusText || "Request failed";
    throw new Error(`${res.status} ${message}`);
  }

  // Nếu không có body trả về (204 No Content), cast về any
  if (res.status === 204) {
    return undefined as unknown as T;
  }

  return (await res.json()) as T;
}

export default {
  get: <T = any>(path: string) => request<T>(path, { method: "GET" }),
  post: <T = any, B = unknown>(path: string, body?: B) =>
    request<T>(path, {
      method: "POST",
      body: body ? JSON.stringify(body) : undefined,
    }),
  put: <T = any, B = unknown>(path: string, body?: B) =>
    request<T>(path, {
      method: "PUT",
      body: body ? JSON.stringify(body) : undefined,
    }),
  delete: <T = any>(path: string) => request<T>(path, { method: "DELETE" }),
};
