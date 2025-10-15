/**
 * Utility functions for generating random names for projects and applications
 */

const ADJECTIVES = [
  'awesome',
  'brilliant',
  'creative',
  'dynamic',
  'elegant',
  'fantastic',
  'innovative',
  'amazing',
  'stellar',
  'epic',
  'modern',
  'smart',
  'sleek',
  'powerful',
  'advanced',
  'cool',
  'fresh',
  'bright',
  'fun',
  'agile',
];

const PROJECT_NOUNS = [
  'project',
  'workspace',
  'platform',
  'solution',
  'system',
  'portal',
  'hub',
  'dashboard',
  'service',
  'application',
  'site',
  'app',
];

const APP_NOUNS = [
  'app',
  'application',
  'client',
  'webapp',
  'service',
  'portal',
  'dashboard',
  'interface',
  'frontend',
  'backend',
];

/**
 * Generates a random project name in the format: adjective-noun-number
 * @returns A random project name (e.g., "awesome-workspace-427")
 */
export function generateRandomProjectName(): string {
  const adjective = ADJECTIVES[Math.floor(Math.random() * ADJECTIVES.length)];
  const noun = PROJECT_NOUNS[Math.floor(Math.random() * PROJECT_NOUNS.length)];
  const number = Math.floor(Math.random() * 1000);

  return `${adjective}-${noun}-${number}`;
}

/**
 * Generates a random app name in the format: adjective-noun-number
 * @returns A random app name (e.g., "brilliant-app-891")
 */
export function generateRandomAppName(): string {
  const adjective = ADJECTIVES[Math.floor(Math.random() * ADJECTIVES.length)];
  const noun = APP_NOUNS[Math.floor(Math.random() * APP_NOUNS.length)];
  const number = Math.floor(Math.random() * 1000);

  return `${adjective}-${noun}-${number}`;
}

/**
 * Generates a framework-specific app name
 * @param frameworkTitle The title of the framework (e.g., "React", "Next.js")
 * @returns A framework-specific app name (e.g., "My React App")
 */
export function generateFrameworkAppName(frameworkTitle?: string): string {
  if (frameworkTitle) {
    return `My ${frameworkTitle} App`;
  }

  // Fallback to random name if no framework title provided
  const adjective = ADJECTIVES[Math.floor(Math.random() * ADJECTIVES.length)];
  return `My ${adjective} App`;
}

/**
 * Generates a random name with custom format
 * @param adjectives Array of adjectives to choose from
 * @param nouns Array of nouns to choose from
 * @param separator Separator character (default: '-')
 * @param includeNumber Whether to include a random number (default: true)
 * @returns A random name with the specified format
 */
export function generateCustomName(
  adjectives: string[],
  nouns: string[],
  separator: string = '-',
  includeNumber: boolean = true,
): string {
  const adjective = adjectives[Math.floor(Math.random() * adjectives.length)];
  const noun = nouns[Math.floor(Math.random() * nouns.length)];

  if (includeNumber) {
    const number = Math.floor(Math.random() * 1000);
    return `${adjective}${separator}${noun}${separator}${number}`;
  }

  return `${adjective}${separator}${noun}`;
}
