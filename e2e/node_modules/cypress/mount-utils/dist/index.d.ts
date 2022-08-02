/**
 * Additional styles to inject into the document.
 * A component might need 3rd party libraries from CDN,
 * local CSS files and custom styles.
 */
export interface StyleOptions {
    /**
     * Creates <link href="..." /> element for each stylesheet
     * @alias stylesheet
     */
    stylesheets: string | string[];
    /**
     * Creates <link href="..." /> element for each stylesheet
     * @alias stylesheets
     */
    stylesheet: string | string[];
    /**
     * Creates <style>...</style> element and inserts given CSS.
     * @alias styles
     */
    style: string | string[];
    /**
     * Creates <style>...</style> element for each given CSS text.
     * @alias style
     */
    styles: string | string[];
    /**
     * Loads each file and creates a <style>...</style> element
     * with the loaded CSS
     * @alias cssFile
     */
    cssFiles: string | string[];
    /**
     * Single CSS file to load into a <style></style> element
     * @alias cssFile
     */
    cssFile: string | string[];
}
export declare const ROOT_SELECTOR = "[data-cy-root]";
export declare const getContainerEl: () => HTMLElement;
/**
 * Remove any style or extra link elements from the iframe placeholder
 * left from any previous test
 *
 */
export declare function cleanupStyles(): void;
/**
 * Injects custom style text or CSS file or 3rd party style resources
 * into the given document.
 */
export declare const injectStylesBeforeElement: (options: Partial<StyleOptions & {
    log: boolean;
}>, document: Document, el: HTMLElement | null) => HTMLElement;
export declare function setupHooks(optionalCallback?: Function): void;
