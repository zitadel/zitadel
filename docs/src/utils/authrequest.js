import React, { useState, useEffect } from "react";

export const AuthRequestContext = React.createContext(null);

export default ({ children }) => {
  const [instance, setInstance] = useState("your-instance");
  const [clientId, setClientId] = useState("your-client-id");
  const [redirectUri, setRedirectUri] = useState("your-redirect-uri");
  const [responseType, setResponseType] = useState("your-response-type");
  const [scope, setScope] = useState("your-scope");

  const [prompt, setPrompt] = useState("your-prompt");
  const [authMethod, setAuthMethod] = useState("your-auth-method");
  const [codeChallenge, setCodeChallenge] = useState("your-code-challenge");
  const [codeVerifier, setCodeVerifier] = useState("your-code-verifier");
  const [loginHint, setLoginHint] = useState("your-login-hint");
  const [idTokenHint, setIdTokenHint] = useState("your-id-token-hint");
  const [organizationId, setOrganizationId] = useState("your-organization-id");

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);

    const instance_param = params.get("instance");
    const client_id = params.get("client_id");
    const redirect_uri = params.get("redirect_uri");
    const response_type = params.get("response_type");
    const scope_param = params.get("scope");

    // optional parameters
    const prompt_param = params.get("prompt");
    const auth_method_param = params.get("auth_method");
    const code_verifier_param = params.get("code_verifier");
    const login_hint = params.get("login_hint");
    const id_token_hint = params.get("id_token_hint");
    const organization_id = params.get("organization_id");

    setInstance(instance_param ?? "https://mydomain-xyza.zitadel.cloud/");
    setClientId(client_id ?? "170086824411201793@yourapp");
    setRedirectUri(
      redirect_uri ?? "http://localhost:8080/api/auth/callback/zitadel"
    );
    setResponseType(response_type ?? "code");
    setScope(scope_param ?? "openid email profile");
    setPrompt(prompt_param ?? "");
    setAuthMethod(auth_method_param ?? "(none) PKCE");
    setCodeVerifier(code_verifier_param ?? "random-string");
    // optional parameters
    setLoginHint(login_hint ?? "johndoe@example.zitadel.cloud");
    setIdTokenHint(id_token_hint ?? "[your-id-token]");
    setOrganizationId(organization_id ?? "168811945419506433");

    if (
      instance_param ||
      client_id ||
      redirect_uri ||
      response_type ||
      scope_param ||
      prompt_param ||
      organization_id ||
      login_hint ||
      id_token_hint
    ) {
      const example = document.getElementById("example");
      if (example) {
        example.scrollIntoView();
      }
    }
  }, []);

  const authRequest = {
    instance: [instance, setInstance],
    clientId: [clientId, setClientId],
    redirectUri: [redirectUri, setRedirectUri],
    responseType: [responseType, setResponseType],
    scope: [scope, setScope],
    prompt: [prompt, setPrompt],
    authMethod: [authMethod, setAuthMethod],
    codeChallenge: [codeChallenge, setCodeChallenge],
    codeVerifier: [codeVerifier, setCodeVerifier],
    loginHint: [loginHint, setLoginHint],
    idTokenHint: [idTokenHint, setIdTokenHint],
    organizationId: [organizationId, setOrganizationId],
  };

  return (
    <AuthRequestContext.Provider value={authRequest}>
      {children}
    </AuthRequestContext.Provider>
  );
};
