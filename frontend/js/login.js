// Entry point halaman Login (di-bundle Vite)
import { apiRequest, showToast } from "./api.js";
import { checkAuth } from "./auth.js";

// Jika sudah punya token valid, langsung alihkan ke dashboard
checkAuth(false);

// Pasang event listener form login
const loginForm = document.getElementById("loginForm");
if (loginForm) {
  loginForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    const submitBtn = loginForm.querySelector("button[type='submit']");
    const originalText = submitBtn.textContent;

    const username = document.getElementById("username").value.trim();
    const password = document.getElementById("password").value;

    if (!username || !password) {
      showToast("Username dan password wajib diisi", "error");
      return;
    }

    try {
      submitBtn.disabled = true;
      submitBtn.textContent = "Memproses...";

      const res = await apiRequest("/auth/login", {
        method: "POST",
        body: JSON.stringify({ username, password }),
      });

      localStorage.setItem("token", res.token);
      showToast("Login berhasil! Mengalihkan...", "success");

      setTimeout(() => {
        window.location.href = "index.html";
      }, 700);
    } catch (err) {
      showToast(err.message, "error");
    } finally {
      submitBtn.disabled = false;
      submitBtn.textContent = originalText;
    }
  });
}
