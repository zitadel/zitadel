import React, { useState } from "react";

export const EnvironmentContext = React.createContext(null);

export default ({ children }) => {
  const [instance, setInstance] = useState("your-instance");
  const [clientId, setClientId] = useState("your-client-id");

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
