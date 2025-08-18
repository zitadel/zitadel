# Future-Proof Theme System Architecture

## 🎯 **Problem Solved**

The new architecture allows multiple ways to inject theme configuration without modifying individual components:

## 🔧 **Three Flexible Approaches**

### 1. **Direct Props Injection** (Most Flexible)

```tsx
// Component accepts theme props directly
<Button roundness="rounded-full" variant="primary">
  Click me
</Button>

// Base component is theme-agnostic
<Button variant="primary"> // Uses defaults
  Click me
</Button>
```

### 2. **Higher-Order Component** (Automatic Injection)

```tsx
// HOC automatically injects theme from environment/context
<ThemedButton variant="primary">
  Click me // Gets roundness from NEXT_PUBLIC_THEME_ROUNDNESS
</ThemedButton>

// Can still override via props
<ThemedButton roundness="rounded-none" variant="primary">
  Override theme
</ThemedButton>
```

### 3. **Context Provider** (Global Theme Management)

```tsx
// App-level theme injection
<ThemeContextProvider customTheme={{ roundness: 'full', layout: 'top-to-bottom' }}>
  <YourApp />
</ThemeContextProvider>

// Or with API-driven themes
<ThemeContextProvider customTheme={apiThemeConfig}>
  <YourApp />
</ThemeContextProvider>
```

## 🚀 **Migration Paths**

### **Current State: Environment Variables**

````tsx
// No component changes needed
### **Current State: Environment Variables**
```tsx
// Create themed component using HOC
import { withTheme } from '@/lib/themeUtils';
import { Button } from '@/components/button';
const ThemedButton = withTheme(Button);

// No component changes needed
### **Current: Environment Variable HOC** (Recommended)
```tsx
import { withButtonTheme } from "@/lib/themeUtils";
import { Button } from "@/components/button";

const ThemedButton = withButtonTheme(Button);

<ThemedButton variant="primary">
  Gets theme from NEXT_PUBLIC_THEME_ROUNDNESS
</ThemedButton>
````

```

```

### **Future: API/Database Themes**

```tsx
// Option 1: Context injection
<ThemeContextProvider customTheme={userTheme}>
  <Button variant="primary">
    Gets theme from context
  </Button>
</ThemeContextProvider>

// Option 2: Direct injection
<Button roundness={userTheme.roundness} variant="primary">
  Gets theme from props
</Button>

// Option 3: HOC with custom logic
const ApiThemedButton = withApiTheme(Button);
<ApiThemedButton variant="primary">
  Gets theme from API
</ApiThemedButton>
```

### **Future: Runtime Theme Switching**

```tsx
function App() {
  const [currentTheme, setCurrentTheme] = useState(defaultTheme);

  return (
    <ThemeContextProvider customTheme={currentTheme}>
      <button onClick={() => setCurrentTheme(darkTheme)}>Switch to Dark</button>
      <Button variant="primary">Theme updates automatically</Button>
    </ThemeContextProvider>
  );
}
```

## ⚡ **Key Benefits**

1. **Zero Component Modifications**: Base components work with or without themes
2. **Multiple Injection Methods**: Props, HOC, Context, or API-driven
3. **Gradual Migration**: Can adopt new methods incrementally
4. **Type Safety**: All theme props are properly typed
5. **Performance**: No unnecessary re-renders or hooks in base components

## 🔄 **Component Architecture**

```
BaseComponent (theme-agnostic)
    ↓
withTheme(BaseComponent) → ThemedComponent (env-aware)
    ↓
withApiTheme(BaseComponent) → ApiThemedComponent (api-aware)
    ↓
withContext(BaseComponent) → ContextThemedComponent (context-aware)
```

Each layer adds functionality without modifying the base component!

## 📚 **Usage Examples**

```tsx
// Immediate: Use themed version with env vars
import { ThemedButton } from "@/components/ThemedButton";

// Future: Use base version with API props
import { Button } from "@/components/button";
<Button roundness={apiTheme.button.roundness} />;

// Future: Use context version
import { Button } from "@/components/button";
import { useTheme } from "@/lib/ThemeContext";

function MyComponent() {
  const { classes } = useTheme();
  return <Button roundness={classes.roundness.button} />;
}
```

This architecture scales from simple environment variables to complex multi-tenant theme systems without breaking changes!
