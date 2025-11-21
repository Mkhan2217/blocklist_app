global.confirm = jest.fn(() => true);

// Mock DOM elements for toast and form
global.document.body.innerHTML = `
  <input id="searchPhone" type="tel" value="+1234567890" />
  <div id="searchResult"></div>
  <div id="toast-container"></div>
  <form id="blockForm">
    <input name="reason" value="test" />
    <input name="store_location" value="store" />
    <input name="check_amount" value="100" />
    <input name="notes" value="notes" />
    <button type="submit">Submit</button>
  </form>
`;

// Mock fetch globally for all tests
global.fetch = jest.fn(() =>
  Promise.resolve({ ok: true, json: () => Promise.resolve({}) })
);
