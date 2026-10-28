import { apiRequest } from "./api.js";

/**
 * Pemeriksaan status autentikasi pengguna
 */
export async function checkAuth(requiresAuth = true) {
  const token = localStorage.getItem("token");

  if (requiresAuth) {
    if (!token) {
      window.location.href = "login.html";
      return null;
    }

    try {
      const res = await apiRequest("/auth/me");
      const userBadge = document.getElementById("currentUserBadge");
      if (userBadge && res.user) {
        userBadge.textContent = `👤 ${res.user.username}`;
      }
      return res.user;
    } catch (err) {
      localStorage.removeItem("token");
      window.location.href = "login.html";
      return null;
    }
  } else {
    // Pada halaman login/register: jika sudah punya token valid, arahkan ke halaman utama
    if (token) {
      try {
        await apiRequest("/auth/me");
        window.location.href = "index.html";
      } catch (err) {
        localStorage.removeItem("token");
      }
    }
  }
}

/**
 * Logout admin
 */
export function logout() {
  localStorage.removeItem("token");
  localStorage.removeItem("admin_user");
  window.location.href = "login.html";
}
