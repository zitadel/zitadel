import React, { useState, useEffect } from "react";

export const AuthRequestContext = React.createContext(null);

export default ({ children }) => {
  const [instance, setInstance] = useState("your-instance");
  const [clientId, setClientId] = useState("your-client-id");
  const [redirectUri, setRedirectUri] = useState("your-redirect-uri");
  const [responseType, setResponseType] = useState("your-response-type");
  const [scope, setScope] = useState("your-scope");
  const [idTokenHint, setIdTokenHint] = useState("your-id-token-hint");

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const instance_param = params.get("instance");
    const client_id = params.get("client-id");
    const redirect_uri = params.get("redirect-uri");
    const response_type = params.get("response_type");
    const scope_param = params.get("scope");

    // optional parameters
    const id_token_hint = params.get("id_token_hint");

    setInstance(instance_param ?? "https://mydomain-xyza.zitadel.cloud/");
    setClientId(client_id ?? "170086824411201793@yourapp");
    setRedirectUri(
      redirect_uri ?? "http://localhost:8080/api/auth/callback/zitadel"
    );
    setResponseType(response_type ?? "code");
    setScope(scope_param ?? "openid email profile");

    // optional parameters
    setIdTokenHint(id_token_hint ?? "[your-id-token-hint]");
  }, []);

  const authRequest = {
    clientId: [clientId, setClientId],
    redirectUri: [redirectUri, setRedirectUri],
    responseType: [responseType, setResponseType],
    scope: [scope, setScope],
    idTokenHint: [idTokenHint, setIdTokenHint],
  };

  return (
    <AuthRequestContext.Provider value={authRequest}>
      {children}
    </AuthRequestContext.Provider>
  );
};
