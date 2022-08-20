from behave import *
from users import create_new_user, login_with_user
    
@given('the client is not authenticated')
def step_impl(context):
    context.scenario.token = None

@given('the client is authenticated')
def step_impl(context):
    """
    NOTE: Should only be used to log in with some client without expectations
    for specific access to any restricted resources.
    """
    user_id = create_new_user(context)
    token = login_with_user(context, user_id)
    context.scenario.client_id = user_id
    context.scenario.token = token


