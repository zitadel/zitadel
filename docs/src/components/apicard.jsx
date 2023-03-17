import React from "react";

import styles from "../css/apicard.module.css";

export function ApiCard({ type, children }) {
  let classes = "";
  switch (type) {
    case "AUTH":
      classes = "bg-green-500/10 dark:bg-green-500/20";
      break;
    case "MGMT":
      classes = "bg-blue-500/10 dark:bg-blue-500/20";
      break;
    case "ADMIN":
      classes = "bg-red-500/10 dark:bg-red-500/20";
      break;
    case "SYSTEM":
      classes = "bg-yellow-500/10 dark:bg-yellow-500/20";
      break;
    case "ASSET":
      classes = "bg-black/10 dark:bg-black/20";
      break;
  }

  return (
    <div
      className={`${styles.apicard} flex mb-4 flex-row p-4 rounded-lg ${classes} `}
    >
      {children}
    </div>
  );
}
