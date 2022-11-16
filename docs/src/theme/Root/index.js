import React from "react";
import EnvironmentProvider from "../../utils/environment";

// Default implementation, that you can customize
export default function Root({ children }) {
  return <EnvironmentProvider>{children}</EnvironmentProvider>;
}
