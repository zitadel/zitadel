export type ZitadelUIProps = {
  dark: boolean;
  children: React.ReactNode;
};

export function ZitadelUIProvider({ dark, children }: ZitadelUIProps) {
  return <div className={`${dark ? "ui-dark" : "ui-light"} `}>{children}</div>;
}
