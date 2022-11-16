import React from "react";
import { EnvironmentContext } from "../utils/environment";

export function Environment() {
  const { clientId, instance, setClientId, setInstance } =
    React.useContext(EnvironmentContext);

  return (
    <div>
      <div>
        <label>Your instance domain</label>
        <input
          id="instance"
          value={instance}
          onChange={(value) => setInstance(value.target.value)}
        />
      </div>

      <div>
        <label>Client ID</label>
        <input
          id="clientid"
          value={clientId}
          onChange={(value) => setClientId(value.target.value)}
        />
      </div>
    </div>
  );
}
