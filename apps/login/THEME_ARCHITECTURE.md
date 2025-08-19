# Current Theme System Architecture

Our theme system provides a simple, environment variable-driven approach for consistent component styling and responsive layout switching.

## 🏗️ **Current Implementation**

### **Environment Variable Configuration**

```bash
# .env.local
NEXT_PUBLIC_THEME_ROUNDNESS=mid          # edgy | mid | full
NEXT_PUBLIC_THEME_LAYOUT=side-by-side    # side-by-side | top-to-bottom
NEXT_PUBLIC_THEME_APPEARANCE=material    # flat | material
NEXT_PUBLIC_THEME_SPACING=regular        # regular | compact
```

### **Core Theme Functions**

```tsx
// Server-safe theme configuration
import { getThemeConfig, getComponentRoundness } from "@/lib/theme";

// Get full theme configuration
const themeConfig = getThemeConfig();
// Returns: { roundness: 'mid', layout: 'side-by-side', appearance: 'material', ... }

// Get component-specific styling
const buttonRoundness = getComponentRoundness("button");
// Returns: "rounded-md" (CSS class)
```

### **Responsive Layout Hook**

```tsx
// Client-side responsive layout detection
import { useResponsiveLayout } from "@/lib/theme-hooks";

function MyComponent() {
  const { isSideBySide, isResponsiveOverride } = useResponsiveLayout();

  return <div className={isSideBySide ? "flex" : "flex-col"}>{/* Layout adapts automatically */}</div>;
}
```

## 🎨 **Component Integration Patterns**

### **Pattern 1: Direct Function Calls** (Current Standard)

```tsx
import { getComponentRoundness } from "@/lib/theme";

export function Button({ children, variant = "primary" }) {
  const roundness = getComponentRoundness("button");

  return (
    <button className={`${roundness} px-4 py-2 ${variant === "primary" ? "bg-blue-500" : "bg-gray-500"}`}>{children}</button>
  );
}
```

### **Pattern 2: Component-Specific Helper Functions**

```tsx
import { getComponentRoundness } from "@/lib/theme";

// Helper function for UserAvatar
function getUserAvatarRoundness(): string {
  return getComponentRoundness("avatarContainer");
}

export function UserAvatar({ loginName, displayName }) {
  const roundness = getUserAvatarRoundness();

  return <div className={`flex border p-1 ${roundness}`}>{/* Avatar content */}</div>;
}
```

### **Pattern 3: Theme-Aware Layout Components**

```tsx
import { useResponsiveLayout } from "@/lib/theme-hooks";

export function DynamicTheme({ children, branding }) {
  const { isSideBySide } = useResponsiveLayout();

  return (
    <ThemeWrapper branding={branding}>
      {isSideBySide ? (
        // Side-by-side layout for desktop
        <div className="flex max-w-[1200px]">
          <div className="w-1/2">{/* Left content */}</div>
          <div className="w-1/2">{/* Right content */}</div>
        </div>
      ) : (
        // Top-to-bottom layout for mobile
        <div className="flex-col max-w-[440px]">{children}</div>
      )}
    </ThemeWrapper>
  );
}
```

## 🎯 **Theme Configuration Structure**

### **Component Roundness Mapping**

```tsx
export interface ComponentRoundnessConfig {
  card: ThemeRoundness; // "rounded-lg" | "rounded-none" | "rounded-3xl"
  button: ThemeRoundness; // "rounded-md" | "rounded-none" | "rounded-full"
  input: ThemeRoundness; // "rounded-md" | "rounded-none" | "rounded-full pl-4"
  image: ThemeRoundness; // "rounded-lg" | "rounded-none" | "rounded-full"
  avatar: ThemeRoundness; // "rounded-lg" | "rounded-none" | "rounded-full"
  avatarContainer: ThemeRoundness; // "rounded-md" | "rounded-none" | "rounded-full"
  themeSwitch: ThemeRoundness; // "rounded-md" | "rounded-none" | "rounded-full"
}
```

### **Responsive Layout Logic**

```tsx
// Automatic layout switching based on screen size
const isSideBySide = themeConfig.layout === "side-by-side" && !isMdOrSmaller;

// md breakpoint: 768px (Tailwind default)
// Below 768px: Always use top-to-bottom layout
// Above 768px: Use configured layout (side-by-side or top-to-bottom)
```

## � **File Structure**

```
src/lib/
├── theme.ts           # Server-safe theme functions
├── theme-hooks.ts     # Client-side responsive hooks
└── themeUtils.tsx     # Legacy utility functions

src/components/
├── dynamic-theme.tsx  # Main responsive layout component
├── theme-wrapper.tsx  # Theme application wrapper
├── button.tsx         # Example themed component
├── card.tsx          # Example themed component
└── user-avatar.tsx   # Example themed component
```

## � **Usage Examples**

### **Adding Theme Support to New Components**

```tsx
import { getComponentRoundness } from "@/lib/theme";

export function NewComponent() {
  // Get theme-appropriate styling
  const roundness = getComponentRoundness("card");

  return <div className={`p-4 ${roundness} bg-white`}>{/* Component content */}</div>;
}
```

### **Using Responsive Layout**

```tsx
import { useResponsiveLayout } from "@/lib/theme-hooks";

export function ResponsiveComponent() {
  const { isSideBySide } = useResponsiveLayout();

  return <div className={isSideBySide ? "text-left" : "text-center"}>Content adapts to layout</div>;
}
```

### **Page Layout Integration**

```tsx
import { DynamicTheme } from "@/components/dynamic-theme";

export default function LoginPage() {
  return (
    <DynamicTheme branding={branding}>
      <div className="flex flex-col space-y-4">
        <h1>Login Title</h1>
        <p>Description text</p>
      </div>

      <div className="w-full">
        <LoginForm />
      </div>
    </DynamicTheme>
  );
}
```

## ⚡ **Key Features**

1. **Environment Variable Configuration**: Simple `.env.local` setup
2. **Server-Safe Functions**: Work in both SSR and client components
3. **Responsive Layout Switching**: Automatic mobile/desktop adaptation
4. **Component-Specific Styling**: Different roundness per component type
5. **Type Safety**: Full TypeScript support
6. **Zero Runtime Dependencies**: No context providers or complex state
7. **SSR Compatible**: No hydration mismatches

## 🔄 **Architecture Benefits**

- **Simple**: Environment variables → CSS classes
- **Fast**: No runtime theme calculations or context switching
- **Reliable**: Server-side rendering compatible
- **Scalable**: Easy to add new theme properties
- **Maintainable**: Clear separation between layout and styling concerns

This architecture provides a solid foundation for environment-driven theming while keeping the implementation simple and performant!
