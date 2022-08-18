from behave import *

@given('a valid set exists with id "{set_id}"')
def step_impl(context, set_id):
    print(f'Set ID: {set_id}')
    return

