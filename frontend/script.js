
/**
 * Frontend Order Lookup (fetch-based)
 * Expects a backend endpoint: GET /order/:id -> returns JSON order
 */

const API_BASE = window.API_BASE || "";

const qs = s => document.querySelector(s);
const orderIdInput   = qs("#orderIdInput");
const orderSection   = qs("#orderSection");
const orderContent   = qs("#orderContent");
const loadingOverlay = qs("#loadingOverlay");
const errorBox       = qs("#errorMessage");
const errorText      = qs("#errorText");

function showLoading(on = true) {
  loadingOverlay?.classList.toggle("hidden", !on);
}
function showError(msg) {
  if (!errorBox || !errorText) return;
  errorText.textContent = msg || "Unknown error";
  errorBox.classList.remove("hidden");
  setTimeout(() => errorBox.classList.add("hidden"), 5000);
}
function hideOrder() {
  orderSection.style.display = "none";
  orderContent.innerHTML = "";
}
window.hideOrder = hideOrder;

function fmtMoney(amount, currency, locale = "ru-RU") {
  try {
    return new Intl.NumberFormat(locale, { style: "currency", currency }).format(amount);
  } catch {
    return `${amount} ${currency || ""}`.trim();
  }
}

function fmtDate(iso, locale = "ru-RU") {
  try {
    const d = new Date(iso);
    return d.toLocaleString(locale, {
      year: "numeric", month: "2-digit", day: "2-digit",
      hour: "2-digit", minute: "2-digit", second: "2-digit"
    });
  } catch {
    return iso;
  }
}

const STATUS_CLASS = {
  1: "status-1",
  2: "status-2",
  3: "status-3",
  4: "status-4",
  5: "status-5"
};

function renderOrder(order) {
  const delivery = order.delivery || {};
  const payment  = order.payment  || {};
  const items    = Array.isArray(order.items) ? order.items : [];

  const deliveryHTML = `
    <div class="order-detail-section">
      <h4>Delivery</h4>
      <div class="order-detail-item"><span class="order-detail-label">Name</span><span class="order-detail-value">${escapeHtml(delivery.name || "-")}</span></div>
      <div class="order-detail-item"><span class="order-detail-label">Phone</span><span class="order-detail-value">${escapeHtml(delivery.phone || "-")}</span></div>
      <div class="order-detail-item"><span class="order-detail-label">Address</span><span class="order-detail-value">${escapeHtml(delivery.address || "-")}</span></div>
      <div class="order-detail-item"><span class="order-detail-label">City/Region</span><span class="order-detail-value"><span class="city-badge">${escapeHtml(delivery.city || "-")}</span> • ${escapeHtml(delivery.region || "-")}</span></div>
      <div class="order-detail-item"><span class="order-detail-label">ZIP</span><span class="order-detail-value">${escapeHtml(delivery.zip || "-")}</span></div>
      <div class="order-detail-item"><span class="order-detail-label">Email</span><span class="order-detail-value">${escapeHtml(delivery.email || "-")}</span></div>
    </div>
  `;

  const paymentHTML = `
    <div class="order-detail-section">
      <h4>Payment</h4>
      <div class="order-detail-item"><span class="order-detail-label">Transaction</span><span class="order-detail-value">${escapeHtml(payment.transaction || "-")}</span></div>
      <div class="order-detail-item"><span class="order-detail-label">Provider</span><span class="order-detail-value">${escapeHtml(payment.provider || "-")}</span></div>
      <div class="order-detail-item"><span class="order-detail-label">Bank</span><span class="order-detail-value">${escapeHtml(payment.bank || "-")}</span></div>
      <div class="order-detail-item"><span class="order-detail-label">Currency</span><span class="order-detail-value">${escapeHtml(payment.currency || "-")}</span></div>
      <div class="order-detail-item"><span class="order-detail-label">Goods Total</span><span class="order-detail-value">${fmtMoney(payment.goods_total ?? 0, payment.currency || "RUB")}</span></div>
      <div class="order-detail-item"><span class="order-detail-label">Delivery Cost</span><span class="order-detail-value">${fmtMoney(payment.delivery_cost ?? 0, payment.currency || "RUB")}</span></div>
      <div class="order-detail-item"><span class="order-detail-label">Custom Fee</span><span class="order-detail-value">${fmtMoney(payment.custom_fee ?? 0, payment.currency || "RUB")}</span></div>
      <div class="order-detail-item"><span class="order-detail-label">Amount</span><span class="order-detail-value total-amount">${fmtMoney(payment.amount ?? 0, payment.currency || "RUB")}</span></div>
    </div>
  `;

  const summaryHTML = `
    <div class="order-detail-section">
      <h4>Order</h4>
      <div class="order-detail-item"><span class="order-detail-label">Order ID</span><span class="order-detail-value order-id">${escapeHtml(order.order_uid)}</span></div>
      <div class="order-detail-item"><span class="order-detail-label">Track #</span><span class="order-detail-value">${escapeHtml(order.track_number || "-")}</span></div>
      <div class="order-detail-item"><span class="order-detail-label">Entry</span><span class="order-detail-value">${escapeHtml(order.entry || "-")}</span></div>
      <div class="order-detail-item"><span class="order-detail-label">Locale</span><span class="order-detail-value">${escapeHtml(order.locale || "-")}</span></div>
      <div class="order-detail-item"><span class="order-detail-label">Delivery Service</span><span class="order-detail-value">${escapeHtml(order.delivery_service || "-")}</span></div>
      <div class="order-detail-item"><span class="order-detail-label">Created</span><span class="order-detail-value">${fmtDate(order.date_created || "-")}</span></div>
    </div>
  `;

  const itemsHTML = `
    <div class="order-detail-section" style="grid-column: 1 / -1;">
      <h4>Items</h4>
      <div class="items-list">
        ${items.map(it => `
          <div class="item-card">
            <div class="item-header">
              <span class="item-name">${escapeHtml(it.name || "-")}</span>
              <span class="item-price">${fmtMoney(it.total_price ?? it.price ?? 0, payment.currency || "RUB")}</span>
            </div>
            <div class="item-details">
              <div>Brand: <strong>${escapeHtml(it.brand || "-")}</strong></div>
              <div>Size: <strong>${escapeHtml(String(it.size ?? "-"))}</strong></div>
              <div>RID: <strong>${escapeHtml(it.rid || "-")}</strong></div>
              <div>Track: <strong>${escapeHtml(it.track_number || "-")}</strong></div>
              <div>Sale: <strong>${escapeHtml(String(it.sale ?? 0))}%</strong></div>
              <div>Status: <span class="status-badge ${STATUS_CLASS[it.status] || ""}">${escapeHtml(String(it.status ?? "-"))}</span></div>
            </div>
          </div>
        `).join("")}
      </div>
    </div>
  `;

  orderContent.innerHTML = `
    <div class="order-detail-grid">
      ${summaryHTML}
      ${deliveryHTML}
      ${paymentHTML}
      ${itemsHTML}
    </div>
  `;
  orderSection.style.display = "block";
}

function escapeHtml(x) {
  return String(x)
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;");
}

async function lookupOrder() {
  const id = (orderIdInput.value || "").trim();
  if (!id) {
    showError("Введите ID заказа (например, b1)");
    return;
  }

  showLoading(true);
  try {
    const url = `${API_BASE}/order/${encodeURIComponent(id)}`;
    const res = await fetch(url, { headers: { "Accept": "application/json" } });
    if (!res.ok) {
      if (res.status === 404) throw new Error("Заказ не найден");
      throw new Error(`Ошибка запроса: ${res.status}`);
    }
    const data = await res.json();
    renderOrder(data);
  } catch (e) {
    hideOrder();
    showError(e.message || "Не удалось получить заказ");
  } finally {
    showLoading(false);
  }
}
window.lookupOrder = lookupOrder;

orderIdInput?.addEventListener("keydown", (e) => {
  if (e.key === "Enter") {
    lookupOrder();
  }
});
