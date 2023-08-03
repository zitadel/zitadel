export type ZitadelReactProps = {
  dark: boolean;
  children: React.ReactNode;
};

export function ZitadelReactProvider({ dark, children }: ZitadelReactProps) {
  return (
    <div className={`${dark ? "ztdl-dark" : "ztdl-light"} `}>{children}</div>
  );
}
