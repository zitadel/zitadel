import { useEffect, useState } from "react";

// Custom hook to read auth record and user profile doc
export function useUserData() {
  const [clientData, setClientData] = useState(null);

  useEffect(() => {
    let unsubscribe;

    return unsubscribe;
  }, [clientData]);

  return { clientData };
}
