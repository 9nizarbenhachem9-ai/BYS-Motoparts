/**
 * Pemeriksaan status autentikasi pengguna
 */
async function checkAuth(requiresAuth = true) {
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
function logout() {
  localStorage.removeItem("token");
  localStorage.removeItem("admin_user");
  window.location.href = "login.html";
}

// Inisialisasi event listener jika berada di halaman Login
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
        body: JSON.stringify({ username, password })
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

// Inisialisasi event listener jika berada di halaman Register
const registerForm = document.getElementById("registerForm");
if (registerForm) {
  registerForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    const submitBtn = registerForm.querySelector("button[type='submit']");
    const originalText = submitBtn.textContent;

    const username = document.getElementById("username").value.trim();
    const password = document.getElementById("password").value;
    const inviteCode = document.getElementById("inviteCode") ? document.getElementById("inviteCode").value.trim() : "";

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
          invite_code: inviteCode
        })
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
