---
title: Customized Texts
---

You are able to customize the texts used from ZITADEL.

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

## Internationalization

ZITADELs support for languages will be extended with time. 
If you need support for a specific language we highly recommend you to write translation files for the missing language.

We have a command to generate translations automatically using [Deepl APIs](https://www.deepl.com/translator). You need to signup and create an API token. The free account should be sufficient to generate all necessary files for ZITADEL. You can get the command [here](https://github.com/caos/zitadel/blob/main/guides/development.md).

ZITADEL loads translations from three files:

 - [Console translations](https://github.com/caos/zitadel/tree/main/console/src/assets/i18n)
 - [Login interface texts](https://github.com/caos/zitadel/tree/main/internal/ui/login/static/i18n)
 - [Email Notifcation texts](https://github.com/caos/zitadel/tree/main/internal/notification/static/i18n)
 - [Common translations](https://github.com/caos/zitadel/tree/main/internal/static/i18n) for success or error toasts

 Make sure you set the correct locale as the name. Later on, language header will determine which file gets displayed.

