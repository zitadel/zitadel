/// <reference types="cypress" />
import Vue from 'vue';
import { VueTestUtilsConfigOptions, Wrapper } from '@vue/test-utils';
import { StyleOptions } from '@cypress/mount-utils';
/**
 * Type for component passed to "mount"
 *
 * @interface VueComponent
 * @example
 *  import Hello from './Hello.vue'
 *         ^^^^^ this type
 *  mount(Hello)
 */
declare type VueComponent = Vue.ComponentOptions<any> | Vue.VueConstructor;
/**
 * Options to pass to the component when creating it, like
 * props.
 *
 * @interface ComponentOptions
 */
declare type ComponentOptions = Record<string, unknown>;
declare type VueLocalComponents = Record<string, VueComponent>;
declare type VueFilters = {
    [key: string]: (value: string) => string;
};
declare type VueMixin = unknown;
declare type VueMixins = VueMixin | VueMixin[];
declare type VuePluginOptions = unknown;
declare type VuePlugin = unknown | [unknown, VuePluginOptions];
/**
 * A single Vue plugin or a list of plugins to register
 */
declare type VuePlugins = VuePlugin[];
/**
 * Additional Vue services to register while mounting the component, like
 * local components, plugins, etc.
 *
 * @interface MountOptionsExtensions
 * @see https://github.com/cypress-io/cypress/tree/master/npm/vue#examples
 */
interface MountOptionsExtensions {
    /**
     * Extra local components
     *
     * @memberof MountOptionsExtensions
     * @see https://github.com/cypress-io/cypress/tree/master/npm/vue#examples
     * @example
     *  import Hello from './Hello.vue'
     *  // imagine Hello needs AppComponent
     *  // that it uses in its template like <app-component ... />
     *  // during testing we can replace it with a mock component
     *  const appComponent = ...
     *  const components = {
     *    'app-component': appComponent
     *  },
     *  mount(Hello, { extensions: { components }})
     */
    components?: VueLocalComponents;
    /**
     * Optional Vue filters to install while mounting the component
     *
     * @memberof MountOptionsExtensions
     * @see https://github.com/cypress-io/cypress/tree/master/npm/vue#examples
     * @example
     *  const filters = {
     *    reverse: (s) => s.split('').reverse().join(''),
     *  }
     *  mount(Hello, { extensions: { filters }})
     */
    filters?: VueFilters;
    /**
     * Optional Vue mixin(s) to install when mounting the component
     *
     * @memberof MountOptionsExtensions
     * @alias mixins
     * @see https://github.com/cypress-io/cypress/tree/master/npm/vue#examples
     */
    mixin?: VueMixins;
    /**
     * Optional Vue mixin(s) to install when mounting the component
     *
     * @memberof MountOptionsExtensions
     * @alias mixin
     * @see https://github.com/cypress-io/cypress/tree/master/npm/vue#examples
     */
    mixins?: VueMixins;
    /**
     * A single plugin or multiple plugins.
     *
     * @see https://github.com/cypress-io/cypress/tree/master/npm/vue#examples
     * @alias plugins
     * @memberof MountOptionsExtensions
     */
    use?: VuePlugins;
    /**
     * A single plugin or multiple plugins.
     *
     * @see https://github.com/cypress-io/cypress/tree/master/npm/vue#examples
     * @alias use
     * @memberof MountOptionsExtensions
     */
    plugins?: VuePlugins;
}
/**
 * Options controlling how the component is going to be mounted,
 * including global Vue plugins and extensions.
 *
 * @interface MountOptions
 */
interface MountOptions {
    /**
     * Vue instance to use.
     *
     * @deprecated
     * @memberof MountOptions
     */
    vue: unknown;
    /**
     * Extra Vue plugins, mixins, local components to register while
     * mounting this component
     *
     * @memberof MountOptions
     * @see https://github.com/cypress-io/cypress/tree/master/npm/vue#examples
     */
    extensions: MountOptionsExtensions;
}
/**
 * Utility type for union of options passed to "mount(..., options)"
 */
declare type MountOptionsArgument = Partial<ComponentOptions & MountOptions & StyleOptions & VueTestUtilsConfigOptions>;
declare global {
    namespace Cypress {
        interface Cypress {
            /**
             * Mounted Vue instance is available under Cypress.vue
             * @memberof Cypress
             * @example
             *  mount(Greeting)
             *  .then(() => {
             *    Cypress.vue.message = 'Hello There'
             *  })
             *  // new message is displayed
             *  cy.contains('Hello There').should('be.visible')
             */
            vue: Vue;
            vueWrapper: Wrapper<Vue>;
        }
    }
}
/**
 * Mounts a Vue component inside Cypress browser.
 * @param {object} component imported from Vue file
 * @example
 *  import Greeting from './Greeting.vue'
 *  import { mount } from '@cypress/vue2'
 *  it('works', () => {
 *    // pass props, additional extensions, etc
 *    mount(Greeting, { ... })
 *    // use any Cypress command to test the component
 *    cy.get('#greeting').should('be.visible')
 *  })
 */
export declare const mount: (component: VueComponent, optionsOrProps?: MountOptionsArgument) => Cypress.Chainable<Cypress.AUTWindow>;
/**
 * Helper function for mounting a component quickly in test hooks.
 * @example
 *  import {mountCallback} from '@cypress/vue2'
 *  beforeEach(mountVue(component, options))
 */
export declare const mountCallback: (component: VueComponent, options?: MountOptionsArgument) => () => Cypress.Chainable<Cypress.AUTWindow>;
export {};
