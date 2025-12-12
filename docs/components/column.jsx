import React from "react";

export default function Column({ children }) {
  return (
    <div className="grid grid-cols-1 gap-4 lg:grid-cols-2">{children}</div>
  );
}
