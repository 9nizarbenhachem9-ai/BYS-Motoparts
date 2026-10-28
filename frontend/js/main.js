// Impor modul (di-bundle oleh Vite)
import { apiRequest, showToast } from "./api.js";
import { checkAuth, logout } from "./auth.js";

// State aplikasi
let state = {
  items: [],
  page: 1,
  limit: 10,
  totalPages: 1,
  totalItems: 0,
  search: "",
  kategori: "",
  selectedItem: null // item yang sedang dipilih untuk aksi / klik kanan
};

// Format mata uang Rupiah
function formatRupiah(number) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0
  }).format(number);
}

// Format tanggal & waktu
function formatDateTime(isoString) {
  if (!isoString) return "-";
  const d = new Date(isoString);
  return d.toLocaleDateString("id-ID", {
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit"
  });
}

/**
 * Memuat data sparepart dari backend API
 */
async function loadSpareparts(page = 1) {
  state.page = page;
  const tbody = document.getElementById("sparepartTableBody");
  tbody.innerHTML = `<tr><td colspan="9" style="text-align:center; padding:30px;">Memuat data sparepart...</td></tr>`;

  try {
    const params = new URLSearchParams({
      page: state.page,
      limit: state.limit,
      search: state.search,
      kategori: state.kategori
    });

    const res = await apiRequest(`/spareparts?${params.toString()}`);
    state.items = res.items || [];
    state.totalItems = res.total_items || 0;
    state.totalPages = res.total_pages || 1;

    renderTable();
    renderPagination();
    updateMetrics();
  } catch (err) {
    tbody.innerHTML = `<tr><td colspan="9" style="text-align:center; color:var(--danger); padding:30px;">${escapeHtml(err.message)}</td></tr>`;
    showToast(err.message, "error");
  }
}

/**
 * Render baris tabel sparepart secara aman dari celah XSS
 */
function renderTable() {
  const tbody = document.getElementById("sparepartTableBody");
  tbody.innerHTML = "";

  if (state.items.length === 0) {
    tbody.innerHTML = `
      <tr>
        <td colspan="9" class="empty-state">
          <p>Belum ada data sparepart yang sesuai.</p>
        </td>
      </tr>
    `;
    return;
  }

  state.items.forEach((item) => {
    const tr = document.createElement("tr");
    tr.dataset.id = item.id;

    // Cek apakah stok menipis (stok <= stok_minimum)
    const isLowStock = item.stok <= item.stok_minimum;

    // Kolom Kode
    const tdKode = document.createElement("td");
    const codeSpan = document.createElement("strong");
    codeSpan.textContent = item.kode;
    tdKode.appendChild(codeSpan);

    // Kolom Nama
    const tdNama = document.createElement("td");
    tdNama.textContent = item.nama;

    // Kolom Kategori
    const tdKategori = document.createElement("td");
    tdKategori.textContent = item.kategori || "-";

    // Kolom Merek
    const tdMerek = document.createElement("td");
    tdMerek.textContent = item.merek || "-";

    // Kolom Harga
    const tdHarga = document.createElement("td");
    tdHarga.textContent = formatRupiah(item.harga);

    // Kolom Stok dengan Penanda Stok Menipis
    const tdStok = document.createElement("td");
    const stokBadge = document.createElement("span");
    stokBadge.className = isLowStock ? "badge badge-warning" : "badge badge-success";
    stokBadge.textContent = isLowStock ? `⚠️ ${item.stok} (Menipis)` : `${item.stok}`;
    tdStok.appendChild(stokBadge);

    // Kolom Stok Min
    const tdMin = document.createElement("td");
    tdMin.textContent = item.stok_minimum;

    // Kolom Lokasi Rak
    const tdRak = document.createElement("td");
    tdRak.textContent = item.lokasi_rak || "-";

    // Kolom Aksi Tombol Alternatif (dapat diklik langsung / layar sentuh)
    const tdAksi = document.createElement("td");
    tdAksi.style.whiteSpace = "nowrap";

    const btnTrx = document.createElement("button");
    btnTrx.className = "btn btn-secondary btn-sm";
    btnTrx.style.marginRight = "6px";
    btnTrx.title = "Barang Masuk / Keluar";
    btnTrx.textContent = "📦 Transaksi";
    btnTrx.addEventListener("click", () => openTransactionModal(item));

    const btnHistory = document.createElement("button");
    btnHistory.className = "btn btn-secondary btn-sm";
    btnHistory.style.marginRight = "6px";
    btnHistory.title = "Riwayat Transaksi";
    btnHistory.textContent = "📜";
    btnHistory.addEventListener("click", () => openHistoryModal(item));

    const btnEdit = document.createElement("button");
    btnEdit.className = "btn btn-secondary btn-sm";
    btnEdit.style.marginRight = "6px";
    btnEdit.title = "Ubah Data";
    btnEdit.textContent = "✏️";
    btnEdit.addEventListener("click", () => openEditModal(item));

    const btnDelete = document.createElement("button");
    btnDelete.className = "btn btn-danger btn-sm";
    btnDelete.title = "Hapus";
    btnDelete.textContent = "🗑️";
    btnDelete.addEventListener("click", () => openDeleteModal(item));

    tdAksi.appendChild(btnTrx);
    tdAksi.appendChild(btnHistory);
    tdAksi.appendChild(btnEdit);
    tdAksi.appendChild(btnDelete);

    tr.appendChild(tdKode);
    tr.appendChild(tdNama);
    tr.appendChild(tdKategori);
    tr.appendChild(tdMerek);
    tr.appendChild(tdHarga);
    tr.appendChild(tdStok);
    tr.appendChild(tdMin);
    tr.appendChild(tdRak);
    tr.appendChild(tdAksi);

    // Event Listener Klik Kanan pada Baris
    tr.addEventListener("contextmenu", (e) => {
      e.preventDefault();
      state.selectedItem = item;
      showContextMenu(e.pageX, e.pageY);
    });

    tbody.appendChild(tr);
  });
}

/**
 * Render tombol paginasi
 */
function renderPagination() {
  const info = document.getElementById("paginationInfo");
  const btnPrev = document.getElementById("btnPrevPage");
  const btnNext = document.getElementById("btnNextPage");

  info.textContent = `Halaman ${state.page} dari ${state.totalPages} (Total ${state.totalItems} sparepart)`;
  btnPrev.disabled = state.page <= 1;
  btnNext.disabled = state.page >= state.totalPages;
}

/**
 * Update kartu ringkasan metrik dashboard
 */
function updateMetrics() {
  document.getElementById("metricTotalItems").textContent = state.totalItems;

  const lowStockCount = state.items.filter(i => i.stok <= i.stok_minimum).length;
  document.getElementById("metricLowStock").textContent = lowStockCount;

  const totalCurrentStock = state.items.reduce((sum, i) => sum + i.stok, 0);
  document.getElementById("metricTotalStock").textContent = totalCurrentStock;
}

/**
 * Custom Context Menu saat Klik Kanan
 */
const contextMenu = document.getElementById("customContextMenu");

function showContextMenu(x, y) {
  contextMenu.style.top = `${y}px`;
  contextMenu.style.left = `${x}px`;
  contextMenu.style.display = "block";
}

function hideContextMenu() {
  if (contextMenu) contextMenu.style.display = "none";
}

document.addEventListener("click", (e) => {
  if (!contextMenu.contains(e.target)) {
    hideContextMenu();
  }
});

// Aksi dari Context Menu
document.getElementById("ctxTrx").addEventListener("click", () => {
  hideContextMenu();
  if (state.selectedItem) openTransactionModal(state.selectedItem);
});

document.getElementById("ctxHistory").addEventListener("click", () => {
  hideContextMenu();
  if (state.selectedItem) openHistoryModal(state.selectedItem);
});

document.getElementById("ctxEdit").addEventListener("click", () => {
  hideContextMenu();
  if (state.selectedItem) openEditModal(state.selectedItem);
});

document.getElementById("ctxDelete").addEventListener("click", () => {
  hideContextMenu();
  if (state.selectedItem) openDeleteModal(state.selectedItem);
});

/**
 * Manajemen Modal Universal
 */
function openModal(id) {
  document.getElementById(id).classList.add("active");
}

function closeModal(id) {
  document.getElementById(id).classList.remove("active");
}

// Modal Tambah Sparepart
function openAddModal() {
  document.getElementById("sparepartModalTitle").textContent = "Tambah Sparepart Baru";
  document.getElementById("sparepartForm").reset();
  document.getElementById("sparepartId").value = "";
  openModal("sparepartModal");
}

// Modal Ubah Sparepart
function openEditModal(item) {
  document.getElementById("sparepartModalTitle").textContent = `Ubah Sparepart: ${item.kode}`;
  document.getElementById("sparepartId").value = item.id;
  document.getElementById("formKode").value = item.kode;
  document.getElementById("formNama").value = item.nama;
  document.getElementById("formKategori").value = item.kategori || "";
  document.getElementById("formMerek").value = item.merek || "";
  document.getElementById("formHarga").value = item.harga;
  document.getElementById("formStokMin").value = item.stok_minimum;
  document.getElementById("formRak").value = item.lokasi_rak || "";
  openModal("sparepartModal");
}

// Simpan Tambah / Ubah
document.getElementById("sparepartForm").addEventListener("submit", async (e) => {
  e.preventDefault();
  const id = document.getElementById("sparepartId").value;
  const isEdit = Boolean(id);

  const payload = {
    kode: document.getElementById("formKode").value.trim(),
    nama: document.getElementById("formNama").value.trim(),
    kategori: document.getElementById("formKategori").value.trim(),
    merek: document.getElementById("formMerek").value.trim(),
    harga: parseInt(document.getElementById("formHarga").value, 10) || 0,
    stok_minimum: parseInt(document.getElementById("formStokMin").value, 10) || 0,
    lokasi_rak: document.getElementById("formRak").value.trim()
  };

  try {
    if (isEdit) {
      await apiRequest(`/spareparts/${id}`, {
        method: "PUT",
        body: JSON.stringify(payload)
      });
      showToast("Data sparepart berhasil diperbarui");
    } else {
      await apiRequest("/spareparts", {
        method: "POST",
        body: JSON.stringify(payload)
      });
      showToast("Sparepart baru berhasil ditambahkan");
    }

    closeModal("sparepartModal");
    loadSpareparts(state.page);
  } catch (err) {
    showToast(err.message, "error");
  }
});

// Modal Transaksi Stok (Barang Masuk / Keluar)
function openTransactionModal(item) {
  state.selectedItem = item;
  document.getElementById("trxSparepartInfo").textContent = `${item.kode} — ${item.nama} (Stok Saat Ini: ${item.stok})`;
  document.getElementById("trxJumlah").value = 1;
  document.getElementById("trxCatatan").value = "";
  document.getElementById("trxTipe").value = "masuk";
  openModal("transactionModal");
}

document.getElementById("transactionForm").addEventListener("submit", async (e) => {
  e.preventDefault();
  if (!state.selectedItem) return;

  const payload = {
    tipe: document.getElementById("trxTipe").value,
    jumlah: parseInt(document.getElementById("trxJumlah").value, 10),
    catatan: document.getElementById("trxCatatan").value.trim()
  };

  try {
    await apiRequest(`/spareparts/${state.selectedItem.id}/transactions`, {
      method: "POST",
      body: JSON.stringify(payload)
    });

    showToast(`Transaksi barang ${payload.tipe} berhasil dicatat!`);
    closeModal("transactionModal");
    loadSpareparts(state.page);
  } catch (err) {
    showToast(err.message, "error");
  }
});

// Modal Riwayat Transaksi
async function openHistoryModal(item) {
  state.selectedItem = item;
  document.getElementById("historyModalTitle").textContent = `Riwayat Transaksi: ${item.kode} — ${item.nama}`;
  const tbody = document.getElementById("historyTableBody");
  tbody.innerHTML = `<tr><td colspan="5" style="text-align:center; padding:20px;">Memuat riwayat transaksi...</td></tr>`;
  openModal("historyModal");

  try {
    const res = await apiRequest(`/spareparts/${item.id}/transactions`);
    const transactions = res.data || [];

    tbody.innerHTML = "";
    if (transactions.length === 0) {
      tbody.innerHTML = `<tr><td colspan="5" class="empty-state">Belum ada transaksi untuk sparepart ini.</td></tr>`;
      return;
    }

    transactions.forEach(trx => {
      const tr = document.createElement("tr");

      const tdTime = document.createElement("td");
      tdTime.textContent = formatDateTime(trx.created_at);

      const tdTipe = document.createElement("td");
      const badge = document.createElement("span");
      badge.className = trx.tipe === "masuk" ? "badge badge-success" : "badge badge-danger";
      badge.textContent = trx.tipe === "masuk" ? "⬇️ Masuk" : "⬆️ Keluar";
      tdTipe.appendChild(badge);

      const tdJumlah = document.createElement("td");
      tdJumlah.style.fontWeight = "bold";
      tdJumlah.textContent = trx.jumlah;

      const tdAdmin = document.createElement("td");
      tdAdmin.textContent = trx.user ? trx.user.username : `User #${trx.user_id}`;

      const tdCatatan = document.createElement("td");
      tdCatatan.textContent = trx.catatan || "-";

      tr.appendChild(tdTime);
      tr.appendChild(tdTipe);
      tr.appendChild(tdJumlah);
      tr.appendChild(tdAdmin);
      tr.appendChild(tdCatatan);
      tbody.appendChild(tr);
    });
  } catch (err) {
    tbody.innerHTML = `<tr><td colspan="5" style="text-align:center; color:var(--danger); padding:20px;">${escapeHtml(err.message)}</td></tr>`;
  }
}

// Modal Konfirmasi Hapus
function openDeleteModal(item) {
  state.selectedItem = item;
  document.getElementById("deleteSparepartText").textContent = `Apakah Anda yakin ingin menghapus sparepart "${item.kode} — ${item.nama}"?`;
  openModal("deleteModal");
}

document.getElementById("btnConfirmDelete").addEventListener("click", async () => {
  if (!state.selectedItem) return;

  try {
    await apiRequest(`/spareparts/${state.selectedItem.id}`, {
      method: "DELETE"
    });

    showToast("Sparepart berhasil dihapus");
    closeModal("deleteModal");
    loadSpareparts(state.page);
  } catch (err) {
    showToast(err.message, "error");
  }
});

// Escape HTML helper
function escapeHtml(str) {
  if (!str) return "";
  return String(str)
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;");
}

// Inisialisasi Event Handlers
document.addEventListener("DOMContentLoaded", async () => {
  // Tombol logout di navbar (pengganti onclick inline)
  document.getElementById("btnLogout").addEventListener("click", logout);

  // Tombol tutup/batal pada semua modal (pengganti onclick inline)
  document.querySelectorAll("[data-close]").forEach((btn) => {
    btn.addEventListener("click", () => closeModal(btn.dataset.close));
  });

  const user = await checkAuth(true);
  if (!user) return;

  // Tombol Tambah
  document.getElementById("btnOpenAddModal").addEventListener("click", openAddModal);

  // Pencarian (Debounce 350ms)
  let searchTimeout;
  document.getElementById("searchInput").addEventListener("input", (e) => {
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(() => {
      state.search = e.target.value.trim();
      loadSpareparts(1);
    }, 350);
  });

  // Filter Kategori
  document.getElementById("filterKategori").addEventListener("change", (e) => {
    state.kategori = e.target.value;
    loadSpareparts(1);
  });

  // Tombol Paginasi
  document.getElementById("btnPrevPage").addEventListener("click", () => {
    if (state.page > 1) loadSpareparts(state.page - 1);
  });

  document.getElementById("btnNextPage").addEventListener("click", () => {
    if (state.page < state.totalPages) loadSpareparts(state.page + 1);
  });

  // Muat data awal
  loadSpareparts(1);
});

