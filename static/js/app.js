function validatePhoneInput(input) {
  const formattedNumber = formatPhoneNumber(input.value);
  const isValid = validatePhoneNumber(formattedNumber);
  const errorMsg = input.parentElement.querySelector(".error-message");
  input.value = formattedNumber;

  // Get the relevant button (either submit or check button)
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
  // Match database constraint: + followed by digit 1-9, then 9-14 more digits
  const phoneRegex = /^\+[1-9][0-9]{9,14}$/;
  return phoneRegex.test(phone);
}

function formatPhoneNumber(phone) {
  // Remove all spaces and special characters except + and digits
  let cleaned = phone.replace(/[^\d+]/g, "");

  // If starts with digits, add +
  if (cleaned.match(/^\d/)) {
    cleaned = "+" + cleaned;
  }

  // Ensure first digit after + is 1-9
  if (cleaned.startsWith("+0")) {
    cleaned = "+" + cleaned.substring(2);
  }

  return cleaned;
}

async function searchPhone() {
  const phoneInput = document.getElementById("searchPhone");
  const result = document.getElementById("searchResult");
  const searchBtn = document.querySelector(".search-box button");

  if (!validatePhoneNumber(phoneInput.value)) {
    result.innerHTML =
      '<p style="color: #ff6b6b;">Please enter a valid phone number </p>';
    result.style.display = "block";
    return;
  }

  searchBtn.disabled = true;
  try {
    const formattedPhone = formatPhoneNumber(phoneInput.value);
    const response = await fetch(
      `/search?phone=${encodeURIComponent(formattedPhone)}`
    );

    if (response.ok) {
      const data = await response.json();
      result.innerHTML = `
      <h3>⚠️ WARNING: Blocked Number</h3>
      <p>Phone: ${data.phone_number}</p>
      <p>Location: ${data.store_location}</p>
      <p>Reason: ${data.reason}</p>
      <p>Date: ${data.incident_date}</p>
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
async function unblockNumber(phone, buttonElement) {
  if (!confirm(`Unblock number: ${phone}?`)) return;

  try {
    const response = await fetch(
      `/unblock?phone=${encodeURIComponent(phone)}`,
      {
        method: "DELETE",
      }
    );

    if (response.ok) {
      // Remove only the row of this button
      const row = buttonElement.closest("tr");
      if (row) row.remove();
    } else if (response.status === 404) {
    } else {
    }
  } catch (err) {
    console.error(err);
  }
}
function formatUSADate(dateString) {
  const date = new Date(dateString);
  if (isNaN(date)) return dateString;

  const options = {
    month: "2-digit",
    day: "2-digit",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    hour12: true,
  };

  return date.toLocaleString("en-US", options);
}

document.addEventListener("DOMContentLoaded", () => {
  // phone validation
  const phoneInputs = document.querySelectorAll('input[type="tel"]');
  phoneInputs.forEach((input) =>
    input.addEventListener("input", () => validatePhoneInput(input))
  );

  // Format CreatedAt column to USA format
  document.querySelectorAll("tbody tr").forEach((row) => {
    const dateCell = row.children[1]; // 2nd column = CreatedAt
    if (dateCell && dateCell.textContent.includes("T")) {
      dateCell.textContent = formatUSADate(dateCell.textContent.trim());
    }
  });

  // unblock button click
  document.addEventListener("click", (e) => {
    const btn = e.target.closest(".unblock-btn-table");
    if (btn) {
      const phone = btn.getAttribute("data-phone");
      unblockNumber(phone, btn);
    }
  });
});
