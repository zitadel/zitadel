import ThemeWrapper from "./ThemeWrapper";

type Props = {
  children: React.ReactNode;
};

export function LayoutProviders({ children }: Props) {
  return (
    // <ThemeProvider
    //   attribute="class"
    //   defaultTheme="system"
    //   storageKey="cp-theme"
    //   value={{ dark: "dark" }}
    // >
    /* @ts-expect-error Server Component */
    <ThemeWrapper>{children}</ThemeWrapper>
    // </ThemeProvider>
  );
}
