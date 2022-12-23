# Section 7: 7. Sending Email Concurrently
 
## 61-1. What we'll cover in this section
Why we would want to send email in the background?

Sending emails can slow things down. If you're connection to a mail service and it's having a bad day for example you're connected to mailgun and for
some reason it's slow to send emails. You don't want whatever you're doing in your code to just stop until that email get sent. Instead,
you wanna sent it of in the background and we're gonna do this using 3 channels but only 2 that are involved with actually sending emails.

We'll send information off to a channel and that will have sth lightening  to that channel and when it receives it, it'll send it off in the background(it''ll
fire off a new goroutine and that sends our email off concurrently). The other channel will listen for errors  and if sth goes wrong, then we send
sth in that error channel and a third one is just used to shut things down.

About cleaning up logic:

If our app is sending emails and doing all sorts of things in the background and someone decides: Hey I need to stop the app for whatever reason,
you don't want to just stop because some things that were queued to happen in the background or that were actually happening gin the background,
they'll just die and you'll never find out about it! So we need to wait until everything that's running in the background is finished

## 62-2. Getting started with the mailer code
Create mailer.go . Install `github.com/vanng822/go-premailer/premailer` and this inline CSS to make email moore compatible with the various
email clients out there when you're sending html formatted email.

For an actual mail package, install: `/xhit/go-simple-mail/v2` (use v2!).

We have mail.html.gohtml and mail.plain.gohtml . Why? Because when we send email, we don't want to send the plain text email as the **only** means
oof communication with our customers. So we'll actually send it in **two** formats. One plain text which is visible to any email client anywhere on
the planet and another one that is actuality gonna be used by the vast majority of them and nicely formatted html message.

## 63-3. Building HTML and Plain Text messages
## 64-4. Sending a message (synchronously)
Let's send email synchronously(not in background) and once we're source, we'll send emails using goroutines.

For this, create `test-email` route.

By creating the handler of test-email, to test that route, assuming we have our docker images running in the background, when we go to
/test-email in the browser, we should get a blank screen, but we should also get an email sent(if everything is ok). 

To run the project:
```shell
make start # (if you already ran make start, run `make restart`)
```

After visiting localhost/test-email, if you go to localhost:8025 to see the mailhog UI, you can see the sent email.

Sending email works, now we need to send email using goroutines in the background.

## 65-5. Getting started sending a message (asynchronously)

## 66-6. Writing a helper function to send email easily

## 67-7. Sending an email on incorrect login

## 68-8. Adding cleanup tasks to the shutdown() function
After running the `docker-compose up -d` , run `make restart` and go to browser and try logging in with wrong credentials and then see if you get an email.
