import React, { useContext, useEffect } from "react";
import { AuthRequestContext } from "../utils/authrequest";
import styles from "../css/authrequest.module.css";
import CodeBlock from "@theme/CodeBlock";

export function SetAuthRequest() {
  const {
    clientId: [clientId, setClientId],
    redirectUri: [redirectUri, setRedirectUri],
    responseType: [responseType, setResponseType],
    scope: [scope, setScope],
    idTokenHint: [idTokenHint, setIdTokenHint],
  } = useContext(AuthRequestContext);

  //   useEffect(() => {
  //     const params = new URLSearchParams(window.location.search);

  //     // required authorize parameters
  //     const client_id = params.get("client_id");
  //     const redirect_uri = params.get("redirect_uri");
  //     const response_type = params.get("response_type");
  //     const scope_param = params.get("scope");

  //     // optional parameters
  //     const id_token_hint = params.get("id_token_hint");

  //     setClientId(client_id ?? "[your-client-id]");
  //     setRedirectUri(redirect_uri ?? "[your-redirect-uri]");
  //     setResponseType(response_type ?? "[your-response-type]");
  //     setScope(scope_param ?? "[your-scope]");
  //     setIdTokenHint(id_token_hint ?? "[your-id-token-hint]");
  //   }, []);

  return (
    <>
      <h5>Required Parameters</h5>

      <div className="grid grid-cols-4">
        <div className={styles.inputwrapper}>
          <label className={styles.label}>Client ID</label>
          <input
            className={styles.input}
            id="client_id"
            value={clientId}
            onChange={(event) => {
              const value = event.target.value;
              setClientId(value);
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
              setRedirectUri(value);
            }}
          />
        </div>

        <div className={styles.inputwrapper}>
          <label className={styles.label}>Response Type</label>
          <input
            className={styles.input}
            id="response_type"
            value={responseType}
            onChange={(event) => {
              const value = event.target.value;
              setResponseType(value);
            }}
          />
        </div>

        <div className={styles.inputwrapper}>
          <label className={styles.label}>Scope</label>
          <input
            className={styles.input}
            id="scope"
            value={scope}
            onChange={(event) => {
              const value = event.target.value;
              setScope(value);
            }}
          />
        </div>
      </div>

      <h5>Optional Parameters</h5>

      <div className={styles.grid}>
        <div className={styles.inputwrapper}>
          <label className={styles.label}>Id Token Hint</label>
          <input
            className={styles.input}
            id="id_token_hint"
            value={idTokenHint}
            onChange={(event) => {
              const value = event.target.value;
              setIdTokenHint(value);
            }}
          />
        </div>
      </div>

      <br />

      <CodeBlock language="ts"></CodeBlock>
    </>
  );
}
