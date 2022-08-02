
/**
 * @cypress/react18 v0.0.0-development
 * (c) 2022 Cypress.io
 * Released under the MIT License
 */

import ReactDOM from 'react-dom/client';
import * as React from 'react';
import 'react-dom';

/******************************************************************************
Copyright (c) Microsoft Corporation.

Permission to use, copy, modify, and/or distribute this software for any
purpose with or without fee is hereby granted.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY
AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM
LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR
OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
PERFORMANCE OF THIS SOFTWARE.
***************************************************************************** */

var __assign = function() {
    __assign = Object.assign || function __assign(t) {
        for (var s, i = 1, n = arguments.length; i < n; i++) {
            s = arguments[i];
            for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p)) t[p] = s[p];
        }
        return t;
    };
    return __assign.apply(this, arguments);
};

var cachedDisplayNames = new WeakMap();
/**
 * Gets the display name of the component when possible.
 * @param type {JSX} The type object returned from creating the react element.
 * @param fallbackName {string} The alias, or fallback name to use when the name cannot be derived.
 * @link https://github.com/facebook/react-devtools/blob/master/backend/getDisplayName.js
 */
function getDisplayName(type, fallbackName) {
    if (fallbackName === void 0) { fallbackName = 'Unknown'; }
    var nameFromCache = cachedDisplayNames.get(type);
    if (nameFromCache != null) {
        return nameFromCache;
    }
    var displayName = null;
    // The displayName property is not guaranteed to be a string.
    // It's only safe to use for our purposes if it's a string.
    // github.com/facebook/react-devtools/issues/803
    if (typeof type.displayName === 'string') {
        displayName = type.displayName;
    }
    if (!displayName) {
        displayName = type.name || fallbackName;
    }
    // Facebook-specific hack to turn "Image [from Image.react]" into just "Image".
    // We need displayName with module name for error reports but it clutters the DevTools.
    var match = displayName.match(/^(.*) \[from (.*)\]$/);
    if (match) {
        var componentName = match[1];
        var moduleName = match[2];
        if (componentName && moduleName) {
            if (moduleName === componentName ||
                moduleName.startsWith(componentName + ".")) {
                displayName = componentName;
            }
        }
    }
    try {
        cachedDisplayNames.set(type, displayName);
    }
    catch (e) {
        // do nothing
    }
    return displayName;
}

const ROOT_SELECTOR = '[data-cy-root]';
const getContainerEl = () => {
    const el = document.querySelector(ROOT_SELECTOR);
    if (el) {
        return el;
    }
    throw Error(`No element found that matches selector ${ROOT_SELECTOR}. Please add a root element with data-cy-root attribute to your "component-index.html" file so that Cypress can attach your component to the DOM.`);
};
/**
 * Remove any style or extra link elements from the iframe placeholder
 * left from any previous test
 *
 */
function cleanupStyles() {
    const styles = document.body.querySelectorAll('[data-cy=injected-style-tag]');
    styles.forEach((styleElement) => {
        if (styleElement.parentElement) {
            styleElement.parentElement.removeChild(styleElement);
        }
    });
    const links = document.body.querySelectorAll('[data-cy=injected-stylesheet]');
    links.forEach((link) => {
        if (link.parentElement) {
            link.parentElement.removeChild(link);
        }
    });
}
/**
 * Insert links to external style resources.
 */
function insertStylesheets(stylesheets, document, el) {
    stylesheets.forEach((href) => {
        const link = document.createElement('link');
        link.type = 'text/css';
        link.rel = 'stylesheet';
        link.href = href;
        link.dataset.cy = 'injected-stylesheet';
        document.body.insertBefore(link, el);
    });
}
/**
 * Inserts a single stylesheet element
 */
function insertStyles(styles, document, el) {
    styles.forEach((style) => {
        const styleElement = document.createElement('style');
        styleElement.dataset.cy = 'injected-style-tag';
        styleElement.appendChild(document.createTextNode(style));
        document.body.insertBefore(styleElement, el);
    });
}
function insertSingleCssFile(cssFilename, document, el, log) {
    return cy.readFile(cssFilename, { log }).then((css) => {
        const style = document.createElement('style');
        style.appendChild(document.createTextNode(css));
        document.body.insertBefore(style, el);
    });
}
/**
 * Reads the given CSS file from local file system
 * and adds the loaded style text as an element.
 */
function insertLocalCssFiles(cssFilenames, document, el, log) {
    return Cypress.Promise.mapSeries(cssFilenames, (cssFilename) => {
        return insertSingleCssFile(cssFilename, document, el, log);
    });
}
/**
 * Injects custom style text or CSS file or 3rd party style resources
 * into the given document.
 */
const injectStylesBeforeElement = (options, document, el) => {
    if (!el)
        return;
    // first insert all stylesheets as Link elements
    let stylesheets = [];
    if (typeof options.stylesheet === 'string') {
        stylesheets.push(options.stylesheet);
    }
    else if (Array.isArray(options.stylesheet)) {
        stylesheets = stylesheets.concat(options.stylesheet);
    }
    if (typeof options.stylesheets === 'string') {
        options.stylesheets = [options.stylesheets];
    }
    if (options.stylesheets) {
        stylesheets = stylesheets.concat(options.stylesheets);
    }
    insertStylesheets(stylesheets, document, el);
    // insert any styles as <style>...</style> elements
    let styles = [];
    if (typeof options.style === 'string') {
        styles.push(options.style);
    }
    else if (Array.isArray(options.style)) {
        styles = styles.concat(options.style);
    }
    if (typeof options.styles === 'string') {
        styles.push(options.styles);
    }
    else if (Array.isArray(options.styles)) {
        styles = styles.concat(options.styles);
    }
    insertStyles(styles, document, el);
    // now load any css files by path and add their content
    // as <style>...</style> elements
    let cssFiles = [];
    if (typeof options.cssFile === 'string') {
        cssFiles.push(options.cssFile);
    }
    else if (Array.isArray(options.cssFile)) {
        cssFiles = cssFiles.concat(options.cssFile);
    }
    if (typeof options.cssFiles === 'string') {
        cssFiles.push(options.cssFiles);
    }
    else if (Array.isArray(options.cssFiles)) {
        cssFiles = cssFiles.concat(options.cssFiles);
    }
    return insertLocalCssFiles(cssFiles, document, el, options.log);
};
function setupHooks(optionalCallback) {
    // Consumed by the framework "mount" libs. A user might register their own mount in the scaffolded 'commands.js'
    // file that is imported by e2e and component support files by default. We don't want CT side effects to run when e2e
    // testing so we early return.
    // System test to verify CT side effects do not pollute e2e: system-tests/test/e2e_with_mount_import_spec.ts
    if (Cypress.testingType !== 'component') {
        return;
    }
    // When running component specs, we cannot allow "cy.visit"
    // because it will wipe out our preparation work, and does not make much sense
    // thus we overwrite "cy.visit" to throw an error
    Cypress.Commands.overwrite('visit', () => {
        throw new Error('cy.visit from a component spec is not allowed');
    });
    // @ts-ignore
    Cypress.on('test:before:run', () => {
        optionalCallback === null || optionalCallback === void 0 ? void 0 : optionalCallback();
        cleanupStyles();
    });
}

/**
 * Inject custom style text or CSS file or 3rd party style resources
 */
var injectStyles = function (options) {
    return function () {
        var el = getContainerEl();
        return injectStylesBeforeElement(options, document, el);
    };
};
var lastMountedReactDom;
/**
 * Create an `mount` function. Performs all the non-React-version specific
 * behavior related to mounting. The React-version-specific code
 * is injected. This helps us to maintain a consistent public API
 * and handle breaking changes in React's rendering API.
 *
 * This is designed to be consumed by `npm/react{16,17,18}`, and other React adapters,
 * or people writing adapters for third-party, custom adapters.
 */
var makeMountFn = function (type, jsx, options, rerenderKey, internalMountOptions) {
    if (options === void 0) { options = {}; }
    if (!internalMountOptions) {
        throw Error('internalMountOptions must be provided with `render` and `reactDom` parameters');
    }
    // Get the display name property via the component constructor
    // @ts-ignore FIXME
    var componentName = getDisplayName(jsx.type, options.alias);
    var displayName = options.alias || componentName;
    var jsxComponentName = "<" + componentName + " ... />";
    var message = options.alias
        ? jsxComponentName + " as \"" + options.alias + "\""
        : jsxComponentName;
    return cy
        .then(injectStyles(options))
        .then(function () {
        var _a, _b, _c;
        var reactDomToUse = internalMountOptions.reactDom;
        lastMountedReactDom = reactDomToUse;
        var el = getContainerEl();
        if (!el) {
            throw new Error([
                "[@cypress/react] \uD83D\uDD25 Hmm, cannot find root element to mount the component. Searched for " + ROOT_SELECTOR,
            ].join(' '));
        }
        var key = rerenderKey !== null && rerenderKey !== void 0 ? rerenderKey : 
        // @ts-ignore provide unique key to the the wrapped component to make sure we are rerendering between tests
        (((_c = (_b = (_a = Cypress === null || Cypress === void 0 ? void 0 : Cypress.mocha) === null || _a === void 0 ? void 0 : _a.getRunner()) === null || _b === void 0 ? void 0 : _b.test) === null || _c === void 0 ? void 0 : _c.title) || '') + Math.random();
        var props = {
            key: key,
        };
        var reactComponent = React.createElement(options.strict ? React.StrictMode : React.Fragment, props, jsx);
        // since we always surround the component with a fragment
        // let's get back the original component
        var userComponent = reactComponent.props.children;
        internalMountOptions.render(reactComponent, el, reactDomToUse);
        if (options.log !== false) {
            Cypress.log({
                name: type,
                type: 'parent',
                message: [message],
                // @ts-ignore
                $el: el.children.item(0),
                consoleProps: function () {
                    return {
                        // @ts-ignore protect the use of jsx functional components use ReactNode
                        props: jsx.props,
                        description: type === 'mount' ? 'Mounts React component' : 'Rerenders mounted React component',
                        home: 'https://github.com/cypress-io/cypress',
                    };
                },
            }).snapshot('mounted').end();
        }
        return (
        // Separate alias and returned value. Alias returns the component only, and the thenable returns the additional functions
        cy.wrap(userComponent, { log: false })
            .as(displayName)
            .then(function () {
            return cy.wrap({
                component: userComponent,
                rerender: function (newComponent) { return makeMountFn('rerender', newComponent, options, key, internalMountOptions); },
                unmount: internalMountOptions.unmount,
            }, { log: false });
        })
            // by waiting, we delaying test execution for the next tick of event loop
            // and letting hooks and component lifecycle methods to execute mount
            // https://github.com/bahmutov/cypress-react-unit-test/issues/200
            .wait(0, { log: false }));
        // Bluebird types are terrible. I don't think the return type can be carried without this cast
    });
};
/**
 * Create an `unmount` function. Performs all the non-React-version specific
 * behavior related to unmounting.
 *
 * This is designed to be consumed by `npm/react{16,17,18}`, and other React adapters,
 * or people writing adapters for third-party, custom adapters.
 */
var makeUnmountFn = function (options, internalUnmountOptions) {
    return cy.then(function () {
        return cy.get(ROOT_SELECTOR, { log: false }).then(function ($el) {
            var _a;
            if (lastMountedReactDom) {
                internalUnmountOptions.unmount($el[0]);
                var wasUnmounted = internalUnmountOptions.unmount($el[0]);
                if (wasUnmounted && options.log) {
                    Cypress.log({
                        name: 'unmount',
                        type: 'parent',
                        message: [(_a = options.boundComponentMessage) !== null && _a !== void 0 ? _a : 'Unmounted component'],
                        consoleProps: function () {
                            return {
                                description: 'Unmounts React component',
                                parent: $el[0],
                                home: 'https://github.com/cypress-io/cypress',
                            };
                        },
                    });
                }
            }
        });
    });
};
// Cleanup before each run
// NOTE: we cannot use unmount here because
// we are not in the context of a test
var preMountCleanup = function () {
    var el = getContainerEl();
    if (el && lastMountedReactDom) {
        lastMountedReactDom.unmountComponentAtNode(el);
    }
};
// Side effects from "import { mount } from '@cypress/<my-framework>'" are annoying, we should avoid doing this
// by creating an explicit function/import that the user can register in their 'component.js' support file,
// such as:
//    import 'cypress/<my-framework>/support'
// or
//    import { registerCT } from 'cypress/<my-framework>'
//    registerCT()
// Note: This would be a breaking change
// it is required to unmount component in beforeEach hook in order to provide a clean state inside test
// because `mount` can be called after some preparation that can side effect unmount
// @see npm/react/cypress/component/advanced/set-timeout-example/loading-indicator-spec.js
setupHooks(preMountCleanup);

var root;
function mount(jsx, options, rerenderKey) {
    if (options === void 0) { options = {}; }
    var internalOptions = {
        reactDom: ReactDOM,
        render: function (reactComponent, el) {
            root = ReactDOM.createRoot(el);
            return root.render(reactComponent);
        },
        unmount: unmount,
    };
    return makeMountFn('mount', jsx, __assign({ ReactDom: ReactDOM }, options), rerenderKey, internalOptions);
}
function unmount(options) {
    if (options === void 0) { options = { log: true }; }
    var internalOptions = {
        // type is ReturnType<typeof ReactDOM.createRoot>
        unmount: function () {
            root.unmount();
            return true;
        },
    };
    return makeUnmountFn(options, internalOptions);
}

export { mount, unmount };
