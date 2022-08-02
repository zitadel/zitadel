/// <reference types="cypress" />
import React from 'react';
import type { MountOptions, UnmountArgs } from '@cypress/react';
export declare function mount(jsx: React.ReactNode, options?: MountOptions, rerenderKey?: string): Cypress.Chainable<import("@cypress/react").MountReturn>;
export declare function unmount(options?: UnmountArgs): Cypress.Chainable<JQuery<HTMLElement>>;
