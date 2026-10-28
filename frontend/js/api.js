// Konfigurasi URL dasar API
// Selalu memakai path relatif /api:
// - Mode development (Vite): /api di-proxy ke backend Go (lihat vite.config.js)
// - Mode produksi: disajikan langsung oleh backend Go di origin yang sama
export const API_BASE = "/api";

/**
 * Wrapper serbaguna untuk panggilan fetch() ke backend API
 */
export async function apiRequest(endpoint, options = {}) {
  const token = localStorage.getItem("token");

  const headers = {
    "Content-Type": "application/json",
    ...(options.headers || {})
  };

  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const config = {
    ...options,
    headers
  };

  try {
    const res = await fetch(`${API_BASE}${endpoint}`, config);

    // Jika belum login / token kedaluwarsa (401)
    if (res.status === 401) {
      localStorage.removeItem("token");
      localStorage.removeItem("admin_user");
      
      const currentPath = window.location.pathname;
      if (!currentPath.includes("login.html") && !currentPath.includes("register.html")) {
        window.location.href = "login.html";
      }
      const data = await res.json().catch(() => ({}));
      throw new Error(data.error || "Sesi telah berakhir, silakan login kembali.");
    }

    const data = await res.json().catch(() => ({}));

    if (!res.ok) {
      throw new Error(data.error || `Terjadi kesalahan (Status ${res.status})`);
    }

    return data;
  } catch (err) {
    throw err;
  }
}

/**
 * Menampilkan pesan toast melayang (notifikasi halus tanpa popup alert)
 */
export function showToast(message, type = "success") {
  let container = document.getElementById("toastContainer");
  if (!container) {
    container = document.createElement("div");
    container.id = "toastContainer";
    container.className = "toast-container";
    document.body.appendChild(container);
  }

  const toast = document.createElement("div");
  toast.className = `toast toast-${type}`;
  toast.textContent = message;

  container.appendChild(toast);

  setTimeout(() => {
    toast.style.opacity = "0";
    toast.style.transition = "opacity 0.3s ease";
    setTimeout(() => {
      if (toast.parentNode) {
        toast.parentNode.removeChild(toast);
      }
    }, 300);
  }, 3500);
}
