/**
 * Run async work over items with a fixed number of in-flight tasks.
 * Avoids unbounded Promise.all which can exhaust ephemeral ports
 * when seeding large USER_AMOUNT datasets.
 */
export async function mapPool<T, R>(
  items: T[],
  concurrency: number, 
  fn: (item: T, index: number) => Promise<R>,
): Promise<R[]> {
  if (items.length === 0) {
    return [];
  }

  const limit = Math.max(1, Math.min(concurrency, items.length));
  const results: R[] = new Array(items.length);
  let nextIndex = 0;

  async function worker(): Promise<void> {
    while (true) {
      const index = nextIndex;
      nextIndex += 1;
      if (index >= items.length) {
        return;
      }
      results[index] = await fn(items[index], index);
    }
  }

  await Promise.all(Array.from({ length: limit }, () => worker()));
  return results;
}
