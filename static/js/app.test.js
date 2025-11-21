const {
  formatPhoneNumber,
  validatePhoneNumber,
  searchPhone,
} = require("./app");

describe("Phone Number Validation", () => {
  test("formatPhoneNumber formats numbers correctly", () => {
    expect(formatPhoneNumber("1234567890")).toBe("+1234567890");
    expect(formatPhoneNumber("+1234567890")).toBe("+1234567890");
    expect(formatPhoneNumber("91987654321")).toBe("+91987654321");
  });

  test("validatePhoneNumber validates correctly", () => {
    expect(validatePhoneNumber("+1234567890")).toBeTruthy();
    expect(validatePhoneNumber("123")).toBeFalsy();
    expect(validatePhoneNumber("abc")).toBeFalsy();
  });
});

describe("API Calls", () => {
  test("searchPhone handles success response", async () => {
    global.fetch = jest.fn(() =>
      Promise.resolve({
        ok: true,
        json: () =>
          Promise.resolve({
            phoneNumber: "+1234567890",
            reason: "test",
          }),
      })
    );

    document.body.innerHTML = `
      <input id="searchPhone" type="tel" value="+1234567890" />
      <div id="searchResult"></div>
      <div class="search-box"><button></button></div>
    `;

    await searchPhone();
    expect(fetch).toHaveBeenCalled();
  });

  test("searchPhone handles 404 response", async () => {
    global.fetch = jest.fn(() =>
      Promise.resolve({
        ok: false,
        status: 404,
      })
    );

    document.body.innerHTML = `
      <input id="searchPhone" type="tel" value="+1234567890" />
      <div id="searchResult"></div>
      <div class="search-box"><button></button></div>
    `;

    await searchPhone();
    expect(fetch).toHaveBeenCalled();
  });
});
