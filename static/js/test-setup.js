// Mock DOM elements
global.document = {
    getElementById: jest.fn(() => ({
        value: '+1234567890',
        style: {},
        classList: {
            add: jest.fn(),
            remove: jest.fn()
        }
    }))
};

// Mock fetch API
global.fetch = jest.fn();
