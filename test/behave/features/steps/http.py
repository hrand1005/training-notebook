from behave import *

@when('the client sends request "{endpoint}"')
def step_impl(context, endpoint):
    # add endpoint to context base_url
    print(f'Endpoint: {endpoint}')
    return

@then('the server responds with status code "{status_code}"')
def step_impl(context, status_code):
    print(f'Status code: {status_code}')
    return
