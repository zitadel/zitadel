import { SignInWithGoogle, SignInWithGitlab } from "@zitadel/react";

export default function IdentityProviders() {
  return (
    <div className="space-y-4 py-4">
      <SignInWithGoogle />
      <SignInWithGitlab />
    </div>
  );
}
