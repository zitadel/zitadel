/// <reference types="cypress" />
import React from 'react';
import type { MountOptions } from './types';
export declare function mount(jsx: React.ReactNode, options?: MountOptions, rerenderKey?: string): Cypress.Chainable<import("./types").MountReturn>;
export declare function unmount(options?: {
    log: boolean;
}): Cypress.Chainable<JQuery<HTMLElement>>;
