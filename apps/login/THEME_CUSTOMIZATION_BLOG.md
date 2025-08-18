# Building a Comprehensive Theme System for Modern Login Applications

_How we transformed a static login interface into a fully customizable, brand-aligned experience_

---

## Introduction

In today's digital landscape, user experience is paramount. Organizations need their authentication flows to not only be secure and reliable but also reflect their unique brand identity. At ZITADEL, we recognized that a one-size-fits-all approach to login interfaces simply wasn't enough. This led us to develop a comprehensive theme customization system that allows organizations to create login experiences that truly represent their brand.

## The Challenge

Before our theme system, customizing the appearance of login interfaces was a complex, developer-intensive process. Organizations had to choose between:

- **Generic interfaces** that worked but didn't match their brand
- **Complex customizations** that required deep technical knowledge and ongoing maintenance
- **Fragmented styling** where different components used inconsistent design patterns

We knew there had to be a better way.

## Our Solution: A Component-First Theme Architecture

We built a sophisticated yet intuitive theme system that operates on four core principles:

### 1. **Environment Variable Configuration**

Simple, declarative configuration through environment variables:

```bash
NEXT_PUBLIC_THEME_ROUNDNESS=mid
NEXT_PUBLIC_THEME_LAYOUT=top-to-bottom
NEXT_PUBLIC_THEME_APPEARANCE=material
NEXT_PUBLIC_THEME_SPACING=regular
```

### 2. **Component-Specific Defaults**

Different components can have their own styling defaults while maintaining overall consistency:

```typescript
export const DEFAULT_COMPONENT_ROUNDNESS: ComponentRoundnessConfig = {
  card: "mid", // Cards use moderate roundness
  button: "mid", // Buttons use moderate roundness
  input: "mid", // Inputs use moderate roundness
  avatar: "full", // Avatars default to full roundness
  avatarContainer: "full", // Avatar containers use full roundness
  themeSwitch: "full", // Theme toggle uses full roundness
};
```

### 3. **Smart Override System**

When a global theme is set, it intelligently overrides component-specific defaults, ensuring consistency when needed:

```typescript
// Without global setting: uses component-specific defaults
// With NEXT_PUBLIC_THEME_ROUNDNESS=full: everything becomes fully rounded
```

### 4. **Future-Proof Architecture**

Built with extensibility in mind, avoiding utility function patterns that could cause conflicts with CSS frameworks like Tailwind CSS.

## Key Features and Benefits

### **Comprehensive Visual Control**

Our theme system provides control over four major visual aspects:

**Roundness Options:**

- **Edgy** (`rounded-none`): Sharp, modern industrial aesthetic
- **Mid** (`rounded-md`/`rounded-lg`): Balanced professional appearance
- **Full** (`rounded-full`): Friendly, approachable design

**Layout Flexibility:**

- **Side-by-side**: Desktop-optimized with brand showcase
- **Top-to-bottom**: Mobile-first, vertical flow

**Appearance Philosophies:**

- **Flat**: Clean, borderless design with subtle shadows
- **Material**: Rich shadows and depth for enhanced interactivity

**Spacing Control:**

- **Regular**: Comfortable, spacious layout
- **Compact**: Dense, efficient use of space

### **Unified Component Integration**

Every interactive element respects the theme system:

✅ **Form Elements**: Buttons, inputs, cards  
✅ **Navigation**: Theme toggles, language switchers  
✅ **Authentication**: IDP buttons (Google, Apple, GitHub, etc.)  
✅ **User Interface**: Avatars, containers, dropdowns

### **Developer Experience Excellence**

**Simple Integration:**

```typescript
// Any component can easily adopt theming
function getComponentRoundness(componentType: keyof ComponentRoundnessConfig): string {
  const themeConfig = getThemeConfig();
  const roundnessLevel = themeConfig.componentRoundness?.[componentType] || themeConfig.roundness;
  return ROUNDNESS_CLASSES[roundnessLevel][componentType];
}
```

**Type Safety:**

```typescript
interface ComponentRoundnessConfig {
  card: ThemeRoundness;
  button: ThemeRoundness;
  input: ThemeRoundness;
  // ... fully typed configuration
}
```

**Consistent Patterns:**
Every component follows the same integration pattern, making the system predictable and maintainable.

## Real-World Impact

### **For Organizations**

- **Brand Consistency**: Login interfaces that perfectly match corporate identity
- **User Experience**: Cohesive design language throughout the authentication flow
- **Flexibility**: Easy switching between design approaches without developer intervention

### **For Developers**

- **Maintainability**: Centralized theme configuration reduces code duplication
- **Extensibility**: Adding new components or theme options is straightforward
- **Performance**: Direct class application avoids utility function overhead

### **For End Users**

- **Familiarity**: Interfaces that feel native to the organization's ecosystem
- **Accessibility**: Consistent interaction patterns across all elements
- **Responsiveness**: Themes that work seamlessly across devices and contexts

## Technical Highlights

### **Smart Fallback Logic**

```typescript
// Uses component-specific defaults when no global theme is set
// Automatically applies global theme when environment variable is provided
const componentRoundness = globalRoundness
  ? {
      card: globalRoundness,
      button: globalRoundness,
      avatar: globalRoundness,
      // ... all components use global setting
    }
  : DEFAULT_COMPONENT_ROUNDNESS; // Use component-specific defaults
```

### **CSS Framework Compatibility**

By using direct class application instead of utility functions, our system avoids conflicts with CSS frameworks and ensures reliable styling.

### **Component Encapsulation**

Each component can be themed independently while maintaining overall design coherence:

```typescript
// IDP buttons automatically inherit theme settings
export const BaseButton = forwardRef<HTMLButtonElement, SignInWithIdentityProviderProps>(function BaseButton(props, ref) {
  const buttonRoundness = getComponentRoundness("button");
  // ... theming applied consistently across all IDP providers
});
```

## Looking to the Future

Our theme system is designed with extensibility in mind. Here's what we're excited to explore next:

### **Custom CSS Properties Integration**

```css
:root {
  --zitadel-primary-color: #4f46e5;
  --zitadel-secondary-color: #e5e7eb;
  --zitadel-font-family: "Inter", sans-serif;
  --zitadel-border-radius: 8px;
}
```

Imagine being able to define custom CSS properties that automatically propagate throughout the entire interface, providing even more granular control over brand expression.

### **Dedicated Theme API**

```typescript
// Future: Runtime theme management
const themeAPI = useZitadelTheme();

await themeAPI.updateTheme({
  roundness: "full",
  appearance: "material",
  customColors: {
    primary: "#your-brand-color",
    secondary: "#your-accent-color",
  },
});
```

A dedicated API would enable:

- **Runtime theme switching** without environment variable changes
- **User preference persistence** across sessions
- **A/B testing** of different design approaches
- **Dynamic branding** based on organization context

### **Advanced Component Theming**

```typescript
// Future: Per-component theme overrides
interface AdvancedThemeConfig {
  global: ThemeSettings;
  components: {
    loginForm: Partial<ThemeSettings>;
    idpButtons: Partial<ThemeSettings>;
    navigation: Partial<ThemeSettings>;
  };
}
```

### **Design Token Integration**

```json
{
  "color": {
    "primary": {
      "50": "#eff6ff",
      "500": "#3b82f6",
      "900": "#1e3a8a"
    }
  },
  "spacing": {
    "xs": "0.5rem",
    "sm": "1rem",
    "md": "1.5rem"
  }
}
```

Integration with design token systems would enable seamless designer-developer collaboration and ensure perfect brand consistency.

### **Component-Based Spacing System**

Just as we implemented component-specific roundness, spacing could follow the same pattern:

```typescript
interface ComponentSpacingConfig {
  card: ThemeSpacing;
  button: ThemeSpacing;
  form: ThemeSpacing;
  navigation: ThemeSpacing;
  avatar: ThemeSpacing;
}

export const DEFAULT_COMPONENT_SPACING: ComponentSpacingConfig = {
  card: "regular", // Cards use standard spacing
  button: "compact", // Buttons use tighter spacing
  form: "regular", // Forms use comfortable spacing
  navigation: "compact", // Navigation uses efficient spacing
  avatar: "compact", // Avatars use minimal spacing
};
```

This would enable scenarios like:

- **Dense navigation bars** with compact spacing while maintaining **comfortable form layouts**
- **Tight button groups** alongside **spacious card layouts**
- **Context-aware spacing** that adapts to component purpose

### **CSS-Based Theme Overrides**

The future of component theming lies in CSS custom properties and data attributes:

#### **CSS Custom Properties Approach**

```css
/* Global theme variables */
:root {
  --zitadel-button-bg: #3b82f6;
  --zitadel-button-text: #ffffff;
  --zitadel-button-radius: 0.5rem;
  --zitadel-button-spacing: 0.75rem 1.5rem;
}

/* Component-specific overrides */
.idp-button {
  --zitadel-button-bg: #1f2937;
  --zitadel-button-radius: 0.375rem;
}

.primary-button {
  --zitadel-button-bg: var(--brand-primary, #3b82f6);
  --zitadel-button-text: var(--brand-primary-contrast, #ffffff);
}
```

#### **Data Attribute System**

```typescript
// Component renders with theme data attributes
<button
  data-theme-component="button"
  data-theme-variant="primary"
  data-theme-size="medium"
  className={getButtonClasses()}
>
  Sign In
</button>
```

```css
/* CSS can target specific theme combinations */
[data-theme-component="button"][data-theme-variant="primary"] {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

[data-theme-component="button"][data-theme-variant="secondary"] {
  background: transparent;
  border: 2px solid currentColor;
  backdrop-filter: blur(10px);
}
```

#### **CSS-in-JS Integration**

```typescript
// Future: Dynamic CSS injection based on theme
const generateThemeCSS = (theme: ThemeConfig) => `
  .zitadel-button {
    background: ${theme.colors.primary};
    border-radius: ${theme.roundness === "full" ? "9999px" : "0.5rem"};
    padding: ${theme.spacing === "compact" ? "0.5rem 1rem" : "0.75rem 1.5rem"};
    transition: all 0.2s ease;
  }
  
  .zitadel-button:hover {
    background: ${theme.colors.primaryHover};
    transform: ${theme.appearance === "material" ? "translateY(-1px)" : "none"};
    box-shadow: ${theme.appearance === "material" ? "0 4px 12px rgba(0,0,0,0.15)" : "none"};
  }
`;
```

### **Advanced Customization Patterns**

#### **Conditional Component Theming**

```typescript
interface ConditionalTheme {
  condition: (context: ThemeContext) => boolean;
  theme: Partial<ThemeConfig>;
}

const contextualThemes: ConditionalTheme[] = [
  {
    condition: (ctx) => ctx.userRole === "admin",
    theme: { appearance: "material", roundness: "edgy" },
  },
  {
    condition: (ctx) => ctx.deviceType === "mobile",
    theme: { spacing: "compact", layout: "top-to-bottom" },
  },
  {
    condition: (ctx) => ctx.timeOfDay === "night",
    theme: { appearance: "dark", roundness: "full" },
  },
];
```

#### **Component Composition Theming**

```typescript
// Theme inheritance and composition
interface CompositeTheme {
  base: ThemeConfig;
  overrides: {
    [componentPath: string]: Partial<ThemeConfig>;
  };
}

const enterpriseTheme: CompositeTheme = {
  base: { roundness: "mid", appearance: "material" },
  overrides: {
    "auth.loginForm": { roundness: "edgy", spacing: "compact" },
    "auth.idpButtons": { appearance: "flat", roundness: "full" },
    "navigation.*": { spacing: "compact" },
    "*.primaryButton": { appearance: "material" },
  },
};
```

### **CSS-First Theme Architecture**

Imagine a future where themes are primarily defined in CSS, with JavaScript acting as the orchestrator:

```css
/* themes/corporate.css */
@layer zitadel-theme {
  :root {
    --zitadel-theme-name: "corporate";
    --zitadel-primary: #1e3a8a;
    --zitadel-secondary: #e5e7eb;
    --zitadel-radius-sm: 0.25rem;
    --zitadel-radius-md: 0.375rem;
    --zitadel-radius-lg: 0.5rem;
  }

  .zitadel-component {
    border-radius: var(--zitadel-radius-md);
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .zitadel-button {
    background: var(--zitadel-primary);
    color: white;
    padding: 0.75rem 1.5rem;
    border-radius: var(--zitadel-radius-lg);

    &:hover {
      background: color-mix(in srgb, var(--zitadel-primary) 90%, black);
      transform: translateY(-1px);
    }

    &[data-variant="secondary"] {
      background: transparent;
      color: var(--zitadel-primary);
      border: 2px solid var(--zitadel-primary);
    }
  }

  .zitadel-card {
    border-radius: var(--zitadel-radius-lg);
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
    background: white;

    @media (prefers-color-scheme: dark) {
      background: #1f2937;
      box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.3);
    }
  }
}
```

```typescript
// JavaScript orchestrates CSS theme loading
class CSSThemeManager {
  async loadTheme(themeName: string) {
    const themeCSS = await import(`./themes/${themeName}.css`);
    this.injectThemeCSS(themeCSS);
    this.updateComponentDataAttributes(themeName);
  }

  private updateComponentDataAttributes(themeName: string) {
    document.documentElement.dataset.zitadelTheme = themeName;
    // Components automatically pick up new styling
  }
}
```

This approach would provide:

- **Designer-friendly** theme creation
- **Performance benefits** through CSS optimization
- **Dynamic theme switching** without JavaScript recalculation
- **Better caching** of theme assets
- **CSS-native features** like `color-mix()`, `@container`, `@layer`

### **Beyond Spacing: Other Component-Based Customizations**

The component-specific approach could extend to many other visual properties:

#### **Component-Based Typography**

```typescript
interface ComponentTypographyConfig {
  heading: "large" | "medium" | "small";
  body: "readable" | "compact" | "dense";
  button: "bold" | "medium" | "light";
  navigation: "uppercase" | "normal" | "small-caps";
}

// Different components could have different typography styles
export const DEFAULT_COMPONENT_TYPOGRAPHY: ComponentTypographyConfig = {
  heading: "large", // Headings are prominent
  body: "readable", // Body text is optimized for reading
  button: "medium", // Buttons use moderate weight
  navigation: "uppercase", // Navigation uses all caps for structure
};
```

#### **Component-Based Animation**

```typescript
interface ComponentAnimationConfig {
  buttons: "none" | "subtle" | "playful";
  forms: "none" | "smooth" | "bouncy";
  navigation: "instant" | "slide" | "fade";
  feedback: "minimal" | "standard" | "expressive";
}

// Example: Enterprise might prefer minimal animations, consumer apps might be more playful
const enterpriseAnimations: ComponentAnimationConfig = {
  buttons: "subtle", // Professional, minimal hover effects
  forms: "smooth", // Clean form transitions
  navigation: "instant", // Fast, efficient navigation
  feedback: "minimal", // Subtle success/error states
};

const consumerAnimations: ComponentAnimationConfig = {
  buttons: "playful", // Fun hover and click effects
  forms: "bouncy", // Engaging form interactions
  navigation: "slide", // Smooth page transitions
  feedback: "expressive", // Clear, animated feedback
};
```

#### **Component-Based Color Systems**

```typescript
interface ComponentColorConfig {
  primary: ColorScheme;
  secondary: ColorScheme;
  accent: ColorScheme;
  semantic: SemanticColors;
}

interface ColorScheme {
  base: string;
  hover: string;
  active: string;
  disabled: string;
}

// Different components could use different color schemes
const componentColors = {
  buttons: "primary",
  links: "accent",
  navigation: "secondary",
  alerts: "semantic",
  forms: "secondary",
};
```

#### **Component-Based Interaction Patterns**

```typescript
interface ComponentInteractionConfig {
  buttons: "click" | "hover" | "tap";
  navigation: "click" | "hover-preview" | "gesture";
  forms: "focus" | "hover" | "touch-optimized";
  tooltips: "hover" | "click" | "disabled";
}

// Mobile-first might emphasize touch interactions
const mobileInteractions: ComponentInteractionConfig = {
  buttons: "tap", // Optimized for touch
  navigation: "gesture", // Swipe navigation
  forms: "touch-optimized", // Large touch targets
  tooltips: "click", // Click to show tooltips
};
```

### **CSS Override Strategies for Enterprise Customers**

Since we currently use Tailwind CSS classes, enterprise customers who need complete styling control will require robust CSS override mechanisms. Here are several approaches we could implement:

#### **CSS Custom Properties with Tailwind Integration**

```css
/* Customer's custom.css - loaded after our styles */
:root {
  /* Override our CSS custom properties */
  --zitadel-primary: #your-brand-color;
  --zitadel-primary-hover: #your-brand-hover;
  --zitadel-border-radius: 12px;
}

/* Direct Tailwind utility overrides */
.zitadel-button {
  @apply bg-[var(--zitadel-primary)] hover:bg-[var(--zitadel-primary-hover)];
  border-radius: var(--zitadel-border-radius) !important;
}

/* Component-specific overrides */
[data-zitadel-component="login-form"] .zitadel-button {
  background: linear-gradient(45deg, #ff6b6b, #4ecdc4) !important;
  border: none !important;
  transform: perspective(1px) translateZ(0) !important;
}
```

#### **CSS Injection API**

```typescript
// Future API for runtime CSS injection
interface ZitadelCSSOverride {
  selector: string;
  styles: CSSStyleDeclaration | string;
  priority?: "low" | "normal" | "high" | "important";
}

class ZitadelCustomizer {
  // Inject custom CSS at runtime
  addCustomCSS(overrides: ZitadelCSSOverride[]) {
    const styleSheet = this.createCustomStyleSheet();

    overrides.forEach(({ selector, styles, priority = "normal" }) => {
      const rule = this.generateCSSRule(selector, styles, priority);
      styleSheet.insertRule(rule);
    });
  }

  // Override specific components
  overrideComponent(componentName: string, styles: Partial<CSSStyleDeclaration>) {
    this.addCustomCSS([
      {
        selector: `[data-zitadel-component="${componentName}"]`,
        styles,
        priority: "high",
      },
    ]);
  }
}

// Usage
const customizer = new ZitadelCustomizer();

customizer.overrideComponent("button", {
  background: "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
  border: "none",
  borderRadius: "8px",
  boxShadow: "0 4px 15px rgba(102, 126, 234, 0.4)",
});
```

#### **Tailwind CSS Override Patterns**

```css
/* Method 1: CSS Layers (Modern approach) */
@layer customer-overrides {
  /* These styles will have higher specificity than our base styles */
  .zitadel-button {
    @apply bg-gradient-to-r from-purple-500 to-blue-500 hover:from-purple-600 hover:to-blue-600;
    @apply shadow-lg hover:shadow-xl transform hover:-translate-y-0.5;
    @apply transition-all duration-200 ease-in-out;
  }

  .zitadel-input {
    @apply border-2 border-gray-300 focus:border-blue-500 focus:ring-2 focus:ring-blue-200;
    @apply bg-white dark:bg-gray-800 dark:border-gray-600;
  }
}

/* Method 2: Higher Specificity Selectors */
.zitadel-login-container .zitadel-button.zitadel-primary {
  /* More specific selector automatically overrides base styles */
  background: #your-custom-color !important;
  border-radius: 12px !important;
}

/* Method 3: Component-Scoped Overrides */
[data-zitadel-theme="corporate"] {
  --tw-bg-primary: #1e3a8a;
  --tw-text-primary: #ffffff;
  --tw-rounded-default: 0.375rem;
}

[data-zitadel-theme="corporate"] .zitadel-button {
  @apply bg-[var(--tw-bg-primary)] text-[var(--tw-text-primary)];
  border-radius: var(--tw-rounded-default);
}
```

#### **Dynamic Tailwind Class Override System**

```typescript
// Future: Runtime Tailwind class replacement
interface TailwindOverride {
  component: string;
  originalClasses: string[];
  overrideClasses: string[];
}

class TailwindOverrideManager {
  private overrides: Map<string, TailwindOverride> = new Map();

  addOverride(override: TailwindOverride) {
    this.overrides.set(override.component, override);
    this.applyOverrides();
  }

  private applyOverrides() {
    this.overrides.forEach((override, componentName) => {
      const elements = document.querySelectorAll(`[data-zitadel-component="${componentName}"]`);

      elements.forEach((element) => {
        // Remove original Tailwind classes
        element.classList.remove(...override.originalClasses);
        // Add custom classes
        element.classList.add(...override.overrideClasses);
      });
    });
  }
}

// Usage
const overrideManager = new TailwindOverrideManager();

overrideManager.addOverride({
  component: "button",
  originalClasses: ["bg-blue-500", "hover:bg-blue-600", "rounded-md"],
  overrideClasses: [
    "bg-gradient-to-r",
    "from-purple-500",
    "to-pink-500",
    "hover:from-purple-600",
    "hover:to-pink-600",
    "rounded-lg",
  ],
});
```

#### **CSS Module Override System**

```typescript
// Future: CSS-in-JS with override support
interface ComponentStyleOverrides {
  [componentName: string]: {
    base?: string;
    variants?: {
      [variantName: string]: string;
    };
    states?: {
      hover?: string;
      focus?: string;
      active?: string;
      disabled?: string;
    };
  };
}

const customerOverrides: ComponentStyleOverrides = {
  button: {
    base: `
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      border: none;
      border-radius: 8px;
      color: white;
      font-weight: 600;
      padding: 12px 24px;
      transition: all 0.2s ease;
    `,
    variants: {
      secondary: `
        background: transparent;
        border: 2px solid #667eea;
        color: #667eea;
      `,
    },
    states: {
      hover: `
        transform: translateY(-2px);
        box-shadow: 0 8px 25px rgba(102, 126, 234, 0.3);
      `,
      active: `
        transform: translateY(0);
        box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
      `,
    },
  },
  input: {
    base: `
      border: 2px solid #e2e8f0;
      border-radius: 8px;
      padding: 12px 16px;
      background: white;
      transition: all 0.2s ease;
    `,
    states: {
      focus: `
        border-color: #667eea;
        box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
        outline: none;
      `,
    },
  },
};
```

#### **CSS Variable Bridge System**

```typescript
// Bridge between our theme system and customer CSS
interface CSSVariableBridge {
  generateCustomProperties(theme: ThemeConfig): string;
  injectCustomerOverrides(customerCSS: string): void;
}

class ZitadelCSSBridge implements CSSVariableBridge {
  generateCustomProperties(theme: ThemeConfig): string {
    return `
      :root {
        /* Expose our theme values as CSS custom properties */
        --zitadel-roundness: ${this.getRoundnessValue(theme.roundness)};
        --zitadel-spacing: ${this.getSpacingValue(theme.spacing)};
        --zitadel-appearance: ${theme.appearance};
        
        /* Component-specific variables */
        --zitadel-button-roundness: ${this.getComponentRoundness("button", theme)};
        --zitadel-card-roundness: ${this.getComponentRoundness("card", theme)};
        --zitadel-input-roundness: ${this.getComponentRoundness("input", theme)};
        
        /* Color system variables */
        --zitadel-primary: #3b82f6;
        --zitadel-primary-hover: #2563eb;
        --zitadel-secondary: #6b7280;
        --zitadel-success: #10b981;
        --zitadel-error: #ef4444;
        --zitadel-warning: #f59e0b;
      }
    `;
  }

  injectCustomerOverrides(customerCSS: string): void {
    const styleElement = document.createElement("style");
    styleElement.textContent = customerCSS;
    styleElement.dataset.zitadelCustomerOverrides = "true";
    document.head.appendChild(styleElement);
  }
}

// Customer usage
const bridge = new ZitadelCSSBridge();

const customerCSS = `
  :root {
    /* Override our defaults */
    --zitadel-primary: #8b5cf6;
    --zitadel-primary-hover: #7c3aed;
    --zitadel-button-roundness: 12px;
  }
  
  .zitadel-button {
    background: var(--zitadel-primary);
    border-radius: var(--zitadel-button-roundness);
    font-family: 'Your Custom Font', sans-serif;
  }
  
  .zitadel-button:hover {
    background: var(--zitadel-primary-hover);
  }
`;

bridge.injectCustomerOverrides(customerCSS);
```

#### **Configuration-Based Override System**

```typescript
// Customer configuration that generates CSS overrides
interface CustomerStyleConfig {
  cssOverrides: {
    [selector: string]: CSSStyleDeclaration | string;
  };
  componentOverrides: {
    [componentName: string]: {
      classes?: string[];
      styles?: CSSStyleDeclaration | string;
    };
  };
  customCSS?: string;
}

const customerConfig: CustomerStyleConfig = {
  cssOverrides: {
    ".zitadel-button": {
      background: "linear-gradient(45deg, #ff6b6b, #4ecdc4)",
      border: "none",
      borderRadius: "8px",
      color: "white",
    },
    ".zitadel-input:focus": {
      borderColor: "#4ecdc4",
      boxShadow: "0 0 0 3px rgba(78, 205, 196, 0.1)",
    },
  },
  componentOverrides: {
    button: {
      classes: ["custom-button", "gradient-bg"],
      styles: "padding: 16px 32px; font-weight: 700;",
    },
    card: {
      styles: "box-shadow: 0 20px 40px rgba(0,0,0,0.1); border: none;",
    },
  },
  customCSS: `
    @keyframes slideIn {
      from { transform: translateX(-100%); }
      to { transform: translateX(0); }
    }
    
    .zitadel-login-form {
      animation: slideIn 0.3s ease-out;
    }
  `,
};
```

### **Implementation Strategy**

1. **Phase 1**: Expose CSS custom properties for all theme values
2. **Phase 2**: Add data attributes to all components for specific targeting
3. **Phase 3**: Implement CSS injection API for runtime overrides
4. **Phase 4**: Build configuration-based override system
5. **Phase 5**: Create visual CSS override builder for non-technical users

This approach would give enterprise customers complete styling control while maintaining the simplicity of our theme system for basic customization needs.

A drag-and-drop interface where organizations could:

- Preview theme changes in real-time
- Export theme configurations
- Share branded templates across teams
- Validate accessibility compliance automatically

## Conclusion

Building a comprehensive theme system has transformed how organizations can approach their authentication interfaces. What started as a need for simple customization has evolved into a powerful, extensible platform that puts brand identity at the forefront of user experience.

The benefits extend far beyond aesthetics. By providing a consistent, predictable theming system, we've:

- **Reduced development overhead** for organizations wanting custom interfaces
- **Improved user experience** through cohesive design language
- **Future-proofed** the system for evolving design needs
- **Democratized customization** for non-technical team members

As we continue to evolve this system, we're excited about the possibilities that lie ahead. The foundation we've built today enables a future where every login experience can be as unique as the organization behind it, without sacrificing security, performance, or maintainability.

---

_Want to see the theme system in action? Check out our [live demo](https://login.zitadel.app) or explore the [implementation details](https://github.com/zitadel/zitadel) in our open-source repository._

**Key Takeaways:**

- ✅ Environment variable-driven configuration
- ✅ Component-specific defaults with global override capability
- ✅ Future-proof architecture avoiding utility function patterns
- ✅ Comprehensive coverage of all UI elements
- ✅ Type-safe, developer-friendly integration
- ✅ Extensible foundation for advanced customization features

---
