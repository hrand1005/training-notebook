## About

This directory is for tests that exceed the scope of simple unit testing, which in classic go fashion 
exist alongside the source code they exercise, ending with `*_test.go`.

## behave

The behave tests use Gherkin language (given-when-then statements) to define integration test cases
in plain english. Thusly, and in contrast to some advice I've encountered for api-testing, the tests do
not strictly inspect http response bodies to test compliance with json schema specifications. Instead 
these tests verify (and document) under what conditions a client can access a resource, and other more
complicated use cases that may require more setup than unit tests can achieve.
The tests answer questions such as...
- can I access this resource without authentication?
- what if I'm authenticated but not the owner of this resource?
- what if the resource is labeled public? private?
- what if the resource doesn't exist?

The test cases are defined in the `/features` folder and use the file extension `.feature`.
