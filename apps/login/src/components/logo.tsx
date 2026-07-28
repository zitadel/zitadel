import { getThemeConfig } from "@/lib/theme";

type Props = {
  darkSrc?: string;
  lightSrc?: string;
  /**
   * Upper bound in px for the rendered height. Defaults to the theme config
   * value (NEXT_PUBLIC_THEME_LOGO_MAX_HEIGHT).
   */
  maxHeight?: number;
};

export function Logo({ lightSrc, darkSrc, maxHeight }: Props) {
  const { logoMaxHeight } = getThemeConfig();
  const cap = maxHeight ?? logoMaxHeight;

  // Constrain both axes instead of forcing fixed dimensions: the intrinsic
  // aspect ratio of the uploaded asset is preserved, so wide wordmarks use the
  // available width while tall marks stay within the height budget. Uploaded
  // logos have no guaranteed aspect ratio, so anything that pins one axis to a
  // constant either distorts the mark or shrinks it to illegibility.
  const style = { maxHeight: `${cap}px`, maxWidth: "100%" };
  const className = "h-auto w-auto object-contain";

  return (
    <>
      {darkSrc && (
        <div className="hidden items-center dark:flex">
          <img className={className} style={style} src={darkSrc} alt="logo" />
        </div>
      )}
      {lightSrc && (
        <div className="flex items-center dark:hidden">
          <img className={className} style={style} src={lightSrc} alt="logo" />
        </div>
      )}
    </>
  );
}
