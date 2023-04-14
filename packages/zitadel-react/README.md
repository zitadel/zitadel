# How to use

### Install

```sh
npm install @zitadel/react
```

or

```sh
yarn add @zitadel/react
```

### Import styles file

To get the styles, import them in `_app.tsx` or global styling file

```
import "@zitadel/react/styles.css";
```

### Setup Dark mode

to set dark theme, wrap your components in a `ui-dark` class.

### Use components

```tsx
import { SignInWithGoogle } from "@zitadel/react";

export default function IdentityProviders() {
  return (
    <div className="py-4">
      <SignInWithGoogle />
    </div>
  );
}
```
