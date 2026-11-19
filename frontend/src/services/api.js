const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080/api";

export function getToken() {
  return localStorage.getItem("driverent_token");
}

export async function apiRequest(path, options = {}) {
  const headers = new Headers(options.headers || {});

  // FormData tidak boleh diberi Content-Type secara manual.
  // Browser akan mengatur multipart boundary secara otomatis.
  if (!(options.body instanceof FormData) && options.body) {
    headers.set("Content-Type", "application/json");
  }

  const token = getToken();

  if (token && !headers.has("Authorization")) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  let response;

  try {
    response = await fetch(`${API_URL}${path}`, {
      ...options,
      headers,
    });
  } catch {
    throw new Error(
      "Tidak dapat terhubung ke server. Pastikan backend Go sedang berjalan.",
    );
  }

  const text = await response.text();

  let data = null;

  if (text) {
    try {
      data = JSON.parse(text);
    } catch {
      data = { message: text };
    }
  }

  if (!response.ok) {
    throw new Error(
      data?.error || data?.message || `Request gagal (${response.status})`,
    );
  }

  return data;
}

// Helper untuk request kendaraan.
// Bisa menerima FormData maupun object JavaScript biasa.
function vehiclePayload(payload) {
  if (payload instanceof FormData) {
    return payload;
  }

  return JSON.stringify(payload);
}

// =========================
// AUTH API
// =========================

export const authApi = {
  login: (payload) =>
    apiRequest("/auth/login", {
      method: "POST",
      body: JSON.stringify(payload),
    }),

  register: (payload) =>
    apiRequest("/auth/register", {
      method: "POST",
      body: JSON.stringify(payload),
    }),

  me: () => apiRequest("/auth/me"),
};

// =========================
// VEHICLE API
// =========================

export const vehicleApi = {
  all: () => apiRequest("/vehicles"),

  byId: (id) => apiRequest(`/vehicles/${id}`),

  create: (payload) =>
    apiRequest("/vehicles", {
      method: "POST",
      body: vehiclePayload(payload),
    }),

  update: (id, payload) =>
    apiRequest(`/vehicles/${id}`, {
      method: "PUT",
      body: vehiclePayload(payload),
    }),

  updateStatus: (id, status) =>
    apiRequest(`/vehicles/${id}/status`, {
      method: "PATCH",
      body: JSON.stringify({ status }),
    }),

  remove: (id) =>
    apiRequest(`/vehicles/${id}`, {
      method: "DELETE",
    }),
};

// =========================
// RESERVATION API
// =========================

export const reservationApi = {
  all: () => apiRequest("/reservations"),

  byId: (id) => apiRequest(`/reservations/${id}`),

  create: (payload) =>
    apiRequest("/reservations", {
      method: "POST",
      body: JSON.stringify(payload),
    }),

  mine: (userId) => apiRequest(`/users/${userId}/reservations`),

  updateStatus: (id, status) =>
    apiRequest(`/reservations/${id}/status`, {
      method: "PATCH",
      body: JSON.stringify({ status }),
    }),

  cancel: (id) =>
    apiRequest(`/reservations/${id}/cancel`, {
      method: "PATCH",
    }),
};

// =========================
// PAYMENT API
// =========================

export const paymentApi = {
  create: (payload) =>
    apiRequest("/payments", {
      method: "POST",
      body: JSON.stringify(payload),
    }),

  all: () => apiRequest("/payments"),

  verify: (id) =>
    apiRequest(`/payments/${id}/verify`, {
      method: "PATCH",
    }),

  reject: (id) =>
    apiRequest(`/payments/${id}/reject`, {
      method: "PATCH",
    }),
};

// =========================
// RENTAL API
// =========================

export const rentalApi = {
  all: () => apiRequest("/rentals"),

  byId: (id) => apiRequest(`/rentals/${id}`),

  create: (payload) =>
    apiRequest("/rentals", {
      method: "POST",
      body: JSON.stringify(payload),
    }),
};

// =========================
// RETURN API
// =========================

export const returnApi = {
  all: () => apiRequest("/returns"),

  create: (payload) =>
    apiRequest("/returns", {
      method: "POST",
      body: JSON.stringify(payload),
    }),
};
