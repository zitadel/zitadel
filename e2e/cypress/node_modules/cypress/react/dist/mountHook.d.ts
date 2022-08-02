/// <reference types="cypress" />
declare type MountHookResult<T> = {
    readonly current: T | null | undefined;
    readonly error: Error | null;
};
/**
 * Mounts a React hook function in a test component for testing.
 *
 */
export declare const mountHook: <T>(hookFn: (...args: any[]) => T) => Cypress.Chainable<MountHookResult<T>>;
export {};
