from behave import *
import yaml


def load_server_settings(path: str) -> dict:
    """
    load_server_config parses a yaml config file. 
    Assumed to match the form of training-notebook server configs.
    """
    with open(path) as f:
        server_config = yaml.load(f, Loader=yaml.FullLoader)

    return server_config.get('server-settings')

def before_all(context):
    """
    before_all gets called by behave before executing features/scenarios/steps.
    Initializes behave context for data common to all scenarios.
    """
    config_path = context.config.userdata.get('server_config')
    server_settings = load_server_settings(config_path)
    context.base_url = server_settings.get('host') + server_settings.get('port') 

def before_scenario(context, scenario):
    context.scenario.url_params = {}
