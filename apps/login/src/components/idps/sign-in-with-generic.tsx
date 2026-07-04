"use client";

import { forwardRef, type CSSProperties } from "react";
import {
  siFacebook,
  siKakaotalk,
  siLine,
  siNaver,
  siQq,
  siSinaweibo,
  siTiktok,
  siVk,
  siWechat,
  siX,
  siZalo,
  type SimpleIcon,
} from "simple-icons";
import { BaseButton, SignInWithIdentityProviderProps } from "./base-button";

// Generic OAuth2/OIDC identity providers have no dedicated branded button, so
// upstream renders them as a plain name while native providers (Google, GitHub,
// Apple, …) show their logo. Match the configured display name to a well-known
// brand glyph (simple-icons) and render the real logo — making social/regional
// providers (Facebook, LINE, WeChat, KakaoTalk, Naver, X, TikTok, QQ, Weibo,
// VK, Zalo) instantly recognisable. Names that don't match keep the original
// name-only button, so there is no behavioural change for unknown providers.
//
// Matches use word boundaries so a token can't hit a substring of an unrelated
// name (e.g. "meta" must not match "Metadata SSO").
const BRAND_ICONS: { match: RegExp; icon: SimpleIcon }[] = [
  { match: /\bfacebook\b|\bmeta\b/i, icon: siFacebook },
  { match: /\bwechat\b|\bweixin\b|微信/i, icon: siWechat },
  { match: /\bkakao(talk)?\b/i, icon: siKakaotalk },
  { match: /\bnaver\b/i, icon: siNaver },
  { match: /\btiktok\b|\bdouyin\b/i, icon: siTiktok },
  { match: /\bweibo\b|微博/i, icon: siSinaweibo },
  { match: /\bzalo\b/i, icon: siZalo },
  { match: /\bline\b/i, icon: siLine },
  { match: /\bqq\b|\btencent\b/i, icon: siQq },
  { match: /\bvk\b|\bvkontakte\b/i, icon: siVk },
  { match: /\b(x|twitter)\b/i, icon: siX },
];

function brandIcon(name?: string): SimpleIcon | undefined {
  if (!name) return undefined;
  return BRAND_ICONS.find((b) => b.match.test(name))?.icon;
}

// simple-icons ship a single brand colour, which can be illegible against one
// theme: near-black brands (X, TikTok #000) vanish on the dark background, and
// near-white ones (KakaoTalk #FFCD00) vanish on the light one. Keep the brand
// colour where it has contrast and only swap the illegible end per theme. The
// fills are exposed as CSS variables + Tailwind `dark:` so the theme switch is
// pure CSS (no JS theme read → no hydration flash).
function brandFills(hex: string): { light: string; dark: string } {
  const r = parseInt(hex.slice(0, 2), 16);
  const g = parseInt(hex.slice(2, 4), 16);
  const b = parseInt(hex.slice(4, 6), 16);
  const luminance = (0.2126 * r + 0.7152 * g + 0.0722 * b) / 255;
  const brand = `#${hex}`;
  return {
    light: luminance > 0.7 ? "#18181b" : brand, // too light for the white login bg
    dark: luminance < 0.2 ? "#ffffff" : brand, // too dark for the dark login bg
  };
}

export const SignInWithGeneric = forwardRef<HTMLButtonElement, SignInWithIdentityProviderProps>(
  function SignInWithGeneric(props, ref) {
    const { children, name = "", className, ...restProps } = props;
    const icon = brandIcon(name);

    // Recognised brand: mirror the native branded buttons (Google/GitHub). The
    // icon row sets the button height (so no name-only "h-[50px]" default is
    // wanted here), and the visible name is the accessible label — the glyph
    // is therefore decorative (aria-hidden, no role/title) to avoid the screen
    // reader announcing the provider name twice.
    if (icon) {
      const fills = brandFills(icon.hex);
      return (
        <BaseButton {...restProps} ref={ref} className={className}>
          <div className="flex h-12 w-12 items-center justify-center">
            <svg
              aria-hidden="true"
              focusable="false"
              viewBox="0 0 24 24"
              className="h-6 w-6 fill-[var(--idp-fill)] dark:fill-[var(--idp-fill-dark)]"
              style={{ "--idp-fill": fills.light, "--idp-fill-dark": fills.dark } as CSSProperties}
              xmlns="http://www.w3.org/2000/svg"
            >
              <path d={icon.path} />
            </svg>
          </div>
          {children ? children : <span className="ml-4">{name}</span>}
        </BaseButton>
      );
    }

    // Unknown provider: upstream name-only button with centred text (#12211).
    return (
      <BaseButton {...restProps} ref={ref} className={className ?? "h-[50px]"}>
        {children ? children : <span className="w-full text-center">{name}</span>}
      </BaseButton>
    );
  },
);
