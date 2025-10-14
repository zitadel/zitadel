"use client";

import { Logo } from "@/components/logo";
import { BrandingSettings } from "@zitadel/proto/zitadel/settings/v2/branding_settings_pb";
import React, { ReactNode, Children } from "react";
import { ThemeWrapper } from "./theme-wrapper";
import { Card } from "./card";
import { useResponsiveLayout } from "@/lib/theme-hooks";

/**
 * DynamicTheme component handles layout switching between traditional top-to-bottom
 * and modern side-by-side layouts based on NEXT_PUBLIC_THEME_LAYOUT.
 *
 * For side-by-side layout:
 * - First child: Goes to left side (title, description, etc.)
 * - Second child: Goes to right side (forms, buttons, etc.)
 * - Single child: Falls back to right side for backward compatibility
 *
 * For top-to-bottom layout:
 * - All children rendered in traditional centered layout
 */
export function DynamicTheme({
  branding,
  children,
}: {
  children: ReactNode | ((isSideBySide: boolean) => ReactNode);
  branding?: BrandingSettings;
}) {
  const { isSideBySide } = useResponsiveLayout();

  // Resolve children immediately to avoid passing functions through React
  const actualChildren: ReactNode = React.useMemo(() => {
    if (typeof children === "function") {
      return (children as (isSideBySide: boolean) => ReactNode)(isSideBySide);
    }
    return children;
  }, [children, isSideBySide]);

  return (
    <ThemeWrapper branding={branding}>
      {isSideBySide
        ? // Side-by-side layout: first child goes left, second child goes right
          (() => {
            const childArray = Children.toArray(actualChildren);
            const leftContent = childArray[0] || null;
            const rightContent = childArray[1] || null;

            // If there's only one child, it's likely the old format - keep it on the right side
            const hasLeftRightStructure = childArray.length === 2;

            return (
              <div className="relative mx-auto w-full max-w-[1100px] py-4 px-8">
                <Card>
                  <div className="flex min-h-[400px]">
                    {/* Left side: First child + branding */}
                    <div className="flex w-1/2 flex-col justify-center p-4 lg:p-8 bg-gradient-to-br from-primary-50 to-primary-100 dark:from-primary-900/20 dark:to-primary-800/20">
                      <div className="max-w-[440px] mx-auto space-y-8">
                        {/* Logo and branding */}
                        {branding && (
                          <Logo
                            lightSrc={branding.lightTheme?.logoUrl}
                            darkSrc={branding.darkTheme?.logoUrl}
                            height={150}
                            width={150}
                          />
                        )}

                        {/* First child content (title, description) - only if we have left/right structure */}
                        {hasLeftRightStructure && (
                          <div className="space-y-4 text-left flex flex-col items-start">
                            {/* Apply larger styling to the content */}
                            <div className="space-y-6 [&_h1]:text-4xl [&_h1]:lg:text-4xl [&_h1]:text-left [&_h1]:text-gray-900 [&_h1]:dark:text-white [&_h1]:leading-tight [&_p]:text-left [&_p]:leading-relaxed [&_p]:text-gray-700 [&_p]:dark:text-gray-300">
                              {leftContent}
                            </div>
                          </div>
                        )}
                      </div>
                    </div>

                    {/* Right side: Second child (form) or single child if old format */}
                    <div className="flex w-1/2 items-center justify-center p-4 lg:p-8">
                      <div className="w-full max-w-[440px]">
                        <div className="space-y-6">{hasLeftRightStructure ? rightContent : leftContent}</div>
                      </div>
                    </div>
                  </div>
                </Card>
              </div>
            );
          })()
        : // Traditional top-to-bottom layout - center title/description, left-align forms
          (() => {
            const childArray = Children.toArray(actualChildren);
            const titleContent = childArray[0] || null;
            const formContent = childArray[1] || null;
            const hasMultipleChildren = childArray.length > 1;

            return (
              <div className="relative mx-auto w-full max-w-[440px] py-4 px-4">
                <Card>
                  <div className="mx-auto flex flex-col items-center space-y-8">
                    <div className="relative flex flex-row items-center justify-center -mb-4">
                      {branding && (
                        <Logo
                          lightSrc={branding.lightTheme?.logoUrl}
                          darkSrc={branding.darkTheme?.logoUrl}
                          height={150}
                          width={150}
                        />
                      )}
                    </div>

                    {hasMultipleChildren ? (
                      <>
                        {/* Title and description - center aligned */}
                        <div className="w-full text-center flex flex-col items-center mb-4">{titleContent}</div>

                        {/* Form content - left aligned */}
                        <div className="w-full">{formContent}</div>
                      </>
                    ) : (
                      // Single child - use original behavior
                      <div className="w-full">{actualChildren}</div>
                    )}

                    <div className="flex flex-row justify-between"></div>
                  </div>
                </Card>
              </div>
            );
          })()}
    </ThemeWrapper>
  );
}
