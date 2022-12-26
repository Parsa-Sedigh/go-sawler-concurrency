# Section 8: 

## 69-1. What we'll cover in this section

Users and plans

- adding & validating user accounts: When a user registers, we'll send them an email with a link in it that they have to click on, in order to
activate their account and for security purposes, that URL, the one that's included in the email, it needs to be signed, so that it can't be
tampered with
- signed URLs in email
- dispelling the list of available subscriptions

## 70-2. Adding mail templates and URL signer code
Let's create an email template. Name it `confirmation-email.html.gohtml` and `confirmation-email.plain.gohtml` .

When we send the activation email, the link is gonna be sth like http://domain/activate?email=<email> . But this is a security loophole. It means
that **anyone** can just start guessing emails to activate. So we need to make the url tamper proof. Create a file named signer.go and copy  the contents
of the resources of this lesson into it.

The github.com/bwmarrin/go-alone package allows us to generate signed text of any sort. In our case, we're gonna sign our URL that we put in email that
we send off with the link to activate the account.

YYou're gonna use a much longer and more secure secret that is used for generating a signed text and also you're not gonna store it in your code,
you'll read it from an environment variable or sth like that.

About `GenerateTokenFromString`: If we pass it a url, it will hand back a url with a hash append to the end.
About `VerifyToken`: We can pass a signed text to it(like a signed link) and that func will verify that it matches the signature that we would
generate from our code. That function will append a hash to the string and later we can look at that hash to see if we generated it ourselves?
This prevents url tamper.

With Expired method, we can have links that expire after some time.

## 71-3. Starting on the handler to create a user
Restart the app by running: `make restart` and test the code of this lesson by registering an account(reload the page on /register too).

## 72-4. Activating a user
## 73-5. Giving user data to our templates

## 74-6. Displaying the Subscription Plans page
After a user chooses a plan(or if you manually insert a row in user_plans table), because of the session, user needs to login again to see
his "selected plan".

## 75-7. Adding a route and trying things out for the Plans page

## 76-8. Writing a stub handler for choosing a plan