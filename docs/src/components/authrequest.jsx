import React, { useContext, useEffect } from "react";
import { AuthRequestContext } from "../utils/authrequest";
import styles from "../css/authrequest.module.css";

export function SetAuthRequest() {
  const {
    clientId: [clientId, setClientId],
    redirectUri: [redirectUri, setRedirectUri],
    responseType: [responseType, setResponseType],
    scope: [scope, setScope],
    idTokenHint: [idTokenHint, setIdTokenHint],
  } = useContext(AuthRequestContext);

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);

    // required authorize parameters
    const client_id = params.get("client_id");
    const redirect_uri = params.get("redirect_uri");
    const response_type = params.get("response_type");
    const scope_param = params.get("scope");

    // optional parameters
    const id_token_hint = params.get("id_token_hint");

    setClientId(client_id ?? "[your-client-id]");
    setRedirectUri(redirect_uri ?? "[your-redirect-uri]");
    setResponseType(response_type ?? "[your-response-type]");
    setScope(scope_param ?? "[scope]");
    setIdTokenHint(id_token_hint ?? "[your-id-token-hint]");
  }, []);

  return (
    <div>
      <div className={styles.inputwrapper}>
        <label className={styles.label}>Client ID</label>
        <input
          className={styles.input}
          id="client_id"
          value={clientId}
          onChange={(event) => {
            const value = event.target.value;
            // setClientId(value);
          }}
        />
      </div>

      <div className={styles.inputwrapper}>
        <label className={styles.label}>Redirect URI</label>
        <input
          className={styles.input}
          id="redirect_uri"
          value={redirectUri}
          onChange={(event) => {
            const value = event.target.value;
            // setRedirectUri(value);
          }}
        />
      </div>
    </div>
  );
}
