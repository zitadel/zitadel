import React from "react";
import AuthRequestProvider from "../../utils/authrequest";
// import EnvironmentProvider from "../../utils/environment";

// Default implementation, that you can customize
export default function Root({ children }) {
  return (
    // <EnvironmentProvider>
    <AuthRequestProvider>{children}</AuthRequestProvider>
    // </EnvironmentProvider>
  );
}
