---
title: Customized Texts
---

You can customize the texts sent by ZITADEL.

## Message Texts
Sometimes users get an email or phone message from ZITADEL (e.g Password Reset Request).
ZITADEL already has some good standard texts, but you might like to customize them for your organization.
To do that, follow these steps:

1. Go to the message text policy in your organization. You will find the different kinds of messages ZITADEL sends. 
1. Choose the template and the language you like to edit. 
1. Change all the texts as you wish.
  As soon as you select an input field, you will see some attribute chips below the field. 
  These are the parameters you can include for this specific message.

![Message Texts](/img/console_message_texts.png)

## Login Texts

Like the message texts, you can also change the texts on the login interface. 

1. Choose the screen and the language you would like to edit. 
   You will see the default texts in the input field.
1. Overwrite the defaults by typing into the box.

![Message Texts](/img/console_login_texts.png)

## Reset to default

If you don't like your customization anymore, select the **reset policy** button.
All your settings will be removed and the default settings of the system will trigger.

## Internationalization

ZITADEL's support for languages will be extended with time. 
If you need support for a specific language, we highly recommend you to write translation files for the missing language.

ZITADEL loads translations from three files:

 - [Console translations](https://github.com/zitadel/zitadel/tree/main/console/src/assets/i18n)
 - [Login interface texts](https://github.com/zitadel/zitadel/tree/main/internal/ui/login/static/i18n)
 - [Email Notifcation texts](https://github.com/zitadel/zitadel/tree/main/internal/notification/static/i18n)
 - [Common translations](https://github.com/zitadel/zitadel/tree/main/internal/static/i18n) for success or error toasts

 Make sure you set the locale as the name. Later on, a language header will determine which file gets displayed.

