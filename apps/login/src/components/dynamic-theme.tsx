"use client";

import { Logo } from "@/components/logo";
import { BrandingSettings } from "@zitadel/proto/zitadel/settings/v2/branding_settings_pb";
import React, { ReactNode, Children } from "react";
import { AppAvatar } from "./app-avatar";
import { ThemeWrapper } from "./theme-wrapper";
import { Card } from "./card";
import { useResponsiveLayout } from "@/lib/theme";

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
  appName,
}: {
  children: ReactNode | ((isSideBySide: boolean) => ReactNode);
  branding?: BrandingSettings;
  appName?: string;
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
      {isSideBySide ? (
        // Side-by-side layout: first child goes left, second child goes right
        (() => {
          const childArray = Children.toArray(actualChildren);
          const leftContent = childArray[0] || null;
          const rightContent = childArray[1] || null;

          // If there's only one child, it's likely the old format - keep it on the right side
          const hasLeftRightStructure = childArray.length === 2;

          return (
            <div className="relative mx-auto w-full max-w-[1200px] py-8">
              <Card>
                <div className="p-6">
                  <div className="flex min-h-[400px]">
                    {/* Left side: First child + branding */}
                    <div className="flex w-1/2 flex-col justify-center p-4 lg:p-8 bg-gradient-to-br from-primary-50 to-primary-100 dark:from-primary-900/20 dark:to-primary-800/20">
                      <div className="max-w-[440px] mx-auto space-y-8">
                        {/* Logo and branding */}
                        {branding && (
                          <div className="space-y-4">
                            <Logo
                              lightSrc={branding.lightTheme?.logoUrl}
                              darkSrc={branding.darkTheme?.logoUrl}
                              height={60}
                              width={120}
                            />
                            {appName && (
                              <h2 className="text-lg font-semibold text-gray-700 dark:text-gray-300">{appName}</h2>
                            )}
                          </div>
                        )}

                        {/* First child content (title, description) - only if we have left/right structure */}
                        {hasLeftRightStructure && (
                          <div className="space-y-4 text-left flex flex-col items-start">{leftContent}</div>
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
                </div>
              </Card>
            </div>
          );
        })()
      ) : (
        // Traditional top-to-bottom layout - keep center alignment
        <div className="relative mx-auto w-full max-w-[440px] py-8">
          <Card>
            <div className="mx-auto flex flex-col items-center space-y-4">
              <div className="relative flex flex-row items-center justify-center gap-8">
                {branding && (
                  <>
                    <Logo
                      lightSrc={branding.lightTheme?.logoUrl}
                      darkSrc={branding.darkTheme?.logoUrl}
                      height={appName ? 100 : 150}
                      width={appName ? 100 : 150}
                    />

                    {appName && <AppAvatar appName={appName} />}
                  </>
                )}
              </div>

              <div className="w-full flex flex-col items-center text-center space-y-4">{actualChildren}</div>
              <div className="flex flex-row justify-between"></div>
            </div>
          </Card>
        </div>
      )}
    </ThemeWrapper>
  );
}
