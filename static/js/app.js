function validatePhoneInput(input) {
  const formattedNumber = formatPhoneNumber(input.value);
  const isValid = validatePhoneNumber(formattedNumber);
  const errorMsg = input.parentElement.querySelector(".error-message");
  // debugger
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

        <button id="unblockBtn">Unblock Number</button>
    `;

      // handle unblock button click
      document
        .getElementById("unblockBtn")
        .addEventListener("click", async () => {
          if (!confirm("Are you sure you want to unblock this number?")) return;

          const unblockRes = await fetch(
            `/unblock?phone=${encodeURIComponent(data.phone_number)}`,
            {
              method: "DELETE",
            }
          );

          if (unblockRes.ok) {
            alert("✅ Number successfully unblocked!");
            result.innerHTML = "";
            phoneInput.value = "";
          } else {
            alert("❌ Failed to unblock number");
          }
        });
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
