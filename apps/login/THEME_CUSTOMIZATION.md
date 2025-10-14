# Login App Theme Customization Guide

This guide helps you customize the appearance of your ZITADEL login application for a personalized user experience.

## Quick Start

1. Copy the theme variables you want to customize from `.env.theme.example` to your `.env.local` file
2. Restart your application
3. Your theme changes will be applied automatically

## Theme Options

### 🔄 Roundness (`NEXT_PUBLIC_THEME_ROUNDNESS`)

Controls how rounded the UI elements appear:

- **`edgy`** - Sharp, rectangular corners (modern tech, corporate)
- **`mid`** - Medium rounded corners (balanced, professional)
- **`full`** - Fully rounded corners (friendly, approachable)

### 📱 Layout (`NEXT_PUBLIC_THEME_LAYOUT`)

Controls the overall page structure:

- **`side-by-side`** - Brand section on left, form on right (desktop view)
- **`top-to-bottom`** - Brand section on top, form below (mobile-first)

### 📏 Spacing (`NEXT_PUBLIC_THEME_SPACING`)

Controls spacing and padding throughout the interface:

- **`regular`** - Standard spacing with comfortable padding (p-6, space-y-6)
- **`compact`** - Tighter spacing for information-dense layouts (p-4, space-y-4)

### 🎨 Appearance (`NEXT_PUBLIC_THEME_APPEARANCE`)

Complete design philosophies:

- **`flat`** - Minimal flat design with cards matching background color, subtle borders, and normal typography
- **`material`** - Material Design inspired with elevated cards, proper contrast, and medium typography

### 🖼️ Background Image (`NEXT_PUBLIC_THEME_BACKGROUND_IMAGE`)

Add a custom background image:

- Use local images: `/images/my-background.jpg` (place in `public/images/`)
- Use external URLs: `https://example.com/background.jpg`
- Leave empty for solid color backgrounds

## Example Configurations

### Tech Startup

```env
NEXT_PUBLIC_THEME_ROUNDNESS=full
NEXT_PUBLIC_THEME_LAYOUT=side-by-side
NEXT_PUBLIC_THEME_SPACING=regular
NEXT_PUBLIC_THEME_APPEARANCE=material
NEXT_PUBLIC_THEME_BACKGROUND_IMAGE=/images/tech-gradient.jpg
```

### Corporate Bank

```env
NEXT_PUBLIC_THEME_ROUNDNESS=edgy
NEXT_PUBLIC_THEME_LAYOUT=top-to-bottom
NEXT_PUBLIC_THEME_SPACING=regular
NEXT_PUBLIC_THEME_APPEARANCE=material
```

### Minimal SaaS

```env
NEXT_PUBLIC_THEME_ROUNDNESS=mid
NEXT_PUBLIC_THEME_LAYOUT=side-by-side
NEXT_PUBLIC_THEME_SPACING=compact
NEXT_PUBLIC_THEME_APPEARANCE=flat
```

### Creative Agency

```env
NEXT_PUBLIC_THEME_ROUNDNESS=full
NEXT_PUBLIC_THEME_LAYOUT=top-to-bottom
NEXT_PUBLIC_THEME_SPACING=regular
NEXT_PUBLIC_THEME_APPEARANCE=material
NEXT_PUBLIC_THEME_BACKGROUND_IMAGE=/images/creative-workspace.jpg
```

### Flat Design App

```env
NEXT_PUBLIC_THEME_ROUNDNESS=mid
NEXT_PUBLIC_THEME_LAYOUT=side-by-side
NEXT_PUBLIC_THEME_SPACING=compact
NEXT_PUBLIC_THEME_APPEARANCE=flat
```

### Material Design System

```env
NEXT_PUBLIC_THEME_ROUNDNESS=mid
NEXT_PUBLIC_THEME_LAYOUT=side-by-side
NEXT_PUBLIC_THEME_SPACING=regular
NEXT_PUBLIC_THEME_APPEARANCE=material
```

## Advanced Customization

For more detailed customization beyond these presets, you can:

1. **Custom CSS**: Add your own CSS files in the `src/styles/` directory
2. **Component Override**: Modify the theme configuration in `src/lib/theme.ts`
3. **Custom Appearances**: Add new appearance options to the `APPEARANCE_STYLES` or `SPACING_STYLES` objects

## Troubleshooting

### Theme not applying

- Ensure you're using `NEXT_PUBLIC_` prefix for all theme variables (`NEXT_PUBLIC_THEME_ROUNDNESS`, `NEXT_PUBLIC_THEME_SPACING`, `NEXT_PUBLIC_THEME_APPEARANCE`, etc.)
- Restart your development server after changing environment variables
- Check that your `.env.local` file is in the root of the login app directory

### Background image not loading

- Verify the image path is correct
- For external URLs, ensure the domain is accessible
- Check browser console for any loading errors

### Layout issues on mobile

- Test your theme on different screen sizes
- The `top-to-bottom` layout is generally more mobile-friendly
- Some combinations work better on certain screen sizes

## Support

For additional customization needs or questions, please refer to the ZITADEL documentation or community forums.
