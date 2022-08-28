from behave import *
from users import create_new_user, login_with_user

INVALID_SET_ID = "0" * 24

### SCENARIO SCOPED STEPS ###

@given('a set exists with id "{_set_id}"')
def step_impl(context, _id):
    context.scenario.url_params["_set_id"] = ""


@given('no set exists with id "{_set_id}"')
def step_impl(context, set_id):
    context.scenario.url_params["_set_id"] = INVALID_SET_ID


@given('the authenticated client is not the owner of the set')
def step_impl(context):
    user_id = create_new_user(context)
    token = login_with_user(context, user_id)
    set_id = create_new_set(context, token)
    context.scenario.url_params["_set_id"] = set_id

    # SET TOKEN of non-owner
    non_owner_id = create_new_user(context)
    non_owner_token = login_with_user(context, non_owner_id)
    context.scenario.token = non_owner_token


@given('the authenticated client is the owner of the set')
def step_impl(context):
    user_id = create_new_user(context)
    token = login_with_user(context, user_id)
    set_id = create_new_set(context, token)
    # SET TOKEN of the owner
    context.scenario.token = token
    context.scenario.url_params["_set_id"] = set_id


def create_new_set(context, token):
    # TODO: IMPLEMENT BDD
    # create user with post request
    # set_id = requests.post(
    #   url=context.base_url+"/sets",
    #   json=set_data
    # )
    set_id = INVALID_SET_ID
    return set_id
