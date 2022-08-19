from behave import *

def before_all(context):
    """
    before_all gets called by behave before executing features/scenarios/steps
    """
    server_config = context.config.userdata.get("server_config")
    print(f'Server Config: {server_config}')
    # context.base_url = config[]
