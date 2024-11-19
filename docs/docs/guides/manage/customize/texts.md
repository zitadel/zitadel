---
title: Customized Texts
---

You are able to customize the texts used from ZITADEL. This is possibly on the instance or organization level.

## Message Texts

Sometimes the users will get an email or phone message from ZITADEL (e.g Password Reset Request).
ZITADEL already has some good standard texts, but maybe you would like to customize it for your organization.

Go to the message text policy on your organization and you will find the different kinds of messages that are sent from ZITADEL. 
Choose the template and the language you like to edit. 
You can now change all the texts from a message. 
As soon as you click into a input field you will see some attribute chips below the field. 
These are the parameters you can include on this specific message.

![Message Texts](/img/console_message_texts.png)

## Login Texts

Like the message texts you are also able to change the texts on the login interface. 
First choose the screen and the language you like to edit. 
You will see the default texts in the input field and you can overwrite them by typing into the box.

![Message Texts](/img/console_login_texts.png)

## Reset to default

If you don't like your customization anymore click the "reset policy" button.
All your settings will be removed and the default settings of the system will trigger.

## Internationalization / i18n

ZITADEL is available in the following languages

- German (de)
- English (en)
- Spanish (es)
- French (fr)
- Indonesian (id)
- Italian (it)
- 日本語 (ja)
- Polish（pl）
- 简体中文（zh）
- Bulgarian (bg)
- Portuguese (pt)
- Macedonian (mk)
- Czech (cs)
- Russian (ru)
- Dutch (nl)
- Swedish (sv)
- Hungarian (hu)

A language is displayed based on your agent's language header.
If a users language header doesn't match any of the supported or [restricted](#restrict-languages) languages, the instances default language will be used.

If you need support for a specific language we highly encourage you to [contribute translation files](https://github.com/zitadel/zitadel/blob/main/CONTRIBUTING.md) for the missing language.

## Restrict Languages

If you only want to enable a subset of the supported languages, you can configure the languages you'd like to allow using the [restrictions API](./restrictions.md).
The login UI and notification messages are only rendered in one of the allowed languages and fallback to the instances default language.
Also, the instances OIDC discovery endpoint will only list the allowed languages in the *ui_locales_supported* field.

All language settings are also configurable in the consoles *Languages* default settings.

![Languages](/img/guides/console/languages.png)
