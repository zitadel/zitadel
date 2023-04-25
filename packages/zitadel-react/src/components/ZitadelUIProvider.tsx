export type ZitadelUIProps = {
  dark: boolean;
  children: React.ReactNode;
};

export function ZitadelUIProvider({ dark, children }: ZitadelUIProps) {
  return (
    <div className={`${dark ? "ztdl-dark" : "ztdl-light"} `}>{children}</div>
  );
}
