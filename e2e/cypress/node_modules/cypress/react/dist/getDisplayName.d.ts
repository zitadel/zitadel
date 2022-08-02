import { JSX } from './createMount';
/**
 * Gets the display name of the component when possible.
 * @param type {JSX} The type object returned from creating the react element.
 * @param fallbackName {string} The alias, or fallback name to use when the name cannot be derived.
 * @link https://github.com/facebook/react-devtools/blob/master/backend/getDisplayName.js
 */
export default function getDisplayName(type: JSX, fallbackName?: string): string;
