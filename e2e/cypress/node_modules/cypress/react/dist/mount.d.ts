/// <reference types="cypress" />
import * as React from 'react';
import * as ReactDOM from 'react-dom';
import { StyleOptions } from '@cypress/mount-utils';
/**
 * Mount a React component in a blank document; register it as an alias
 * To access: use an alias or original component reference
 * @function   mount
 * @param      {React.ReactElement}  jsx - component to mount
 * @param      {MountOptions}  [options] - options, like alias, styles
 * @see https://github.com/bahmutov/@cypress/react
 * @see https://glebbahmutov.com/blog/my-vision-for-component-tests/
 * @example
 ```
  import Hello from './hello.jsx'
  import { mount } from '@cypress/react'
  it('works', () => {
    mount(<Hello onClick={cy.stub()} />)
    // use Cypress commands
    cy.contains('Hello').click()
  })
 ```
 **/
export declare const mount: (jsx: React.ReactNode, options?: MountOptions) => globalThis.Cypress.Chainable<MountReturn>;
/**
 * Removes the mounted component. Notice this command automatically
 * queues up the `unmount` into Cypress chain, thus you don't need `.then`
 * to call it.
 * @see https://github.com/cypress-io/cypress/tree/develop/npm/react/cypress/component/basic/unmount
 * @example
  ```
  import { mount, unmount } from '@cypress/react'
  it('works', () => {
    mount(...)
    // interact with the component using Cypress commands
    // whenever you want to unmount
    unmount()
  })
  ```
 */
export declare const unmount: (options?: {
    log: boolean;
}) => globalThis.Cypress.Chainable<JQuery<HTMLElement>>;
/**
 * Creates new instance of `mount` function with default options
 * @function   createMount
 * @param      {MountOptions}  [defaultOptions] - defaultOptions for returned `mount` function
 * @returns    new instance of `mount` with assigned options
 * @example
 * ```
 * import Hello from './hello.jsx'
 * import { createMount } from '@cypress/react'
 *
 * const mount = createMount({ strict: true, cssFile: 'path/to/any/css/file.css' })
 *
 * it('works', () => {
 *   mount(<Hello />)
 *   // use Cypress commands
 *   cy.get('button').should('have.css', 'color', 'rgb(124, 12, 109)')
 * })
 ```
 **/
export declare const createMount: (defaultOptions: MountOptions) => (element: React.ReactElement, options?: Partial<StyleOptions & MountReactComponentOptions> | undefined) => globalThis.Cypress.Chainable<MountReturn>;
/** @deprecated Should be removed in the next major version */
export default mount;
export interface ReactModule {
    name: string;
    type: string;
    location: string;
    source: string;
}
export interface MountReactComponentOptions {
    alias: string;
    ReactDom: typeof ReactDOM;
    /**
     * Log the mounting command into Cypress Command Log,
     * true by default.
     */
    log: boolean;
    /**
     * Render component in React [strict mode](https://reactjs.org/docs/strict-mode.html)
     * It activates additional checks and warnings for child components.
     */
    strict: boolean;
}
export declare type MountOptions = Partial<StyleOptions & MountReactComponentOptions>;
export interface MountReturn {
    /**
     * The component that was rendered.
     */
    component: React.ReactNode;
    /**
     * Rerenders the specified component with new props. This allows testing of components that store state (`setState`)
     * or have asynchronous updates (`useEffect`, `useLayoutEffect`).
     */
    rerender: (component: React.ReactNode) => globalThis.Cypress.Chainable<MountReturn>;
    /**
     * Removes the mounted component.
     * @see `unmount`
     */
    unmount: () => globalThis.Cypress.Chainable<JQuery<HTMLElement>>;
}
/**
 * The `type` property from the transpiled JSX object.
 * @example
 * const { type } = React.createElement('div', null, 'Hello')
 * const { type } = <div>Hello</div>
 */
export interface JSX extends Function {
    displayName: string;
}
export declare namespace Cypress {
    interface Cypress {
        stylesCache: any;
        React: string;
        ReactDOM: string;
        Styles: string;
        modules: ReactModule[];
    }
    interface Chainable<Subject> {
        state: (key: string) => any;
        /**
         * Mount a React component in a blank document; register it as an alias
         * To access: use an alias or original component reference
         *  @function   cy.mount
         *  @param      {Object}  jsx - component to mount
         *  @param      {string}  [Component] - alias to use later
         *  @example
        ```
        import Hello from './hello.jsx'
        // mount and access by alias
        cy.mount(<Hello />, 'Hello')
        // using default alias
        cy.get('@Component')
        // using specified alias
        cy.get('@Hello').its('state').should(...)
        // using original component
        cy.get(Hello)
        ```
        **/
        get<S = any>(alias: string | symbol | Function, options?: Partial<any>): Chainable<any>;
    }
}
