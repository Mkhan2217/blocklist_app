/* ---------------------------------------------------------
   Phone Validation
--------------------------------------------------------- */
function validatePhoneInput(input) {
  const formattedNumber = formatPhoneNumber(input.value);
  const isValid = validatePhoneNumber(formattedNumber);
  const errorMsg = input.parentElement.querySelector(".error-message");
  input.value = formattedNumber;

  // Enable/Disable related button
  let button;
  if (input.id === "searchPhone") {
    button = document.querySelector(".search-box button");
  } else if (input.form) {
    button = input.form.querySelector('button[type="submit"]');
  }

  if (isValid) {
    input.classList.remove("invalid-input");
    input.classList.add("valid-input");
    if (errorMsg) errorMsg.style.display = "none";
    if (button) button.disabled = false;
  } else {
    input.classList.remove("valid-input");
    input.classList.add("invalid-input");
    if (errorMsg) errorMsg.style.display = "block";
    if (button) button.disabled = true;
  }

  return isValid;
}

function validatePhoneNumber(phone) {
  const phoneRegex = /^\+[1-9][0-9]{9,14}$/; // + followed by 10–15 digits
  return phoneRegex.test(phone);
}

function formatPhoneNumber(phone) {
  let cleaned = phone.replace(/[^\d+]/g, "");

  if (cleaned.match(/^\d/)) {
    cleaned = "+" + cleaned;
  }

  if (cleaned.startsWith("+0")) {
    cleaned = "+" + cleaned.substring(2);
  }

  return cleaned;
}

/* ---------------------------------------------------------
   Search Blocked Number
--------------------------------------------------------- */
async function searchPhone() {
  const phoneInput = document.getElementById("searchPhone");
  const result = document.getElementById("searchResult");
  const searchBtn = document.querySelector(".search-box button");

  if (!validatePhoneNumber(phoneInput.value)) {
    result.innerHTML =
      '<p style="color: #ff6b6b;">Please enter a valid phone number</p>';
    result.style.display = "block";
    return;
  }

  searchBtn.disabled = true;

  try {
    const formattedPhone = formatPhoneNumber(phoneInput.value);

    const response = await fetch(
      `/api/blocklist?phone=${encodeURIComponent(formattedPhone)}`
    );

    if (response.ok) {
      const data = await response.json();
      result.innerHTML = `
        <h3>⚠️ WARNING: Blocked Number</h3>
        <p>Phone: ${data?.phoneNumber}</p>
        <p>Location: ${data?.storeLocation}</p>
        <p>Reason: ${data?.reason}</p>
        <p>Date: ${formatUSADate(data?.incidentDate)}</p>
      `;
    } else if (response.status === 404) {
      result.innerHTML = "<p>Number not found in blocklist</p>";
    } else {
      throw new Error("Search failed");
    }
  } catch (error) {
    console.error("Error:", error);
    result.innerHTML = '<p style="color: #ff6b6b;">Error checking number</p>';
  } finally {
    result.style.display = "block";
    searchBtn.disabled = false;
  }
}

/* ---------------------------------------------------------
   Block Number (SAVE) — POST /api/blocklist
--------------------------------------------------------- */
async function handleBlockSubmit(e) {
  e.preventDefault();

  const phone = document.getElementById("phoneInput").value.trim();
  const form = e.target;

  const payload = {
    phoneNumber: phone,
    reason: form.reason.value,
    storeLocation: form.store_location.value,
    checkAmount: parseFloat(form.check_amount.value),
    notes: form.notes.value,
  };

  if (!validatePhoneNumber(phone)) {
    showToast("Please enter a valid international phone number", "error");

    return;
  }

  try {
    const res = await fetch("/api/blocklist", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    if (res.ok) {
      showToast("Number blocked successfully!", "success");
      resetSearchBox();
      setTimeout(() => location.reload(), 800);
    } else {
      showToast("Failed to block number.", "error");
    }
  } catch (err) {
    console.error("Error saving:", err);
    showToast("Server error while saving", "error");
  }
}

/* ---------------------------------------------------------
   Unblock Number — DELETE /api/blocklist
--------------------------------------------------------- */
async function unblockNumber(phone, buttonElement) {
  if (!confirm(`Unblock number: ${phone}?`)) return;

  try {
    const response = await fetch(
      `/api/blocklist?phone=${encodeURIComponent(phone)}`,
      { method: "DELETE" }
    );

    if (response.ok) {
      const row = buttonElement.closest("tr");
      if (row) row.remove();
      showToast("Number unblocked successfully!", "success");
      resetSearchBox();
    }
  } catch (err) {
    console.error("Unblock error:", err);
  }
}

/* ---------------------------------------------------------
   Format Date
--------------------------------------------------------- */
function formatUSADate(dateString) {
  const date = new Date(dateString);
  if (isNaN(date)) return dateString;

  return date.toLocaleString("en-US", {
    month: "2-digit",
    day: "2-digit",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    hour12: true,
  });
}

/* ---------------------------------------------------------
   DOM Events
--------------------------------------------------------- */
document.addEventListener("DOMContentLoaded", () => {
  // Phone validation
  const phoneInputs = document.querySelectorAll('input[type="tel"]');
  phoneInputs.forEach((input) =>
    input.addEventListener("input", () => validatePhoneInput(input))
  );

  // Format dates in table
  document.querySelectorAll("tbody tr").forEach((row) => {
    const dateCell = row.children[1];
    if (dateCell && dateCell.textContent.includes("T")) {
      dateCell.textContent = formatUSADate(dateCell.textContent.trim());
    }
  });

  // Save / Block form submission
  const blockForm = document.getElementById("blockForm");
  if (blockForm) {
    blockForm.addEventListener("submit", handleBlockSubmit);
  }

  // Unblock button event
  document.addEventListener("click", (e) => {
    const btn = e.target.closest(".unblock-btn-table");
    if (btn) {
      const phone = btn.getAttribute("data-phone");
      unblockNumber(phone, btn);
    }
  });
});
function showToast(message, type = "success") {
  const container = document.getElementById("toast-container");

  const toast = document.createElement("div");
  toast.classList.add("toast");

  if (type === "success") toast.classList.add("toast-success");
  if (type === "error") toast.classList.add("toast-error");
  if (type === "warning") toast.classList.add("toast-warning");

  toast.textContent = message;

  container.appendChild(toast);

  setTimeout(() => {
    toast.remove();
  }, 3000);
}
function resetSearchBox() {
  const searchInput = document.getElementById("searchPhone");
  const resultBox = document.getElementById("searchResult");

  if (searchInput) searchInput.value = "";
  if (resultBox) {
    resultBox.innerHTML = "";
    resultBox.style.display = "none";
  }
}
// Export for Jest testing
if (typeof module !== "undefined") {
  module.exports = {
    formatPhoneNumber,
    validatePhoneNumber,
    searchPhone,
    handleBlockSubmit,
    unblockNumber,
    formatUSADate,
    validatePhoneInput,
    resetSearchBox,
    showToast,
  };
}