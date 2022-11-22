import React, { useState, useEffect } from "react";

export const EnvironmentContext = React.createContext(null);

export default ({ children }) => {
  const [instance, setInstance] = useState("your-instance");
  const [clientId, setClientId] = useState("your-client-id");

  useEffect(() => {
    const params = new URLSearchParams(window.location.search); // id=123
    const clientId = params.get("clientId");
    const instance = params.get("instance");

    const localClientId = localStorage.getItem("clientId");
    const localInstance = localStorage.getItem("instance");

    setClientId(clientId ?? localClientId ?? "");
    setInstance(instance ?? localInstance ?? "");
  }, []);

  const environment = {
    instance: [instance, setInstance],
    clientId: [clientId, setClientId],
  };

  return (
    <EnvironmentContext.Provider value={environment}>
      {children}
    </EnvironmentContext.Provider>
  );
};
