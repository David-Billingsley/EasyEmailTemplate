#EasyEmail
    To use this package put var Email EasyEmail.Email.  
    After this is in function you can call the functions by doing Email.Email_Body_Only()

    This libary is being devloped to help make it easier to send emails using go.  
    This libary is under devolpment if there is a send method you would like, please add it into the issue tracker.

## Email_Body_Only
    This function generates an email with body only using the email template file provided.

Parms:
(Name - Data type - Defenition)
* Sender - String - Who is sending the email
* Password - String - Senders Password
* SmtpAdd - string - Senders email SMTP address
*  smtpHost - string - Senders email SMTP host
* smtpPort - string - Senders email SMTP Port
* templname - string - the .html template file to use for the email body.
* recivers - []string - the list of recivers
* subject - string - the subject of the email
* bodytext - string - the text you would like passed to the body

Returns:
    String or Error
    
Example: 
<!-- template.html
<!DOCTYPE html>
<html>

<body>
    {{.Message}}
</body>

</html>
-->

## Email_W_Attachments
    Send email with an attachment

    Parms:
(Name - Data type - Defenition)
  * Sender - String - Who is sending the email
  * Password - String - Senders Password
  * SmtpAdd - string - Senders email SMTP address
  * smtpHost - string - Senders email SMTP host
  * smtpPort - string - Senders email SMTP Port
  * recivers - []string - the list of recivers
  * subject - string - the subject of the email
  * bodytext - string - the text you would like passed to the body
  * attachmentPath - string - the path of the file you would like to attach

Returns:
    String or Error