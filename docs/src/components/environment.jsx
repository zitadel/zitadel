import React, { useContext, useEffect } from "react";
import { EnvironmentContext } from "../utils/environment";
import styles from "../css/environment.module.css";
import Interpolate from "@docusaurus/Interpolate";
import CodeBlock from "@theme/CodeBlock";

export function SetEnvironment() {
  const {
    instance: [instance, setInstance],
    clientId: [clientId, setClientId],
  } = useContext(EnvironmentContext);

  useEffect(() => {
    const params = new URLSearchParams(window.location.search); // id=123
    const clientId = params.get("clientId");
    const instance = params.get("instance");

    const localClientId = localStorage.getItem("clientId");
    const localInstance = localStorage.getItem("instance");

    setClientId(clientId ?? localClientId ?? "");
    setInstance(instance ?? localInstance ?? "");
  }, []);

  function setAndSaveInstance(value) {
    if (instance !== value) {
      localStorage.setItem("instance", value);
      setInstance(value);
    }
  }

  function setAndSaveClientId(value) {
    if (clientId !== value) {
      localStorage.setItem("clientId", value);
      setClientId(value);
    }
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
              localStorage.removeItem("instance");
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
              setAndSaveClientId(value);
            } else {
              localStorage.removeItem("clientId");
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

  return <div>{variable}</div>;
}

export function EnvInterpolate({ children }) {
  const {
    instance: [instance],
    clientId: [clientId],
  } = useContext(EnvironmentContext);

  return (
    <Interpolate
      values={{
        clientId,
        instance,
      }}
    >
      {children}
    </Interpolate>
  );
}

export function EnvCode({
  language,
  title,
  code,
  showLineNumbers = false,
  children,
}) {
  const {
    instance: [instance],
    clientId: [clientId],
  } = useContext(EnvironmentContext);

  return (
    <CodeBlock
      language={language}
      title={title}
      showLineNumbers={showLineNumbers}
    >
      <Interpolate
        values={{
          clientId,
          instance,
        }}
      >
        {children}
      </Interpolate>
    </CodeBlock>
  );
}
