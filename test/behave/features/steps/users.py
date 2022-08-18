from behave import *

@given('a valid user exists with id "{user_id}"')
def step_impl(context, user_id):
    print(f'User ID: {user_id}')
    return

