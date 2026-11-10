// Mock for server-only package in tests
// The real package throws an error when imported outside of a React Server Component context
// This empty mock allows tests to run without that restriction
export {};
