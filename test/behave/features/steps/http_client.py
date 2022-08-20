from behave import *
import requests

@when('the client sends GET request to "{endpoint}"')
def step_impl(context, endpoint):
    for k, v in context.scenario.url_params.items():
        endpoint = endpoint.replace(k, v)

    token = context.scenario.token
    headers = {
            'Authorization': f'Bearer {token}'
            }
    context.scenario.response = requests.get(context.base_url+endpoint, headers=headers)

@then('the server responds with status code "{status_code}"')
def step_impl(context, status_code):
    got_status = int(status_code)
    expected_status = context.scenario.response.status_code 
    assert got_status == expected_status, f'GOT: {got_status}\nEXPECTED: {expected_status}'

