from behave import *

INVALID_USER_ID = "0" * 24

@given('a user exists with id "{_user_id}"')
def step_impl(context, _id):
    context.scenario.url_params["_user_id"] = ""

@given('the user does not exist')
def step_impl(context):
    context.scenario.url_params["_user_id"] = INVALID_USER_ID

@given('the authenticated client is not that user')
def step_impl(context):
    user_id = create_new_user(context)
    context.scenario.url_params["_user_id"] = user_id
    # set token for a different user
    different_user_id = create_new_user(context)
    different_user_token = login_with_user(context, different_user_id)
    context.scenario.token = different_user_token

@given('the authenticated client is that user')
def step_impl(context):
    user_id = create_new_user(context)
    context.scenario.url_params["_user_id"] = user_id
    # set auth token for that user
    context.scenario.token = login_with_user(context, user_id)

def create_new_user(context):
    # TODO: IMPLEMENT BDD
    # create user with post request
    # user_id = requests.post(
    #   url=context.base_url+"/users",
    #   json=user_data
    # )
    user_id = INVALID_USER_ID
    return user_id

def login_with_user(context, user_id):
    # TODO: IMPLEMENT BDD
    # resp = requests.post(
    #   url=context.base_url+"/login",
    #   json=username_and_password
    # )
    # token = token_from_resp(resp)
    token = "token"
    return token
