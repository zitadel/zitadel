import * as t from 'tcomb';

export * from 'tcomb';

// Augment 'tcomb': Add getValidationErrorMessage.
declare module 'tcomb' {
    export interface Type<T> {
        /**
         * Allows customization of the error message for a type.
         * 
         * (Extension from tcomb-validation)
         * 
         * @param actual Current value
         * @param path Path to validate
         * @param context Additional metadata.
         */
        getValidationErrorMessage(actual: T, path: Path, context: any): string;
    }
}

/**
 * Defines a path through the properties of an
 * object (string, property name) or array (number, index).
 */
type Path = Array<string | number>;
type Predicate<T> = (value: T) => boolean;

export interface ValidationError {
    /**
     * Error message.
     */
    message: string;
    /**
     * Current value.
     */
    actual: any;
    /**
     * Expected type.
     */
    expected: t.Type<any>;
    /**
     * Path to the property/index that failed validation.
     */
    path: Path;
}

/**
 * Result of a validation.
 */
export interface ValidationResult {
    /**
     * True if there are no validation errors. False otherwise.
     */
    isValid(): boolean;
    /**
     * Returns the first error, if any. Null otherwise.
     */
    firstError(): ValidationError | null;
    /**
     * Contains the validation errors, if any. Empty if none.
     */
    errors: Array<ValidationError>;
}

/**
 * Options for the validate function.
 */
interface ValidateOptions {
    /**
     * Path prefix for validation.
     */
    path?: Path;
    /**
     * Data passed to getValidationErrorMessage.
     */
    context?: any;
    /**
     * If true, no additional properties are allowed
     * when validating structs.
     * 
     * Defaults to false.
     */
    strict?: boolean;
}

/**
 * Validates an object and returns the validation result.
 * @param value The value to validate.
 * @param type The type to validate against.
 * @param options Validation options. Optional.
 */
export function validate<T>(value: any, type: t.Type<T>, options?: ValidateOptions): ValidationResult;

