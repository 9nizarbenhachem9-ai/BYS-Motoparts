// Entry point halaman Register (di-bundle Vite)
import { apiRequest, showToast } from "./api.js";
import { checkAuth } from "./auth.js";

// Jika sudah punya token valid, langsung alihkan ke dashboard
checkAuth(false);

// Pasang event listener form registrasi
const registerForm = document.getElementById("registerForm");
if (registerForm) {
  registerForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    const submitBtn = registerForm.querySelector("button[type='submit']");
    const originalText = submitBtn.textContent;

    const username = document.getElementById("username").value.trim();
    const password = document.getElementById("password").value;
    const inviteCode = document.getElementById("inviteCode")
      ? document.getElementById("inviteCode").value.trim()
      : "";

    if (!username || !password) {
      showToast("Username dan password wajib diisi", "error");
      return;
    }

    if (password.length < 8) {
      showToast("Password minimal 8 karakter", "error");
      return;
    }

    try {
      submitBtn.disabled = true;
      submitBtn.textContent = "Mendaftarkan...";

      await apiRequest("/auth/register", {
        method: "POST",
        body: JSON.stringify({
          username,
          password,
          invite_code: inviteCode,
        }),
      });

      showToast("Registrasi berhasil! Silakan login.", "success");
      setTimeout(() => {
        window.location.href = "login.html";
      }, 1000);
    } catch (err) {
      showToast(err.message, "error");
    } finally {
      submitBtn.disabled = false;
      submitBtn.textContent = originalText;
    }
  });
}
