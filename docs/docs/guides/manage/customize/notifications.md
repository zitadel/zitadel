---
title: SMS, SMTP and HTTP Provider for Notifications 
---

All Notifications send as SMS and Email are customizable as that you can define your own providers, 
which then send the notifications out. These providers can also be defined as an HTTP type, 
and the text and content, which is used to send the SMS's and Emails will get send to a Webhook as JSON.

With this everything can be customized or even custom logic can be implemented to use a not yet supported provider by ZITADEL.

## How it works

There is a default provider configured in ZITADEL Cloud, both for SMS's and Emails, but this default providers can be changed through the respective API's.

This API's are provided on an instance level:
- [SMS Providers](/apis/resources/admin/sms-provider)
- [Email Providers](/apis/resources/admin/email-provider)

To use a non-default provider just add, and then activate. There can only be 1 provider be activated at the same time.

## Resulting messages

In case of the Twilio and SMTP providers, the messages will be sent as before, in case of the HTTP providers the content of the messages is the same but as a HTTP call.

Here an example of the body of an Email sent via HTTP provider:
```json
{
  "contextInfo": {
    "eventType": "user.human.initialization.code.added",
    "provider": {
      "id": "285181292935381355",
      "description": "test"
    },
    "recipientEmailAddress": "example@zitadel.com"
  },
  "templateData": {
    "title": "Zitadel - Initialize User",
    "preHeader": "Initialize User",
    "subject": "Initialize User",
    "greeting": "Hello GivenName FamilyName,",
    "text": "This user was created in Zitadel. Use the username Username to login. Please click the button below to finish the initialization process. (Code 0M53RF) If you didn't ask for this mail, please ignore it.",
    "url": "http://example.zitadel.com/ui/login/user/init?authRequestID=\u0026code=0M53RF\u0026loginname=Username\u0026orgID=275353657317327214\u0026passwordset=false\u0026userID=285181014567813483",
    "buttonText": "Finish initialization",
    "primaryColor": "#5469d4",
    "backgroundColor": "#fafafa",
    "fontColor": "#000000",
    "fontFamily": "-apple-system, BlinkMacSystemFont, Segoe UI, Lato, Arial, Helvetica, sans-serif",
    "footerText": "InitCode.Footer"
  },
  "args": {
    "changeDate": "2024-09-16T10:58:50.73237+02:00",
    "code": "0M53RF",
    "creationDate": "2024-09-16T10:58:50.73237+02:00",
    "displayName": "GivenName FamilyName",
    "firstName": "GivenName",
    "lastEmail": "example@zitadel.com",
    "lastName": "FamilyName",
    "lastPhone": "+41791234567",
    "loginNames": [
      "Username"
    ],
    "nickName": "",
    "preferredLoginName": "Username",
    "userName": "Username",
    "verifiedEmail": "example@zitadel.com",
    "verifiedPhone": ""
  }
}
```

There are 3 elements to this message:
- contextInfo, with information on why this message is sent like the Event, which Email or SMS provider is used and which recipient should receive this message
- templateData, with all texts and format information which can be used with a template to produce the desired message
- args, with the information provided to the user which can be used in the message to customize 
