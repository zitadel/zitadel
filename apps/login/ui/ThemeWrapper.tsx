import { getBranding } from "#/lib/zitadel";
import { server } from "../lib/zitadel";

const ThemeWrapper = async ({ children }: any) => {
  console.log("hehe");

  const defaultClasses = "bg-background-light-600 dark:bg-background-dark-600";

  try {
    const policy = await getBranding(server);

    const darkStyles = {
      backgroundColor: `${policy?.backgroundColorDark}`,
      color: `${policy?.fontColorDark}`,
    };

    const lightStyles = {
      backgroundColor: `${policy?.backgroundColor}`,
      color: `${policy?.fontColor}`,
    };

    console.log(policy);

    return (
      <div className={defaultClasses} style={darkStyles}>
        {children}
      </div>
    );
  } catch (error) {
    console.error(error);

    return <div className={defaultClasses}>{children}</div>;
  }
};

export default ThemeWrapper;
