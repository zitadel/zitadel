export type ZitadelNextProps = {
  dark: boolean;
  children: React.ReactNode;
};

export function ZitadelNextProvider({ dark, children }: ZitadelNextProps) {
  return (
    <div className={`${dark ? "ztdl-dark" : "ztdl-light"} `}>{children}</div>
  );
}
