# lob

## Problem Statement

Here at Lob we have a dashboard where you can create a postcard. In it, there is an address field where you enter the recipient's address. We'd like
to provide suggestions based on the partial string the user has typed so far based on the addresses in a pre-defined data-set. In real-life there are
many strategies to accomplish this, but for the purposes of this interview, we will use an API endpoint that takes partial strings and returns a list
of matches. The format & ordering are up to you, the API's designer. For this exercise we have provided a sample data-set in addresses.json.

### Requirements

Please design and implement a set of APIs that allow us to:

1. Query with a partial address string and receive all matching addresses in the data-set.
2. Manage the addresses in the data-set.

These requirements are deliberately vague. Sometimes in engineering we have ill-defined problems and it is up to us to provide clarity. Use your best
judgment and we would like you to call out any decisions and/or assumptions you have made in the phone call portion of the interview.

### Example Responses

For the input "1600", we expect to get the 3 addresses with 1600 Holloway Ave from the data-set.

For the input "MD" we expect to get the 3400 N. Charles St., Baltimore, MD 21218 address.

### Acceptance Criteria

Your server must run and accept requests from the API client of your choice (localhost is fine) and we'll ask you to demonstrate it working and
discuss your approach during the follow up call.

***We will be asking for you to develop a small incremental feature during the call, so make sure it is easy and comfortable to modify and run your
code in your environment for the follow up call.***

# Instructions (Linux)

To compile & run, from project root run:

`(cd src; go build; ./lob -addr localhost:8080 -file api/testdata/addresses.json)`

To test, run from another terminal:

`curl http://localhost:8080/address?query=1600 | jq`
