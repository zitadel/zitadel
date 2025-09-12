/**
 * Utility functions for working with TypeScript enums in Angular forms
 *
 * Example usage:
 * ```typescript
 * // Get dropdown options
 * const bindingOptions = getEnumKeys(SAMLBinding);
 *
 * // Convert backend enum values to form-friendly string keys
 * const formConfig = convertEnumValuesToKeys(backendConfig, {
 *   binding: SAMLBinding,
 *   nameIdFormat: SAMLNameIDFormat
 * });
 *
 * // Find specific enum key
 * const defaultBinding = getEnumKeyFromValue(SAMLBinding, SAMLBinding.SAML_BINDING_POST);
 * ```
 */

// Type constraint for TypeScript numeric enums
type NumericEnum = Record<string, string | number> & Record<number, string>;

/**
 * Get string keys from a TypeScript enum, excluding numeric reverse mappings
 * @param enumObject The enum object
 * @returns Array of string keys
 * @example
 * ```typescript
 * const bindingOptions = getEnumKeys(SAMLBinding);
 * // Returns: ['SAML_BINDING_UNSPECIFIED', 'SAML_BINDING_POST', ...]
 * ```
 */
export function getEnumKeys<T extends NumericEnum>(enumObject: T): string[] {
  return Object.keys(enumObject).filter((key) => isNaN(Number(key)));
}

/**
 * Find the string key for a given numeric enum value
 * @param enumObject The enum object
 * @param value The numeric enum value
 * @returns The corresponding string key or undefined if not found
 * @example
 * ```typescript
 * const key = getEnumKeyFromValue(SAMLBinding, 1);
 * // Returns: 'SAML_BINDING_POST'
 * ```
 */
export function getEnumKeyFromValue<T extends NumericEnum>(enumObject: T, value: number): string | undefined {
  return Object.keys(enumObject).find((key) => enumObject[key] === value && isNaN(Number(key)));
}

/**
 * Convert enum values to string keys for form controls
 * @param config Object containing enum values to convert
 * @param enumMappings Object mapping config property names to enum objects
 * @returns Modified config object with string keys
 * @example
 * ```typescript
 * const formConfig = convertEnumValuesToKeys(samlConfig, {
 *   binding: SAMLBinding,
 *   nameIdFormat: SAMLNameIDFormat
 * });
 * // Converts: { binding: 1, nameIdFormat: 2 }
 * // To: { binding: 'SAML_BINDING_POST', nameIdFormat: 'SAML_NAME_ID_FORMAT_PERSISTENT' }
 * ```
 */
export function convertEnumValuesToKeys<T extends Record<string, unknown>>(
  config: T,
  enumMappings: { [K in keyof Partial<T>]: NumericEnum },
): T {
  const converted = { ...config };

  for (const [propertyName, enumObject] of Object.entries(enumMappings)) {
    const typedPropertyName = propertyName as keyof T;
    if (converted[typedPropertyName] !== undefined && typeof converted[typedPropertyName] === 'number') {
      const key = getEnumKeyFromValue(enumObject, converted[typedPropertyName] as number);
      if (key) {
        (converted as Record<string, unknown>)[propertyName] = key;
      }
    }
  }

  return converted;
}
