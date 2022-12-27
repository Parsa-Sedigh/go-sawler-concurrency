# Section 9: 9. Adding Concurrency to Choosing a Plan

## 77-1. What we'll cover in this section
We're gonna write a handler that fires off a couple of goroutines. One goroutine will generate an invoice(and will run in the background OFC) and we'll
also have another goroutine fired off that will open a PDF file, a user manual. So we'll open an existing PDF and we'lll modify it, save that PDF and then
send it off to the user as an attachment and these things will run concurrently. We will also subscribe the user to the plan.

## 78-2. Getting the plan id, the plan, and the user
## 79-3. Generating an Invoice
## 80-4. Generating a manual
Install: `phpdave11/gofpdf` which allows us to create a PDF and `phpdave11/gofpdf/contrib/gofpdi` allows us to open an existing PDF and use that as a template.

We need to create a folder called tmp because the generating invoice code needs it(it's gonna write PDFs to that directory).

## 81-5. Trying things out, subscribing a user, updating  the session, and redirecting
 After running the docker containers of this app, run: `make start` or `make restart` to start it up.
 
In the MIME tab of mailhog, you can see the PDF files(attachments) to download.

Why we just didn't put the two goroutines in SubscribeToPlan in one goroutine?
Because this way, generating an invoice and generating a manual, they both run concurrently. If we make them 11 goroutine, they would run sequentially. This way,
we have them running at the same time. The only situation where we couldn't do this(two goroutines approach), is that if we needed some info from the first
goroutine in order to make the second goroutine run.