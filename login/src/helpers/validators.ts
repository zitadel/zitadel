export function symbolValidator(value: string): boolean {
  const REGEXP = /[^a-zA-Z0-9]/gi;
  return REGEXP.test(value);
}

export function numberValidator(value: string): boolean {
  const REGEXP = /[0-9]/g;
  return REGEXP.test(value);
}

export function upperCaseValidator(value: string): boolean {
  const REGEXP = /[A-Z]/g;
  return REGEXP.test(value);
}

export function lowerCaseValidator(value: string): boolean {
  const REGEXP = /[a-z]/g;
  return REGEXP.test(value);
}
