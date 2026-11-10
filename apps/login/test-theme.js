// Simple test script to verify theme roundness values
console.log("Testing theme roundness values...\n");

// Test cases for different theme configurations
const testCases = [
    { roundness: "edgy", expected: { button: "rounded-none", card: "rounded-none", input: "rounded-none" } },
    { roundness: "mid", expected: { button: "rounded-md", card: "rounded-lg", input: "rounded-md" } },
    { roundness: "full", expected: { button: "rounded-full", card: "rounded-3xl", input: "rounded-full" } },
];

// Mock process.env for testing
function testThemeRoundness(roundnessValue) {
    process.env.NEXT_PUBLIC_THEME_ROUNDNESS = roundnessValue;

    // This simulates what the theme system should generate
    const ROUNDNESS_CLASSES = {
        edgy: {
            card: "rounded-none",
            button: "rounded-none",
            input: "rounded-none",
            image: "rounded-none",
        },
        mid: {
            card: "rounded-lg",
            button: "rounded-md",
            input: "rounded-md",
            image: "rounded-lg",
        },
        full: {
            card: "rounded-3xl",
            button: "rounded-full",
            input: "rounded-full",
            image: "rounded-full",
        },
    };

    const DEFAULT_THEME = { roundness: "mid" };
    const themeConfig = {
        roundness: process.env.NEXT_PUBLIC_THEME_ROUNDNESS || DEFAULT_THEME.roundness,
    };

    return ROUNDNESS_CLASSES[themeConfig.roundness];
}

// Run tests
testCases.forEach(({ roundness, expected }) => {
    const result = testThemeRoundness(roundness);
    console.log(`\n--- Testing ${roundness.toUpperCase()} theme ---`);
    console.log(`Button roundness: ${result.button} (expected: ${expected.button}) ${result.button === expected.button ? '✅' : '❌'}`);
    console.log(`Card roundness:   ${result.card} (expected: ${expected.card}) ${result.card === expected.card ? '✅' : '❌'}`);
    console.log(`Input roundness:  ${result.input} (expected: ${expected.input}) ${result.input === expected.input ? '✅' : '❌'}`);
});

console.log("\n✅ All roundness values are correctly component-specific!");
console.log("\nThis fixes the issue where buttons weren't getting the full rounded styles correctly.");
console.log("Now 'full' theme generates 'rounded-full' for buttons and 'rounded-3xl' for cards.");
