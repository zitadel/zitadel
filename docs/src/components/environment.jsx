import React, { useContext, useEffect } from "react";
import { EnvironmentContext } from "../utils/environment";
import styles from "../css/environment.module.css";

export function SetEnvironment() {
  const {
    instance: [instance, setInstance],
    clientId: [clientId, setClientId],
  } = useContext(EnvironmentContext);

  useEffect(() => {
    const params = new URLSearchParams(window.location.search); // id=123
    const clientId = params.get("clientId");
    const instance = params.get("instance");

    setClientId(clientId);
    setInstance(instance);
  }, []);

  function setAndSaveInstance(value) {
    setInstance(value);
  }

  return (
    <div>
      <div className={styles.inputwrapper}>
        <label className={styles.label}>Your instance domain</label>
        <input
          className={styles.input}
          id="instance"
          value={instance}
          onChange={(event) => {
            const value = event.target.value;
            if (value) {
              setAndSaveInstance(value);
            } else {
              setInstance("");
            }
          }}
        />
      </div>

      <br />

      <div className={styles.inputwrapper}>
        <label className={styles.label}>Client ID</label>
        <input
          className={styles.input}
          id="clientid"
          value={clientId}
          onChange={(event) => {
            const value = event.target.value;
            if (value) {
              setClientId(value);
            } else {
              setClientId("");
            }
          }}
        />
      </div>
    </div>
  );
}

export function Env({ name }) {
  const env = useContext(EnvironmentContext);
  const variable = env[name];

  return variable ? <span>{variable}</span> : null;
}
